package main
import (	
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strings"
	"crypto/rand"
	"fmt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func Slugify(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, " ", "-")

	re := regexp.MustCompile(`[^a-z0-9-]+`)
	s = re.ReplaceAllString(s, "")

	s = strings.Trim(s, "-")

	randomString, err := generateRandomString(5)
	if err != nil {
		return s
	}

	slug := fmt.Sprintf("%s-%s", s, randomString)

	return slug
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result.WriteByte(letters[num.Int64()])
	}

	return result.String(), nil
}
