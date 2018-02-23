package handlers

import (
	"database/sql"
	"github.com/nikita-vanyasin/go-web-course/common"
	"github.com/nikita-vanyasin/go-web-course/video"
)

type IsoContext struct {
	VideoRepository video.RepositoryInterface
	VideoStorage    video.StorageInterface

	db *sql.DB
}

func CreateContext(envSettings *common.EnvironmentSettings) *IsoContext {
	db := common.OpenSQLConnection(envSettings.SQLConnectionString)
	repo := video.CreateRepository(db)
	storage := video.CreateStorage(envSettings.ContentFolderPath)
	return &IsoContext{VideoRepository: repo, VideoStorage: storage, db: db}
}

func (context *IsoContext) Shutdown() error {
	return context.db.Close()
}
