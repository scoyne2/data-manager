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

func (m *mockFeedService) UpdateFeedStatus(processDate string, recordCount int, errorCount int, status string, fileName string, vendor string, feedName string) (int, error) {
	return 123, nil
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

	// Test updateFeedStatus field
	updateFeedStatusField := mutation.Fields()["updateFeedStatus"]
	assert.NotNil(t, updateFeedStatusField)
	assert.NotNil(t, updateFeedStatusField.Args)
}