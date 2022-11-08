package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/scoyne2/data-manager-api/feed"
)

var addFeedArgs = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.Int),
	},
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

func generateRootMutation(fs *feed.FeedService) *graphql.Object {

	mutationFields := graphql.Fields{
		"addFeed": generateGraphQLField(feedType, fs.AddFeed, "Add a new feed", addFeedArgs),
	}
	mutationConfig := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}

	return graphql.NewObject(mutationConfig)
}