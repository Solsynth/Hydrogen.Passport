package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) string {
	data, _ := bcrypt.GenerateFromPassword([]byte(raw), 12)
	return string(data)
}

func VerifyPassword(text string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(text)) == nil
}
