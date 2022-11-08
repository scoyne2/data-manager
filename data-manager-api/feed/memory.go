package feed

import (
	"errors"
	"sync"
)

type InMemoryRepository struct {
	feeds []Feed
	sync.Mutex
}

func NewMemoryRepository() *InMemoryRepository {
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

	return &InMemoryRepository{
		feeds: feeds,
	}
}

func (imr *InMemoryRepository) GetFeeds() ([]Feed, error) {
	return imr.feeds, nil
}

func (imr *InMemoryRepository) GetFeed(feedName string) (Feed, error) {
	for _, feed := range imr.feeds {
		if feed.FeedName == feedName {
			return feed, nil
		}
	}
	return Feed{}, errors.New("no such feed exists")
}

func (imr *InMemoryRepository) AddFeed(feed Feed) (Feed, error) {
	imr.feeds = append(imr.feeds, feed)
	return Feed{}, errors.New("something went wrong adding feed")
}