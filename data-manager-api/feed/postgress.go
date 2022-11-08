package feed

import (
	"errors"
	"sync"
)

// TODO spin up and connect to a postgress DB

type PostgressRepository struct {
	feeds []Feed
	sync.Mutex
}

func NewPostgressRepository() *PostgressRepository {
	feeds := []Feed{
		{
			ID:         1,
			Vendor:     "GoodRx",
			FeedName:   "Claims",
			FeedMethod: "SFTP",
		},
		{
			ID:         2,
			Vendor:     "The Advisory Board",
			FeedName:   "Physicians",
			FeedMethod:	"S3",
		},
	}

	return &PostgressRepository{
		feeds: feeds,
	}
}

func (imr *PostgressRepository) GetFeeds() ([]Feed, error) {
	return imr.feeds, nil
}

func (imr *PostgressRepository) GetFeed(feedName string) (Feed, error) {
	for _, feed := range imr.feeds {
		if feed.FeedName == feedName {
			return feed, nil
		}
	}
	return Feed{}, errors.New("no such feed exists")
}

func (imr *PostgressRepository) AddFeed(feed Feed) (Feed, error) {
	imr.feeds = append(imr.feeds, feed)
	return Feed{}, errors.New("something went wrong adding feed")
}