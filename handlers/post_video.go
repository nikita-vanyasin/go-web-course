package handlers

import (
	"github.com/segmentio/ksuid"
	"io"
	"net/http"
	"os"
)

func (context *IsoContext) postVideo(w http.ResponseWriter, r *http.Request) {
	fileInfo, fileHeader, err := r.FormFile("file[]")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ksuid.New().String()

	fileName := fileHeader.Filename
	folderPath := "content/" + id

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.OpenFile(folderPath+"/index.mp4", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, fileInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var item = VideoItem{
		Id:        id,
		Name:      fileName,
		Duration:  0,
		Thumbnail: "/" + folderPath + "/screen.jpg",
		Url:       "/" + folderPath + "/index.mp4",
	}

	db := context.DB

	q := `INSERT INTO video ( video_key, title, duration, thumbnail_url,  url)
         VALUES (?, ?, ?, ?, ?)`
	rows, err := db.Query(q, item.Id, item.Name, item.Duration, item.Thumbnail, item.Url)
	if err == nil {
		rows.Close()
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
