package handlers

import (
	"net/http"
)

func (context *IsoContext) postVideo(w http.ResponseWriter, r *http.Request) {
	fileInfo, fileHeader, err := r.FormFile("file[]")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		http.Error(w, "invalid content-type", http.StatusBadRequest)
		return
	}

	item, err := context.VideoStorage.Save(fileHeader.Filename, fileInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = context.VideoRepository.Insert(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
