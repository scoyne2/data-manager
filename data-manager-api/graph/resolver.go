package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/scoyne2/data-manager-api/graph/model"
)

type Resolver struct{
	feeds []*model.Feed
	feed *model.Feed
	fileStats []*model.FileStats
	fileStat *model.FileStats
}