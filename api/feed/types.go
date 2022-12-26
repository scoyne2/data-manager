package feed

type Feed struct {
	ID         	int    `json:"id"`
	Vendor     	string `json:"vendor"`
	FeedName   	string `json:"feed_name"`
	FeedMethod	string `json:"feed_method"`
}

type FeedStatus struct {
	ID         	int    `json:"id"`
	ProcessDate string `json:"process_date"`
	RecordCount int    `json:"record_count"`
	ErrorCount	int    `json:"error_count"`
	Status	    string `json:"feed_status"`
}

type FeedStatusResults struct {
	ID         	int    `json:"id"`
	ProcessDate string `json:"process_date"`
	RecordCount int    `json:"record_count"`
	ErrorCount	int    `json:"error_count"`
	Status	    string `json:"feed_status"`
	Vendor     	string `json:"vendor"`
	FeedName   	string `json:"feed_name"`
	FeedMethod	string `json:"feed_method"`
}

type FeedStatusAggregate struct {
	Files     	int `json:"files"`
	Rows   		int `json:"rows"`
	Errors		int `json:"errors"`
}

type Repository interface {
	GetFeeds() ([]Feed, error)
	GetFeed(id int) (Feed, error)
	UpdateFeed(feed Feed) (Feed, error)
	AddFeed(feed Feed) (string, error)
	DeleteFeed(id int) (string, error)

	GetFeedStatuses() ([]FeedStatusResults, error)
	GetFeedStatusesAggregate(startDate string, endDate string) (FeedStatusAggregate, error)
}
