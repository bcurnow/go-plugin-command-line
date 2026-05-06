package plugin

import (
	"fmt"

	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	// Removed import to avoid import cycle
)

var logger = logging.Logger().Named("plugin")

func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
