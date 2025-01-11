package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kraanter/blackjack/pkg/manager"
)

type ContextPlayer string

const ContextUserKey ContextPlayer = "user"
const CookiePlayerIdKey = "PlayerId"

func AuthMiddleware(noAuth bool) func(res http.ResponseWriter, req *http.Request) bool {
	return func(res http.ResponseWriter, req *http.Request) bool {
		value := GetUserFromReq(req)
		fmt.Printf("mid value %v\n", value)
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
			http.Error(res, "Unauthorized: Invalid User", http.StatusUnauthorized)
			return true
		}
		return false
	}
}

type UserCookie = string

type AuthUser struct {
	Cookie        UserCookie             `json:"cookie"`
	Player        *manager.ManagedPlayer `json:"player"`
	Ctx           context.Context        `json:"-"`
	cancelContext context.CancelFunc
}

var UserMap = make(map[UserCookie]*AuthUser)

func GetUserFromReq(req *http.Request) *AuthUser {
	for _, cookie := range req.Cookies() {
		if cookie.Name == CookiePlayerIdKey {
			cookieval, ok := UserMap[cookie.Value]
			if !ok {
				continue
			}

			return cookieval
		}
	}

	return nil
}

func getUserFromContext(ctx context.Context) *AuthUser {
	value := ctx.Value(ContextUserKey)
	fmt.Printf("contextval %v  \n", value)
	switch value.(type) {
	case *AuthUser:
		return value.(*AuthUser)
	default:
		return nil
	}
}

func RegisterUser(player *manager.ManagedPlayer, ctx context.Context) UserCookie {
	authUser := createAuthUser(player, ctx)

	UserMap[authUser.Cookie] = authUser

	return authUser.Cookie
}

func createAuthUser(player *manager.ManagedPlayer, ctx context.Context) *AuthUser {
	ctx, cancel := context.WithCancel(ctx)
	return &AuthUser{Cookie: createCookie(), Player: player, cancelContext: cancel, Ctx: ctx}
}

func createCookie() string {
	return uuid.New().String()
}
