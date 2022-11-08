package feed

type Repository interface {
	GetFeeds() ([]Feed, error)
	GetFeed(feedName string) (Feed, error)
	AddFeed(feed Feed) (Feed, error)
}