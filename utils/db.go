package utils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDbConnection is function to create New DB Connection
func NewDbConnection(host, port, user, password, authSrc string) (*mongo.Client, *context.Context, error) {

	mongoDBConnectionURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", user, password, host, port, authSrc)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnectionURI))

	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	err = client.Connect(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, &ctx, nil
}
