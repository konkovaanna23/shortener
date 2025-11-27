package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const cookieName = "user"
const userIDContextKey = "userID"

func GetUserID(r *http.Request) (string, bool) {
	v := r.Context().Value(userIDContextKey)
	id, ok := v.(string)
	return id, ok
}

func sign(key string, loadString string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(loadString))
	sum := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum)
}

func verify(key, loadString, sig string) bool {
	expected := sign(key, loadString)
	return hmac.Equal([]byte(expected), []byte(sig))
}

func GenerateUserID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func parseCookieValue(v string) (userID, sig string, err error) {
	parts := strings.Split(v, ".")
	if len(parts) != 2 {
		return "", "", errors.New("некорректный формат cookie")
	}
	return parts[0], parts[1], nil
}

func AuthMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var userID string
			var needSetCookie bool

			c, err := r.Cookie(cookieName)
			fmt.Printf("%v: \n", r.Cookies())
			if err == nil {
				id, sig, err := parseCookieValue(c.Value)
				if err != nil || id == "" {
					http.Error(w,
						http.StatusText(http.StatusUnauthorized),
						http.StatusUnauthorized,
					)
					return
				}

				if verify(key, id, sig) {
					userID = id
				} else {
					needSetCookie = true
				}
			} else {
				needSetCookie = true
			}

			if needSetCookie {
				id, err := GenerateUserID()
				if err != nil {
					http.Error(w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
					return
				}
				userID = id
				sig := sign(key, userID)
				value := userID + "." + sig

				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    value,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
				})
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
