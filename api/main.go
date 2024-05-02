package main

import (
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	config "kanban-go/config"

	"context"
)

func main() {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON([]byte(config.FB_SECRET_CREDENTIAL())))
	if err != nil {
		panic("error initializing app: " + err.Error())
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}

	defer client.Close()
}
