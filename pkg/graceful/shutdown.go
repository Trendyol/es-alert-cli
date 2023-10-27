package graceful

import (
	"context"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Shutdown shuts down the app gracefully when receiving an os.Interrupt or syscall.SIGTERM signal.
func Shutdown(timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Infof("shutting down cli app with %s timeout", timeout)

	//close something here
}
