package db

import (
	"context"
	"fmt"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/common/enums/env"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var (
	YtSearchServiceDb *gorm.DB
)

func connectToDb(user, password, dbname, host, port string, ctx context.Context) *gorm.DB {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s", user, password, dbname, host, port)
	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	log.HandleErrorWithSuccessMessage(err, ctx, "DB connected - "+dbname)
	return db
}

func ConnectDatabase(ctx context.Context) {
	user := os.Getenv(env.PostgresUser)
	port := os.Getenv(env.PostgresPort)
	host := os.Getenv(env.PostgresHost)
	dbname := os.Getenv(env.PostgresDatabase)
	password := os.Getenv(env.PostgresPassword)
	YtSearchServiceDb = connectToDb(user, password, dbname, host, port, ctx)
	gormAutoMigrations(ctx)
}
