package lib

import (
	"crypto/rand"
	"errors"
	"fmt"
	"go-news/consts"
	"go-news/database"
	"go-news/session"
	"go-news/types"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func FetchAuthData(r *http.Request) types.AuthData {
	user, err := FetchSession(r)

	if err != nil {
		return types.AuthData{
			IsLoggedIn: false,
			IsAdmin:    false,
		}
	}

	return types.AuthData{
		IsLoggedIn: true,
		IsAdmin:    user.Role == "ADMIN",
	}
}

func FetchSession(r *http.Request) (types.User, error) {
	// if authMiddleware already called, on calling this fun reading from context
	// would avoid re-query user from db

	userPointer, ok := r.Context().Value(consts.AUTH_COOKIE).(*types.User)

	if ok && userPointer != nil {
		return *userPointer, nil
	}

	session, _ := session.Store.Get(r, consts.AUTH_COOKIE)
	user_id, ok := session.Values[consts.AUTH_COOKIE].(int64)

	if user_id == -1 || !ok {
		return types.User{}, errors.New("not logged in")
	}

	var user types.User

	err := database.GetDB().QueryRow("select user_id,email,name,role from users where user_id=?", user_id).Scan(&user.UserId, &user.Email, &user.Name, &user.Role)

	if err != nil {
		return types.User{}, errors.New("not logged in")
	}

	return user, nil
}

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

func GetAllMyPostsUtil(user_id int64) []types.Post {
	var posts []types.Post

	rows, err := database.GetDB().Query(`
		select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
		where user_id=?
	`, user_id)

	if err != nil {
		return posts
	}

	defer rows.Close()

	for rows.Next() {
		var post types.Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	return posts
}
func GetAllPostsUtil() []types.Post {
	var posts []types.Post
	rows, err := database.GetDB().Query(`
		select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
	`)
	if err != nil {
		fmt.Println("er: ", err)
		return posts
	}
	if err != nil {
		return posts
	}
	defer rows.Close()

	for rows.Next() {
		var post types.Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	return posts
}

func GetAllPostsExceptUtil(user_id int64) []types.Post {
	var posts []types.Post

	rows, err := database.GetDB().Query(`
		select p.post_id, p.title, p.slug, p.body, p.author_id, u.name
		from posts p
		inner join users u
		on p.author_id=u.user_id
		where not p.author_id=? 
	`, user_id)

	if err != nil {
		return posts
	}

	defer rows.Close()

	for rows.Next() {
		var post types.Post
		_ = rows.Scan(&post.PostId, &post.Title, &post.Slug, &post.Body, &post.AuthorId, &post.Author)

		posts = append(posts, post)
	}

	return posts
}

func IsAdmin(r *http.Request) bool {
	user, ok := r.Context().Value(consts.AUTH_COOKIE).(*types.User)
	if user == nil || !ok || user.Role != "ADMIN" {
		return false
	}
	return true
}
