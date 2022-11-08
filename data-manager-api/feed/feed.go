package feed

type Feed struct {
	ID         	int `json:"id"`
	Vendor     	string `json:"vendor"`
	FeedName   	string `json:"feedName"`
	FeedMethod	string `json:"feedMethod"`
}