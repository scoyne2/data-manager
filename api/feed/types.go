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
	FileName   	string `json:"file_name"`
	FeedID      int    `json:"feed_id"`
}

type FeedStatusUpdate struct {
	ProcessDate string `json:"process_date"`
	RecordCount int    `json:"record_count"`
	ErrorCount	int    `json:"error_count"`
	Status	    string `json:"feed_status"`
	FileName   	string `json:"file_name"`
	Vendor     	string `json:"vendor"`
	FeedName   	string `json:"feed_name"`
}

type FeedStatusResults struct {
	ID         	int    `json:"id" db:"id"`
	ProcessDate string `json:"process_date" db:"process_date"`
	RecordCount int    `json:"record_count" db:"record_count"`
	ErrorCount	int    `json:"error_count" db:"error_count"`
	Status	    string `json:"feed_status" db:"feed_status"`
	Vendor     	string `json:"vendor" db:"vendor"`
	FeedName   	string `json:"feed_name" db:"feed_name"`
	FeedMethod	string `json:"feed_method" db:"feed_method"`
	FileName   	string `json:"file_name" db:"file_name"`
}

type FeedStatusResultsDetailed struct {
	ID         	  int                 `json:"id" db:"id"`
	ProcessDate   string              `json:"process_date" db:"process_date"`
	RecordCount   int                 `json:"record_count" db:"record_count"`
	ErrorCount	  int                 `json:"error_count" db:"error_count"`
	Status	      string              `json:"feed_status" db:"feed_status"`
	Vendor     	  string              `json:"vendor" db:"vendor"`
	FeedName   	  string              `json:"feed_name" db:"feed_name"`
	FeedMethod	  string              `json:"feed_method" db:"feed_method"`
	FileName   	  string              `json:"file_name" db:"file_name"`
	PreviousFeeds []FeedStatusResults `json:"previous_feeds" db:"previous_feeds"`
}

type FeedStatusAggregate struct {
	Files     	int `json:"files"`
	Rows   		int `json:"rows"`
	Errors		int `json:"errors"`
}

type Repository interface {
	AddFeed(feed Feed) (string, error)
	UpdateFeedStatus(feedStatusUpdate FeedStatusUpdate) (string, error)
	GetFeedStatusDetails() ([]FeedStatusResultsDetailed, error)
	GetFeedStatusesAggregate(startDate string, endDate string) (FeedStatusAggregate, error)
}
