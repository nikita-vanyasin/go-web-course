package video

import (
	"database/sql"
)

type RepositoryInterface interface {
	List(searchStringParam string, skip uint64, limit uint64) ([]*Item, error)
	RetrieveByKey(key string) (*Item, error)
	Insert(item *Item) error
	Update(item *Item) error
	GetUnprocessedItem() (*Item, error)
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
		 video_key AS ID,
         title AS Name,
         duration AS Duration,
         thumbnail_url AS Thumbnail,
         url AS URL,
         status AS Status
       FROM
         video
	   WHERE
	     title LIKE CONCAT('%', ?, '%') AND 
	     NOT status IN (?, ?)
	   LIMIT ?, ?
    `, searchStringParam, StatusError, StatusDeleted, skip, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*Item, 0)
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Duration, &item.Thumbnail, &item.URL, &item.Status)
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
		 video_key AS ID,
         title AS Name,
         duration AS Duration,
         thumbnail_url AS Thumbnail,
         url AS URL,
         status AS Status
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

	var item *Item
	for rows.Next() {
		item, err = populateItem(rows)
		if err != nil {
			return nil, err
		}
	}
	return item, nil
}

func (repo *Repository) Insert(item *Item) error {
	q := `
		INSERT INTO video (video_key, title, duration, thumbnail_url,  url, status)
        VALUES (?, ?, ?, ?, ?, ?)`
	rows, err := repo.db.Query(q, item.ID, item.Name, item.Duration, item.Thumbnail, item.URL, item.Status)
	if err == nil {
		rows.Close()
	}
	return err
}

func (repo *Repository) Update(item *Item) error {
	q := `
		UPDATE video
		SET 
			title = ?, duration = ?, thumbnail_url = ?, url = ?, status = ?
        WHERE video_key = ?
`
	rows, err := repo.db.Query(q, item.Name, item.Duration, item.Thumbnail, item.URL, item.Status, item.ID)
	if err == nil {
		rows.Close()
	}
	return err
}

func (repo *Repository) GetUnprocessedItem() (*Item, error) {
	rows, err := repo.db.Query(`
       SELECT
		 video_key AS ID,
         title AS Name,
         duration AS Duration,
         thumbnail_url AS Thumbnail,
         url AS URL,
         status AS Status
       FROM
		 video
       WHERE
         status = ?
       LIMIT 1
    `, StatusCreated)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item *Item
	for rows.Next() {
		item, err = populateItem(rows)
		if err != nil {
			return nil, err
		}
	}
	return item, nil
}

func populateItem(rows *sql.Rows) (*Item, error) {
	item := &Item{}
	err := rows.Scan(&item.ID, &item.Name, &item.Duration, &item.Thumbnail, &item.URL, &item.Status)
	return item, err
}
