package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/scoyne2/data-manager-api/feed"
)

var feedType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Feed",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "The ID that is used to identify unique feeds",
		},
		"vendor": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the vendor",
		},
		"feedName": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the feed",
		},
		"feedMethod": &graphql.Field{
			Type:        graphql.String,
			Description: "The method by which the feed is received",
		},
	},
},
)

// GenerateSchema will create a GraphQL Schema and set the Resolvers found in the GopherService
// For all the needed fields
func GenerateSchema(fs *feed.FeedService) (*graphql.Schema, error) {

	// RootQuery
	fields := graphql.Fields{
		// We define the Feeds query
		"feeds": &graphql.Field{
			// It will return a list of feedType, a List is an Slice
			Type: graphql.NewList(feedType),
			// We change the Resolver to use the feedRepo instead, allowing us to access all Feeds
			Resolve: fs.ResolveFeeds,
			// Description explains the field
			Description: "Query all Feeds",
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// RootMutation
	rootMutation := generateRootMutation(fs)

	// Now combine all Objects into a Schema Configuration
	schemaConfig := graphql.SchemaConfig{
		// Query is the root object query schema
		Query: graphql.NewObject(rootQuery),
		// Appliy the Mutation to the schema
		Mutation: rootMutation,
	}
	// Create a new GraphQL Schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

// generateGraphQLField is a generic builder factory to create graphql fields
func generateGraphQLField(output graphql.Output, resolver graphql.FieldResolveFn, description string, args graphql.FieldConfigArgument) *graphql.Field {
	return &graphql.Field{
		Type:        output,
		Resolve:     resolver,
		Description: description,
		Args:        args,
	}
}