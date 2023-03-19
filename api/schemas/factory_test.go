package schemas_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/scoyne2/data-manager-api/feed"
	"github.com/scoyne2/data-manager-api/schemas"
)

type mockFeedService struct {}

func (m *mockFeedService) ResolveFeeds(p graphql.ResolveParams) (interface{}, error) {
	// Mock implementation for ResolveFeeds
	return nil, nil
}

func (m *mockFeedService) ResolveFeed(p graphql.ResolveParams) (interface{}, error) {
	// Mock implementation for ResolveFeed
	return nil, nil
}

func (m *mockFeedService) ResolveFeedStatuses(p graphql.ResolveParams) (interface{}, error) {
	// Mock implementation for ResolveFeedStatuses
	return nil, nil
}

func (m *mockFeedService) ResolveFeedStatuseAggregate(p graphql.ResolveParams) (interface{}, error) {
	// Mock implementation for ResolveFeedStatuseAggregate
	return nil, nil
}

func TestGenerateSchema(t *testing.T) {
	// Setup
	fs := &feed.FeedService{}

	// Execution
	schema, err := schemas.GenerateSchema(fs)

	// Verification
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if schema == nil {
		t.Error("expected a schema, got nil")
	}
}