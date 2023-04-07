package socialmedia

import (
	"os"
)

func Start() {
	client_id := os.Getenv("GOOGLE_DATA_CLIENT_ID")
	if client_id == "" {
		panic("GOOGLE_DATA_CLIENT_ID not set")
	}

	// youtube.NewService(context.Background(), option.WithCredentials(&google.Credentials{
	// TokenSource: ,
	// }))
}
