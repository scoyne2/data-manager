package feed

import (
	"errors"

	"github.com/graphql-go/graphql"
)

type Resolver interface {
	ResolveFeed(p graphql.ResolveParams) (interface{}, error)
	ResolveFeeds(p graphql.ResolveParams) (interface{}, error)
	ResolveFeedStatuses(p graphql.ResolveParams) (interface{}, error)
	ResolveFeedStatuseAggregate(p graphql.ResolveParams) (interface{}, error)
}

type FeedService struct {
	repo Repository
}

func NewService(repo Repository,) FeedService {
	return FeedService{
		repo: repo,
	}
}

func (fs FeedService) ResolveFeeds(p graphql.ResolveParams) (interface{}, error) {
	feeds, err := fs.repo.GetFeeds()
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
	feed, err := fs.repo.GetFeed(id)
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

	return fs.repo.UpdateFeed(feed)
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

	return fs.repo.AddFeed(feed)
}

func (fs FeedService) DeleteFeed(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	return fs.repo.DeleteFeed(id)
}

func (fs FeedService) ResolveFeedStatuses(p graphql.ResolveParams) (interface{}, error) {
	feedstatuses, err := fs.repo.GetFeedStatuses()
	if err != nil {
		return nil, err
	}
	return feedstatuses, nil
}

func (fs FeedService) ResolveFeedStatuseAggregate(p graphql.ResolveParams) (interface{}, error) {
	startDate := p.Args["startDate"].(string)
	endDate := p.Args["endDate"].(string)
	fsAgg, err := fs.repo.GetFeedStatusesAggregate(startDate, endDate)
	if err != nil {
		return nil, err
	}
	return fsAgg, nil
}