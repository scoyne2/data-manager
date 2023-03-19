package schemas

import (
	"testing"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/scoyne2/data-manager-api/feed"
)

type mockFeedService struct{}

func (m *mockFeedService) AddFeed(vendor string, feedName string, feedMethod string) (int, error) {
	return 123, nil
}

func (m *mockFeedService) UpdateFeed(vendor string, feedName string, feedMethod string) error {
	return nil
}

func (m *mockFeedService) DeleteFeed(id int) error {
	return nil
}

func TestGenerateRootMutation(t *testing.T) {
	fs := &feed.FeedService{}
	mutation := generateRootMutation(fs)

	assert.NotNil(t, mutation)

	// Test addFeed field
	addFeedField := mutation.Fields()["addFeed"]
	assert.NotNil(t, addFeedField)
	assert.NotNil(t, addFeedField.Args)

	assert.Equal(t, graphql.NewNonNull(graphql.String), addFeedField.Args[0].Type)
	assert.Equal(t, graphql.NewNonNull(graphql.String), addFeedField.Args[1].Type)
	assert.Equal(t, graphql.NewNonNull(graphql.String), addFeedField.Args[2].Type)

	// Test updateFeed field
	updateFeedField := mutation.Fields()["updateFeed"]
	assert.NotNil(t, updateFeedField)
	assert.NotNil(t, updateFeedField.Args)
	assert.Equal(t, graphql.NewNonNull(graphql.String), updateFeedField.Args[0].Type)
	assert.Equal(t, graphql.NewNonNull(graphql.String), updateFeedField.Args[1].Type)
	assert.Equal(t, graphql.NewNonNull(graphql.String), updateFeedField.Args[2].Type)

	// Test deleteFeed field
	deleteFeedField := mutation.Fields()["deleteFeed"]
	assert.NotNil(t, deleteFeedField)
	assert.NotNil(t, deleteFeedField.Args)
	assert.Equal(t, graphql.NewNonNull(graphql.Int), deleteFeedField.Args[0].Type)
}
