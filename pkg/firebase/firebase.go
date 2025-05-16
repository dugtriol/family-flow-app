package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

const (
	serviceAccountKeyPath = "family-flow-app-4900a-firebase-adminsdk-fbsvc-4aa2ce52fc.json"
)

func Init(ctx context.Context) *firebase.App {
	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}

	// Проверка FCM клиента
	// client, err := app.Messaging(ctx)
	// if err != nil {
	// 	log.Fatalf("error initializing Firebase Messaging client: %v", err)
	// }
	// log.Println("Firebase Messaging client initialized successfully")

	// // Пример тестового вызова FCM
	// _, err = client.Send(
	// 	ctx, &messaging.Message{
	// 		Notification: &messaging.Notification{
	// 			Title: "Test Notification - 1",
	// 			Body:  "This is a test notification to verify Firebase setup.",
	// 		},
	// 		Token: "cIcGYigbTdaNSebLqCvlva:APA91bGo_CrN7n3zojl5ZZwAKi7tMhYHqBEa39HN21AGDvBpH5RAediM6rqPx_8VdRuOaQrhDEE5un6u_aA5Xe8ztAqxZLXZncgEfm5RQYgC0aMY1DHJGoU", // Замените на реальный токен устройства для теста
	// 	},
	// )
	// if err != nil {
	// 	log.Printf("Test FCM notification failed (this is expected if token is invalid): %v", err)
	// } else {
	// 	log.Println("Test FCM notification sent successfully")
	// }

	return app
}
