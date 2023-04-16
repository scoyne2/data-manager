package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/scoyne2/data-manager-api/feed"
)

var addFeedArgs = graphql.FieldConfigArgument{
	"vendor": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"feedName": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"feedMethod": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
}

var deleteFeedArgs = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.Int),
	},
}


var updateFeedStatusArgs = graphql.FieldConfigArgument{
	"processDate": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"recordCount": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.Int),
	},
	"errorCount": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.Int),
	},
	"status": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"fileName": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"vendor": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
	"feedName": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.String),
	},
}

func generateRootMutation(fs *feed.FeedService) *graphql.Object {

	mutationFields := graphql.Fields{
		"addFeed": generateGraphQLField(graphql.String, fs.AddFeed, "Add a new feed", addFeedArgs),
		"updateFeed": generateGraphQLField(feedType, fs.UpdateFeed, "Update an existing feed", addFeedArgs),
		"deleteFeed": generateGraphQLField(graphql.String, fs.DeleteFeed, "Delete existing feed", deleteFeedArgs),
		"updateFeedStatus": generateGraphQLField(graphql.String, fs.UpdateFeedStatus, "Update an existing feed status", updateFeedStatusArgs),
	}
	mutationConfig := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}

	return graphql.NewObject(mutationConfig)
}