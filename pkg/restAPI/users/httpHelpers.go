package users

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const CookiePlayerIdKey = "PlayerId"

func AuthMiddleware(noAuth bool) (isAllowed func(res http.ResponseWriter, req *http.Request) bool) {
	return func(res http.ResponseWriter, req *http.Request) bool {
		value := GetUserFromReq(req)
		if (value == nil) != noAuth {
			if !noAuth {
				c := &http.Cookie{
					Name:     CookiePlayerIdKey,
					Path:     "/",
					Value:    "",
					Expires:  time.Unix(0, 0),
					HttpOnly: true,
					Secure:   true,
				}
				http.SetCookie(res, c)
			}
			return false
		}
		return true
	}
}

func GetUserFromReq(req *http.Request) *AuthUser {
	for _, cookie := range req.Cookies() {
		if cookie.Name == CookiePlayerIdKey {
			cookieval, ok := UserMap[cookie.Value]
			if !ok {
				continue
			}

			cookieval.lastSeen = time.Now()

			return cookieval
		}
	}

	return nil
}

func createCookie() UserCookie {
	return uuid.New().String()
}
