package users

import (
	"time"
)

func ensureUserActive(user *AuthUser) {
	inactiveTimer := time.NewTimer(30 * time.Second) // Timer for inactivity
	defer inactiveTimer.Stop()                       // Stop the timer when done

	for {
		select {
		case <-inactiveTimer.C:
			// User is inactive for 30 seconds
			removeAuthUser(user)
			return
		case <-user.Ctx.Done():
			// User context is canceled
			removeAuthUser(user)
			return
		case <-time.After(time.Second):
			// Reset the timer if user is active
			if time.Since(user.lastSeen) < 30*time.Second || user.writerCancel != nil {
				if !inactiveTimer.Stop() {
					<-inactiveTimer.C // Drain the timer channel
				}
				inactiveTimer.Reset(10 * time.Second)
			} else {
				// User is now gone for longer than 30 seconds
				// TODO: might want to send user a message if they are gone for this long -_-
			}
		}
	}
}
