package main

import (
	"testing"
)

func TestValidSignature(t *testing.T) {
	body := []byte(`{"action":"completed"}`)
	secret := "test-secret"
	if validSignature(body, "", secret) { t.Fatal("empty signature accepted") }
}
