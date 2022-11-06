package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/scoyne2/data-manager-api/graph/generated"
	"github.com/scoyne2/data-manager-api/graph/model"
)

// CreateFeed is the resolver for the createFeed field.
func (r *mutationResolver) CreateFeed(ctx context.Context, input model.NewFeed) (*model.Feed, error) {
	feed := &model.Feed{
		ID:       fmt.Sprintf("T%d", rand.Int()),
		Vendor:   input.Vendor,
		FeedName: input.FeedName,
		FeedType: input.FeedType,
	}
	r.feeds = append(r.feeds, feed)
	return feed, nil
}

// CreateFileStats is the resolver for the createFileStats field.
func (r *mutationResolver) CreateFileStats(ctx context.Context, input model.NewFileStats) (*model.FileStats, error) {
	fileStat := &model.FileStats{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Name:        input.Name,
		Date:        input.Date,
		Vendor:      input.Vendor,
		Errors:      input.Errors,
		Status:      input.Status,
		Designation: input.Designation,
	}
	r.fileStats = append(r.fileStats, fileStat)
	return fileStat, nil
}

// Feed is the resolver for the feed field.
func (r *queryResolver) Feed(ctx context.Context, id string) (*model.Feed, error) {
	for _, feed := range r.feeds {
		if id == feed.ID {
			return feed, nil
		}
	}
	return nil, nil
}

// Feeds is the resolver for the feeds field.
func (r *queryResolver) Feeds(ctx context.Context) ([]*model.Feed, error) {
	return r.feeds, nil
}

// FileStat is the resolver for the fileStat field.
func (r *queryResolver) FileStat(ctx context.Context, id string) (*model.FileStats, error) {
	for _, fileStat := range r.fileStats {
		if id == fileStat.ID {
			return fileStat, nil
		}
	}
	return nil, nil
}

// FileStats is the resolver for the fileStats field.
func (r *queryResolver) FileStats(ctx context.Context) ([]*model.FileStats, error) {
	return r.fileStats, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
