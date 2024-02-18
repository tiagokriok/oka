package repositories

import "github.com/tiagokriok/oka/internal/storages"

type Link struct {
	ID  string `json:"id,omitempty" param:"id"`
	URL string `json:"url"`
	Key string `json:"key,omitempty" param:"key"`
}

type LinkRepository struct {
	db *storages.MysqlDB
}

func NewLinkRepository(db *storages.MysqlDB) *LinkRepository {
	return &LinkRepository{
		db,
	}
}

func (lr *LinkRepository) CreateLink(link *Link) error {
	_, err := lr.db.Exec("INSERT INTO links (id, `key`, url) VALUES (?, ?, ?)", link.ID, link.Key, link.URL)
	if err != nil {
		return err
	}

	return nil
}

func (lr *LinkRepository) GetLinkByKey(key string) (*Link, error) {
	var link Link

	err := lr.db.Get(&link, "SELECT * FROM links WHERE `key`=?", key)
	if err != nil {
		return &link, err
	}

	return &link, nil
}
