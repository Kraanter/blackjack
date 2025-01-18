package users

import (
	"context"
	"fmt"

	"github.com/kraanter/blackjack/pkg/manager"
)

type ContextPlayer string

const ContextUserKey ContextPlayer = "user"

type UserCookie = string

type AuthUser struct {
	Cookie        UserCookie             `json:"cookie"`
	Player        *manager.ManagedPlayer `json:"player"`
	Ctx           context.Context        `json:"-"`
	cancelContext context.CancelFunc     `json:"-"`
}

var UserMap = make(map[UserCookie]*AuthUser)

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
