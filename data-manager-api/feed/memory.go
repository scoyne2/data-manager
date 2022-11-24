package feed

import (
	"errors"
	"fmt"
	"sync"
)

type InMemoryRepository struct {
	feeds map[int]Feed
	sync.Mutex
}

func NewMemoryRepository() *InMemoryRepository {
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

	return &InMemoryRepository{
		feeds: feeds,
	}
}

func (imr *InMemoryRepository) GetFeeds() ([]Feed, error) {
	var feeds []Feed
	for _, feed := range imr.feeds {
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (imr *InMemoryRepository) GetFeed(id int) (Feed, error) {
	feed, ok := imr.feeds[id]
	if ok {
		return feed, nil
	}
	return Feed{}, fmt.Errorf("Feed Id %x does not exist", id)
}

func (imr *InMemoryRepository) UpdateFeed(feed Feed) (Feed, error) {
	_, ok := imr.feeds[feed.ID]
	if ok {
		imr.feeds[feed.ID] = feed
		return feed, nil
	}
	return Feed{}, errors.New("no such feed exists")
}

func (imr *InMemoryRepository) AddFeed(feed Feed) (Feed, error) {
	_, ok := imr.feeds[feed.ID]
	if ok {
		return Feed{}, fmt.Errorf("Feed Id %x already exists", feed.ID)
	}
	imr.feeds[feed.ID] = feed
	return feed, nil
}

func (imr *InMemoryRepository) DeleteFeed(id int) (string, error) {
	_, ok := imr.feeds[id]
	if ok {
		delete(imr.feeds, id)
		return fmt.Sprintf("Feed Id %x deleted", id), nil
	}
	return "", fmt.Errorf("Feed Id %x does not exist", id)
}