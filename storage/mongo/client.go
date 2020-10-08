package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClientFromURL(url string) (*mongo.Client, error) {
	return mongo.NewClient(options.Client().ApplyURI(url))
}
