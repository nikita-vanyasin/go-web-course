package video

import (
	"database/sql"
)

type RepositoryInterface interface {
	List(searchStringParam string, skip uint64, limit uint64) ([]*Item, error)
	RetrieveByKey(key string) (*Item, error)
	Insert(item *Item) error
}

type Repository struct {
	db *sql.DB
}

func CreateRepository(db *sql.DB) RepositoryInterface {
	var repo = new(Repository)
	repo.db = db
	return repo
}

func (repo *Repository) List(searchStringParam string, skip uint64, limit uint64) ([]*Item, error) {

	rows, err := repo.db.Query(`
       SELECT
		 video_key AS Id,
         title AS Name,
         duration AS Duration,
         thumbnail_url AS Thumbnail,
         url AS Url
       FROM
         video
	   WHERE
	     title LIKE CONCAT('%', ?, '%')
	   LIMIT ?, ?
    `, searchStringParam, skip, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*Item, 0)
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Name, &item.Duration, &item.Thumbnail, &item.Url)
		if err != nil {
			return nil, err
		}
		list = append(list, &item)
	}

	return list, nil
}

func (repo *Repository) RetrieveByKey(key string) (*Item, error) {
	rows, err := repo.db.Query(`
       SELECT
		 video_key AS Id,
         title AS Name,
         duration AS Duration,
         thumbnail_url AS Thumbnail,
         url AS Url
       FROM
		video
       WHERE
        video_key = ?
       LIMIT 1
    `, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item Item
	for rows.Next() {
		err := rows.Scan(&item.Id, &item.Name, &item.Duration, &item.Thumbnail, &item.Url)
		if err != nil {
			return nil, err
		}
	}
	return &item, nil
}

func (repo *Repository) Insert(item *Item) error {
	q := `INSERT INTO video (video_key, title, duration, thumbnail_url,  url)
        VALUES (?, ?, ?, ?, ?)`
	rows, err := repo.db.Query(q, item.Id, item.Name, item.Duration, item.Thumbnail, item.Url)
	if err == nil {
		rows.Close()
	}
	return err
}
