package feed

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	_ "github.com/lib/pq"
)

type PostgressRepository struct {
	db *sql.DB
	sync.Mutex
}

const (
	DB_HOST     = "postgres"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)


func NewPostgressRepository() *PostgressRepository {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	return &PostgressRepository{
		db: db,
	}
}

func (pr *PostgressRepository) GetFeeds() ([]Feed, error) {
	rows, err := pr.db.Query("SELECT id, vendor, feed_name, feed_method FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var f Feed
		err = rows.Scan(&f.ID, &f.Vendor, &f.FeedName, &f.FeedMethod)
		if err != nil {
		return nil, err
		}
		feeds = append(feeds, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (pr *PostgressRepository) GetFeed(id int) (Feed, error) {
	rows := pr.db.QueryRow("SELECT id, vendor, feed_name, feed_method FROM feeds WHERE id=$1;", id)
	var feed Feed
	err := rows.Scan(&feed.ID, &feed.Vendor, &feed.FeedName, &feed.FeedMethod)
	if err != nil {
		return Feed{}, err
	}
	return feed, nil
}

func (pr *PostgressRepository) UpdateFeed(feed Feed) (Feed, error) {
	sqlStatement := `
	UPDATE feeds
	SET vendor = $2, feed_name = $3, feed_method = $4
	WHERE id = $1;`
	_, err := pr.db.Exec(sqlStatement, feed.ID, feed.Vendor, feed.FeedName, feed.FeedMethod)
	if err != nil {
		return Feed{}, err
	}
	return feed, nil

}

func (imr *PostgressRepository) AddFeed(feed Feed) (Feed, error) {
	// _, ok := imr.feeds[feed.ID]
	// if ok {
	// 	return Feed{}, errors.New("Feed ID already exists")
	// }
	// imr.feeds[feed.ID] = feed
	// return feed, nil
	return Feed{}, errors.New("not yet implemented")

}

func (imr *PostgressRepository) DeleteFeed(id int) (string, error) {
	// _, ok := imr.feeds[id]
	// if ok {
	// 	delete(imr.feeds, id)
	// 	return fmt.Sprintf("Feed Id %x deleted", id), nil
	// }
	// return "", fmt.Errorf("Feed Id %x does not exist", id)
	return "", errors.New("not yet implemented")
}