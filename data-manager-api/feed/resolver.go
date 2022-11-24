package feed

import (
	"errors"

	"github.com/graphql-go/graphql"
)

type Resolver interface {
	ResolveFeed(p graphql.ResolveParams) (interface{}, error)
	ResolveFeeds(p graphql.ResolveParams) (interface{}, error)
}

type FeedService struct {
	feeds Repository
}

func NewService(repo Repository,) FeedService {
	return FeedService{
		feeds: repo,
	}
}

func (fs FeedService) ResolveFeeds(p graphql.ResolveParams) (interface{}, error) {
	feeds, err := fs.feeds.GetFeeds()
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (fs FeedService) ResolveFeed(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("id has to be an int")
	}
	feed, err := fs.feeds.GetFeed(id)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func (fs FeedService) UpdateFeed(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	vendor := p.Args["vendor"].(string)
	feedName := p.Args["feedName"].(string)
	feedMethod := p.Args["feedMethod"].(string)

	var feed Feed
	feed.ID = id
	feed.Vendor = vendor
	feed.FeedName = feedName
	feed.FeedMethod = feedMethod

	return fs.feeds.UpdateFeed(feed)
}



func (fs FeedService) AddFeed(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	vendor := p.Args["vendor"].(string)
	feedName := p.Args["feedName"].(string)
	feedMethod := p.Args["feedMethod"].(string)

	var feed Feed
	feed.ID = id
	feed.Vendor = vendor
	feed.FeedName = feedName
	feed.FeedMethod = feedMethod

	return fs.feeds.AddFeed(feed)
}

