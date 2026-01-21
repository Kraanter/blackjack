package users

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kraanter/blackjack/pkg/manager"
	"github.com/kraanter/blackjack/pkg/restAPI/games"
)

type ContextPlayer string

const ContextUserKey ContextPlayer = "user"

type UserCookie = string

type AuthUser struct {
	Cookie UserCookie             `json:"cookie"`
	Player *manager.ManagedPlayer `json:"player"`
	Ctx    context.Context        `json:"-"`

	writerContext context.Context
	writerCancel  context.CancelFunc
	cancelContext context.CancelFunc
	lastSeen      time.Time
}

var UserMap = make(map[UserCookie]*AuthUser)

func getUserFromContext(ctx context.Context) *AuthUser {
	value := ctx.Value(ContextUserKey)

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

func RemoveAuthUser(user *AuthUser) {
	if user != nil {
		user.cancelContext()
	}
}

var userNotFoundError = fmt.Errorf("Could not find user")

func removeAuthUser(user *AuthUser) {
	if user == nil {
		return
	}

	defer user.cancelContext()

	gameId := user.Player.GameId

	// TODO: Store the balance of the user somewhere
	_, err := user.Player.Leave()
	if err != nil {
		println("Error while leaving game", err.Error())
		return
	}

	delete(UserMap, user.Cookie)

	games.GameManager.RemoveGame(gameId)
}

func createAuthUser(player *manager.ManagedPlayer, ctx context.Context) *AuthUser {
	ctx, cancel := context.WithCancel(ctx)
	var once sync.Once
	cancelFunc := func() {
		once.Do(func() {
			cancel()
		})
	}
	user := &AuthUser{Cookie: createCookie(), Player: player, cancelContext: cancelFunc, Ctx: ctx, lastSeen: time.Now()}
	go ensureUserActive(user)
	return user
}

type writerContextType string

var writerContextKey = writerContextType("Writer")

func (user *AuthUser) SetUserWriter(writer http.ResponseWriter, requestContext context.Context) {
	if user.writerCancel != nil {
		user.writerCancel()
	}

	go func() {
		<-requestContext.Done()
		user.writerCancel = nil
	}()

	user.writerContext, user.writerCancel = context.WithCancel(context.WithValue(requestContext, writerContextKey, writer))
}

func (user *AuthUser) WriteContext() context.Context {
	return user.writerContext
}

func (user *AuthUser) GetUserWriter() http.ResponseWriter {
	writer := user.writerContext.Value(writerContextKey).(http.ResponseWriter)
	return writer
}
