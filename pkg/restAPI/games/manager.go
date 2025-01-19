package games

import (
	"github.com/kraanter/blackjack/pkg/manager"
)

var managerSettings = createManagerSettings()
var GameManager = manager.CreateManager(managerSettings)

func createManagerSettings() *manager.Settings {
	settings := manager.CreateSettings()

	return settings
}
