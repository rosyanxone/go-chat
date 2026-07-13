package main

import (
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// Run once: `go run ./cmd/vapidkeygen`
// Paste the printed values to .env as
// VAPID_PUBLIC_KEY / VAPID_PRIVATE_KEY.
func main() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()

	if err != nil {
		log.Fatalf("failed to generate VAPID keys: %v", err)
	}

	fmt.Println("VAPID_PUBLIC_KEY=" + publicKey)
	fmt.Println("VAPID_PRIVATE_KEY=" + privateKey)
}
