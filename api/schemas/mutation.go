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
	"emrApplicationID": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"emrStepID": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

func generateRootMutation(fs *feed.FeedService) *graphql.Object {

	mutationFields := graphql.Fields{
		"addFeed": generateGraphQLField(graphql.String, fs.AddFeed, "Add a new feed", addFeedArgs),
		"updateFeedStatus": generateGraphQLField(graphql.String, fs.UpdateFeedStatus, "Update an existing feed status", updateFeedStatusArgs),
	}
	mutationConfig := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}

	return graphql.NewObject(mutationConfig)
}