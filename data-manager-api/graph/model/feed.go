package model

type Feed struct {
	ID         string `json:"id"`
	Vendor     string `json:"vendor"`
	FeedName   string   `json:"feedName"`
	FeedType   string `json:"feedType"`
}