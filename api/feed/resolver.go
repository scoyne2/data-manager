package feed

import (
	"github.com/graphql-go/graphql"
)

type Resolver interface {
	ResolveGetFeedStatusDetails(p graphql.ResolveParams) (interface{}, error)
	ResolveFeedStatuseAggregate(p graphql.ResolveParams) (interface{}, error)
	ResolveDataPreview(p graphql.ResolveParams) (interface{}, error)
}

type FeedService struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface,) FeedService {
	return FeedService{
		repo: repo,
	}
}

func (fs FeedService) UpdateFeedStatus(p graphql.ResolveParams) (interface{}, error) {
	processDate := p.Args["processDate"].(string)
	recordCount := p.Args["recordCount"].(int)
	errorCount := p.Args["errorCount"].(int)
	status := p.Args["status"].(string)
	fileName := p.Args["fileName"].(string)
	vendor := p.Args["vendor"].(string)
	feedName := p.Args["feedName"].(string)


	var feedStatusUpdate FeedStatusUpdate
	feedStatusUpdate.ProcessDate = processDate
	feedStatusUpdate.RecordCount = recordCount
	feedStatusUpdate.ErrorCount = errorCount
	feedStatusUpdate.Status = status
	feedStatusUpdate.FileName = fileName
	feedStatusUpdate.Vendor = vendor
	feedStatusUpdate.FeedName = feedName

	return fs.repo.UpdateFeedStatus(feedStatusUpdate)
}


func (fs FeedService) AddFeed(p graphql.ResolveParams) (interface{}, error) {
	vendor := p.Args["vendor"].(string)
	feedName := p.Args["feedName"].(string)
	feedMethod := p.Args["feedMethod"].(string)

	var feed Feed
	feed.Vendor = vendor
	feed.FeedName = feedName
	feed.FeedMethod = feedMethod

	return fs.repo.AddFeed(feed)
}


func (fs FeedService) ResolveGetFeedStatusDetails(p graphql.ResolveParams) (interface{}, error) {
	fsDetails, err := fs.repo.GetFeedStatusDetails()
	if err != nil {
		return nil, err
	}
	return fsDetails, nil
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

func (fs FeedService) ResolveDataPreview(p graphql.ResolveParams) (interface{}, error) {
	vendor := p.Args["vendor"].(string)
	feedName := p.Args["feedName"].(string)
	fileName := p.Args["fileName"].(string)
	s3Bucket := p.Args["s3Bucket"].(string)
	dataPreview, err := fs.repo.GetDataPreview(vendor, feedName, fileName, s3Bucket)
	if err != nil {
		return nil, err
	}
	return dataPreview, nil
}