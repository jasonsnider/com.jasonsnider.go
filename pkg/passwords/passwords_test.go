package passwords

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if len(hash) == 0 {
		t.Fatalf("HashPassword returned an empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatalf("CheckPasswordHash returned false for a valid password and hash")
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Fatalf("CheckPasswordHash returned true for an invalid password and hash")
	}
}
