package auth

import (
	"encoding/base64"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

func argon2_password(password string, salt string) string {
	hashed_pw := argon2.Key([]byte(password), []byte(salt), 3, 32*1024, 1, 32)
	encoded := base64.StdEncoding.EncodeToString(hashed_pw)
	return encoded
}

func GetHashedPassword(password string) (salt string, hashed_pw string) {
	generated_salt := uuid.New().String()
	hashed_pw = argon2_password(password, generated_salt)
	return generated_salt, hashed_pw
}
