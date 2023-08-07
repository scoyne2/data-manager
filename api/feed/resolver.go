package feed

import (
	"github.com/graphql-go/graphql"
)

type Resolver interface {
	ResolveGetFeedStatusDetails(p graphql.ResolveParams) (interface{}, error)
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

func (fs FeedService) UpdateFeedStatus(p graphql.ResolveParams) (interface{}, error) {
	processDate := p.Args["processDate"].(string)
	recordCount := p.Args["recordCount"].(int)
	errorCount := p.Args["errorCount"].(int)
	status := p.Args["status"].(string)
	fileName := p.Args["fileName"].(string)
	vendor := p.Args["vendor"].(string)
	feedName := p.Args["feedName"].(string)

	var emrApplicationID string
	if p.Args["emrApplicationID"] != nil {
		emrApplicationID = p.Args["emrApplicationID"].(string)
	} 
	
	var emrStepID string
	if p.Args["emrStepID"] != nil {
		emrStepID = p.Args["emrStepID"].(string)
	}

	var feedStatusUpdate FeedStatusUpdate
	feedStatusUpdate.ProcessDate = processDate
	feedStatusUpdate.RecordCount = recordCount
	feedStatusUpdate.ErrorCount = errorCount
	feedStatusUpdate.Status = status
	feedStatusUpdate.FileName = fileName
	feedStatusUpdate.Vendor = vendor
	feedStatusUpdate.FeedName = feedName
	feedStatusUpdate.EMRApplicationID = emrApplicationID
	feedStatusUpdate.EMRStepID = emrStepID

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