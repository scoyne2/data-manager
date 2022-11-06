// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type FileStats struct {
	ID          string `json:"id"`
	Rows        int    `json:"rows"`
	Name        string `json:"name"`
	Date        string `json:"date"`
	Vendor      string `json:"vendor"`
	Errors      int    `json:"errors"`
	Status      string `json:"status"`
	Designation string `json:"designation"`
}

type NewFeed struct {
	Vendor   string `json:"vendor"`
	FeedName string `json:"feedName"`
	FeedType string `json:"feedType"`
}

type NewFileStats struct {
	Rows        int    `json:"rows"`
	Name        string `json:"name"`
	Date        string `json:"date"`
	Vendor      string `json:"vendor"`
	Errors      int    `json:"errors"`
	Status      string `json:"status"`
	Designation string `json:"designation"`
}
