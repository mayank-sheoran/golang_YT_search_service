package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"time"
)

type YoutubeDataService struct {
	Service            *youtube.Service
	Ctx                context.Context
	PreDefinedKeywords []string
	ApiKeys            []string
	FetchInterval      time.Duration
	LastFetchedAt      time.Time
}

var (
	client = &YoutubeDataService{
		Service: nil,
		Ctx:     nil,
		PreDefinedKeywords: []string{
			"gaming", "news",
		},
		ApiKeys:       []string{},
		FetchInterval: time.Minute,
		LastFetchedAt: time.Now().Add(time.Minute * -1),
	}
)

const (
	quotaExceededError = "googleapi: Error 403: The request cannot be completed because you have exceeded your <a href=\"/youtube/v3/getting-started#quota\">quota</a>., quotaExceeded"
)

func NewYoutubeDataV3Client() *YoutubeDataService {
	return client
}

func (yt *YoutubeDataService) setupNewApiKey() {
	var err error
	if len(yt.ApiKeys) == 0 {
		log.HandleError(errors.New("No valid YT API key remaining"), yt.Ctx, true)
	}
	yt.Service, err = youtube.NewService(yt.Ctx, option.WithAPIKey(yt.ApiKeys[0]))
	log.HandleError(err, yt.Ctx, true)
}

func (yt *YoutubeDataService) Run() {
	yt.Ctx = utils.GetContextWithFlowName(context.Background(), "Youtube Data service")
	yt.setupNewApiKey()
	log.Info("Youtube Data service started", yt.Ctx)
	for {
		var videosMetaData []*models.VideoMetaData
		for _, keyword := range yt.PreDefinedKeywords {
			videosMetaData = append(videosMetaData, yt.fetchLatestVideos(keyword)...)
		}
		db.BulkInsertWithClause(
			db.YtSearchServiceDb, videosMetaData, 5000,
			"Added videos published after: "+yt.LastFetchedAt.Format(time.RFC3339), models.VideoMetaDataUpsertClause(),
			yt.Ctx,
		)
		log.Info(
			"Waiting for next fetch interval at "+yt.LastFetchedAt.Add(yt.FetchInterval).Format(time.RFC3339), yt.Ctx,
		)
		time.Sleep(yt.FetchInterval)
	}
}

func (yt *YoutubeDataService) fetchLatestVideos(keyword string) []*models.VideoMetaData {
	log.Info(
		fmt.Sprintf(
			"Fetching latest videos for keyword: %s | published after: %s", keyword,
			yt.LastFetchedAt.Format(time.TimeOnly),
		), yt.Ctx,
	)
	call := yt.Service.Search.List([]string{"id,snippet"}).
		Q(keyword).
		Type("video").
		Order("date").
		PublishedAfter(yt.LastFetchedAt.Format(time.RFC3339)).
		MaxResults(10)
	yt.LastFetchedAt = time.Now().Add(time.Second)
	response, err := call.Do()
	if err != nil && err.Error() == quotaExceededError {
		log.Warn("Quote exceeded, changing API Key", yt.Ctx)
		yt.ApiKeys = yt.ApiKeys[1:]
		yt.setupNewApiKey()
		yt.fetchLatestVideos(keyword)
	}

	var videosMetaData []*models.VideoMetaData
	for _, item := range response.Items {
		parsedTime, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		log.HandleError(err, yt.Ctx, false)
		videosMetaData = append(
			videosMetaData, &models.VideoMetaData{
				ID:           item.Id.VideoId,
				Title:        item.Snippet.Title,
				Description:  item.Snippet.Description,
				PublishedAt:  parsedTime,
				ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
			},
		)
	}

	if len(videosMetaData) == 0 {
		log.Info(
			"No new videos found for keyword: "+keyword+" | published after: "+yt.LastFetchedAt.Format(time.TimeOnly),
			yt.Ctx,
		)
	}
	return videosMetaData
}
