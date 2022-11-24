package feed

import (
	"errors"
	"sync"
)

// TODO spin up and connect to a postgress DB
type PostgressRepository struct {
	feeds map[int]Feed
	sync.Mutex
}

func NewPostgressRepository() *PostgressRepository {
	feeds:= make(map[int]Feed)
	feed1 := Feed{
			ID:         1,
			Vendor:     "GoodRx",
			FeedName:   "Claims",
			FeedMethod: "SFTP",
		}
	feeds[feed1.ID] = feed1

	feed2 :=  Feed{
		ID:         2,
		Vendor:     "The Advisory Board",
		FeedName:   "Physicians",
		FeedMethod:	"S3",
	}
	feeds[feed2.ID] = feed2

	return &PostgressRepository{
		feeds: feeds,
	}
}

func (imr *PostgressRepository) GetFeeds() ([]Feed, error) {
	var feeds []Feed
	for _, feed := range imr.feeds {
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (imr *PostgressRepository) GetFeed(id int) (Feed, error) {
	feed, ok := imr.feeds[id]
	if ok {
		return feed, nil
	}
	return Feed{}, errors.New("no such feed exists")
}

func (imr *PostgressRepository) UpdateFeed(feed Feed) (Feed, error) {
	_, ok := imr.feeds[feed.ID]
	if ok {
		imr.feeds[feed.ID] = feed
		return feed, nil
	}
	return Feed{}, errors.New("no such feed exists")
}

func (imr *PostgressRepository) AddFeed(feed Feed) (Feed, error) {
	_, ok := imr.feeds[feed.ID]
	if ok {
		return Feed{}, errors.New("Feed ID already exists")
	}
	imr.feeds[feed.ID] = feed
	return feed, nil
}