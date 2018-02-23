package video

import (
	"github.com/segmentio/ksuid"
	"io"
	"os"
	"path"
)

const indexFileName = "index.mp4"
const thumbnailFileName = "screen.jpg"

type StorageInterface interface {
	Save(fileName string, writer io.Reader) (*Item, error)
	GetFilePath(item *Item) string
	GetThumbnailPath(item *Item) string
}

type storage struct {
	contentFolderPath string
}

func CreateStorage(contentFolderPath string) StorageInterface {
	storage := new(storage)
	storage.contentFolderPath = contentFolderPath
	return storage
}

func (s storage) GetFilePath(item *Item) string {
	dirPath := path.Dir(s.contentFolderPath)
	return dirPath + item.URL
}

func (s storage) GetThumbnailPath(item *Item) string {
	dirPath := path.Dir(s.contentFolderPath)
	return dirPath + item.Thumbnail
}

func (s storage) Save(fileName string, reader io.Reader) (*Item, error) {

	id := ksuid.New().String()
	dirPath := path.Dir(s.contentFolderPath)
	baseDir := path.Base(s.contentFolderPath) + "/" + id
	folderPath := dirPath + "/" + baseDir

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
		ID:        id,
		Name:      fileName,
		Duration:  0,
		Status:    StatusCreated,
		Thumbnail: "/" + baseDir + "/" + thumbnailFileName,
		URL:       "/" + baseDir + "/" + indexFileName,
	}
	return &item, nil
}
