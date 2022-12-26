package feed

import (
	"database/sql"
	"fmt"
	"strconv"
	"os"
	"sync"
	_ "github.com/lib/pq"
)

type PostgressRepository struct {
	db *sql.DB
	sync.Mutex
}


var POSTGRES_HOST string = os.Getenv("POSTGRES_HOST")
var POSTGRES_PORT string = os.Getenv("POSTGRES_PORT")
var POSTGRES_USER string = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD string= os.Getenv("POSTGRES_PASSWORD")
var POSTGRES_DB_NAME string = os.Getenv("POSTGRES_DB_NAME")

func NewPostgressRepository() *PostgressRepository {
	port, err := strconv.Atoi(POSTGRES_PORT)
    if err != nil {
        panic(err)
    }
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    POSTGRES_HOST, port, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB_NAME)

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
	// Check if feed exists first
	rows := pr.db.QueryRow("SELECT id, vendor, feed_name, feed_method FROM feeds WHERE id=$1;", feed.ID)
	var f Feed
	e := rows.Scan(&f.ID, &f.Vendor, &f.FeedName, &f.FeedMethod)
	if e == sql.ErrNoRows {
		return Feed{}, fmt.Errorf("could not compelte update, ID %x does not exist", feed.ID)
	}
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

func (pr *PostgressRepository) AddFeed(feed Feed) (string, error) {
	sqlStatement := `
	INSERT INTO feeds (vendor, feed_name, feed_method)
	VALUES ($1, $2, $3);`

	_, err := pr.db.Exec(sqlStatement, feed.Vendor, feed.FeedName, feed.FeedMethod)
	if err != nil {
		return "", err
	}
	return "Feed Added", nil
}

func (pr *PostgressRepository) DeleteFeed(id int) (string, error) {
	sqlStatement := `
	DELETE FROM feeds
	WHERE id = $1;`

	_, err := pr.db.Exec(sqlStatement, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Feed Id %x deleted", id), nil
}


func (pr *PostgressRepository) GetFeedStatuses() ([]FeedStatusResults, error) {
	sqlStatement := `
	SELECT fs.id, fs.process_date, fs.record_count, fs.error_count, fs.feed_status,
	f.vendor, f.feed_name, f.feed_method
	FROM feed_status fs
	INNER JOIN feeds f
	ON fs.feed_id = f.id`
	rows, err := pr.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedStatus []FeedStatusResults
	for rows.Next() {
		var f FeedStatusResults
		err = rows.Scan(&f.ID, &f.ProcessDate, &f.RecordCount, &f.ErrorCount, &f.Status, &f.Vendor, &f.FeedName, &f.FeedMethod)
		if err != nil {
		return nil, err
		}
		feedStatus = append(feedStatus, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedStatus, nil
}

func (pr *PostgressRepository) GetFeedStatusesAggregate(startDate string, endDate string) (FeedStatusAggregate, error) {
	sqlStatement := `
	SELECT 
		COUNT(DISTINCT fs.id) AS files,
		SUM(fs.record_count) AS rows,
		SUM(fs.error_count) AS errors
	FROM feed_status fs
	WHERE process_date BETWEEN TO_DATE($1,'YYYY-MM-DD') AND TO_DATE($2,'YYYY-MM-DD')`
	rows := pr.db.QueryRow(sqlStatement, startDate, endDate)
	var fsAgg FeedStatusAggregate
	err := rows.Scan(&fsAgg.Files, &fsAgg.Rows, &fsAgg.Errors)
	if err != nil {
		return FeedStatusAggregate{}, err
	}
	return fsAgg, nil
}