package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

const (
	serviceAccountKeyPath = "family-flow-app-4900a-firebase-adminsdk-fbsvc-9e2d528ca1.json"
)

func Init(ctx context.Context) *firebase.App {
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	return app
}
