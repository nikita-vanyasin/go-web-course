package video

import (
	"github.com/segmentio/ksuid"
	"io"
	"os"
)

const indexFileName = "index.mp4"
const thumbnailFileName = "screen.jpg"

type StorageInterface interface {
	Save(fileName string, writer io.Reader) (*Item, error)
}

type Storage struct {
	contentFolderPath string
}

func CreateStorage(contentFolderPath string) StorageInterface {
	storage := new(Storage)
	storage.contentFolderPath = contentFolderPath
	return storage
}

func (storage Storage) Save(fileName string, reader io.Reader) (*Item, error) {

	id := ksuid.New().String()
	folderPath := storage.contentFolderPath + "/" + id

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(folderPath+"/"+indexFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return nil, err
	}

	var item = Item{
		Id:        id,
		Name:      fileName,
		Duration:  0,
		Thumbnail: "/" + folderPath + "/" + thumbnailFileName,
		Url:       "/" + folderPath + "/" + indexFileName,
	}
	return &item, nil
}
