package signal

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HandleSignals(logger *slog.Logger, funcs ...func()) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs

	logger.Info("received signal", "signal", sig)

	for _, fn := range funcs {
		fn()
	}

	logger.Info("signals handled")

	time.Sleep(1 * time.Second)
	os.Exit(0)
}
