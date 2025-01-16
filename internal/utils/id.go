package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

var alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// NewID generates a new ID using the default alphabet and length (21 chars)
func NewID() string {
	return gonanoid.Must()
}

// NewIDShort uses a custom alphabet (no "_-") to generate a shorter ID of lentgh 10
func NewIDShort() string {
	return gonanoid.MustGenerate(alphabet, 10)
}

