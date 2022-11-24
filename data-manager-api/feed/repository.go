package feed

type Repository interface {
	GetFeeds() ([]Feed, error)
	GetFeed(id int) (Feed, error)
	UpdateFeed(feed Feed) (Feed, error)
	AddFeed(feed Feed) (Feed, error)
	DeleteFeed(id int) (string, error)
}