package apiserver

import (
	"context"
	"net/http"
	"notification-service/internal/app/store/mongostore"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Start(config *Config) error {
	clientOptions := options.Client().ApplyURI(config.DatabaseURL)
	ctx := context.Background()
	defer ctx.Done()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	database := client.Database("restapi_dev")
	collection := database.Collection("inventory")

	store := mongostore.New(ctx, collection)
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}
