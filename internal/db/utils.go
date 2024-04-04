package db

import (
	"context"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"sync"
)

func BulkInsertWithClause[T any](
	db *gorm.DB, modelInstruments []*T, blockSize int, blockInsertSuccessMessage string, clause clause.Expression,
	ctx context.Context,
) {
	wg := sync.WaitGroup{}
	for i := 0; i < len(modelInstruments); i += blockSize {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			blockEnd := int(math.Min(float64(index+blockSize), float64(len(modelInstruments))))
			block := modelInstruments[index:blockEnd]
			result := db.Clauses(clause).Create(block)
			if blockInsertSuccessMessage != "" {
				log.HandleErrorWithSuccessMessage(result.Error, ctx, blockInsertSuccessMessage)
			} else {
				log.HandleError(result.Error, ctx, false)
			}
		}(i)
	}
	wg.Wait()
}
