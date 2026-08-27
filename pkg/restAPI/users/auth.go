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
	mu     sync.RWMutex
	Cookie UserCookie             `json:"cookie"`
	Player *manager.ManagedPlayer `json:"player"`
	Ctx    context.Context        `json:"-"`

	writerContext context.Context
	writerCancel  context.CancelFunc
	writerID      uint64
	cancelContext context.CancelFunc
	lastSeen      time.Time
}

var (
	userMapMu sync.RWMutex
	UserMap   = make(map[UserCookie]*AuthUser)
)

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

	userMapMu.Lock()
	UserMap[authUser.Cookie] = authUser
	userMapMu.Unlock()

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

	userMapMu.Lock()
	delete(UserMap, user.Cookie)
	userMapMu.Unlock()

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
	user.mu.Lock()
	oldCancel := user.writerCancel
	user.writerID++
	writerID := user.writerID
	user.writerContext, user.writerCancel = context.WithCancel(context.WithValue(requestContext, writerContextKey, writer))
	user.mu.Unlock()
	if oldCancel != nil {
		oldCancel()
	}

	go func() {
		<-requestContext.Done()
		user.mu.Lock()
		if user.writerID == writerID {
			user.writerCancel = nil
		}
		user.mu.Unlock()
	}()
}

func (user *AuthUser) WriteContext() context.Context {
	user.mu.RLock()
	defer user.mu.RUnlock()
	return user.writerContext
}

func (user *AuthUser) GetUserWriter() http.ResponseWriter {
	user.mu.RLock()
	defer user.mu.RUnlock()
	writer := user.writerContext.Value(writerContextKey).(http.ResponseWriter)
	return writer
}
