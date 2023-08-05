package schemas

import (
	"github.com/graphql-go/graphql"
	"github.com/scoyne2/data-manager-api/feed"
)

var feedStatusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeedStatus",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "The ID that is used to identify unique feed status",
		},
		"process_date": &graphql.Field{
			Type:        graphql.String,
			Description: "The date the feed was processed",
		},
		"record_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of records successfully processed",
		},
		"error_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of records that had errors",
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "The status of the feed",
		},
		"vendor": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the vendor",
		},
		"feed_name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the feed",
		},
		"feed_method": &graphql.Field{
			Type:        graphql.String,
			Description: "The method by which the feed is received",
		},
		"file_name": &graphql.Field{
			Type:        graphql.String,
			Description: "The file received",
		},
	},
},
)

var feedStatusDetailsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeedStatusDetails",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "The ID that is used to identify unique feed status",
		},
		"process_date": &graphql.Field{
			Type:        graphql.String,
			Description: "The date the feed was processed",
		},
		"record_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of records successfully processed",
		},
		"error_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of records that had errors",
		},
		"status": &graphql.Field{
			Type:        graphql.String,
			Description: "The status of the feed",
		},
		"vendor": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the vendor",
		},
		"feed_name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the feed",
		},
		"feed_method": &graphql.Field{
			Type:        graphql.String,
			Description: "The method by which the feed is received",
		},
		"previous_feeds": &graphql.Field{
			Type:        graphql.NewList(feedStatusType),
			Description: "List of previous feed events",
		},
	},
},
)

var feedStatusAggregateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeedStatusAggregate",
	Fields: graphql.Fields{
		"files": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of files processed",
		},
		"rows": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of records processed",
		},
		"errors": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of records that had errors",
		},
	},
},
)

var dataPreviewType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DataPreview",
	Fields: graphql.Fields{
		"columns": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "A List of column names from the file",
		},
		"values": &graphql.Field{
			Type:        graphql.NewList(graphql.NewList(graphql.String)),
			Description: "A List of Lists which contain the values from the file",
		},
	},
},
)
// GenerateSchema will create a GraphQL Schema and set the Resolvers found in the FeedService
// For all the needed fields
func GenerateSchema(fs *feed.FeedService) (*graphql.Schema, error) {

	// RootQuery
	fields := graphql.Fields{
		"feedstatusesdetailed": &graphql.Field{
			Type: graphql.NewList(feedStatusDetailsType),
			Resolve: fs.ResolveGetFeedStatusDetails,
			Description: "Query all Feed Statuses With Details",
		},
		"feedstatuseaggregates": &graphql.Field{
			Type: feedStatusAggregateType,
			Resolve: fs.ResolveFeedStatuseAggregate,
			Args: graphql.FieldConfigArgument{
				"startDate": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"endDate": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "Aggregate of Feed Statuses",
		},
		"datapreview": &graphql.Field{
			Type: dataPreviewType,
			Resolve: fs.ResolveDataPreview,
			Args: graphql.FieldConfigArgument{
				"vendor": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"feedName": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"fileName": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"s3Bucket": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "Data Preview of a given file",
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