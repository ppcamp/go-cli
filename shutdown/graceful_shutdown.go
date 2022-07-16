package shutdown

import (
	"context"
	"os"
	"os/signal"
)

type function func(context.Context) error

// Graceful is a blocking function that will spawn a new goroutine to process your function.
// When the user cancel the current context or when the user type Ctrl+C in the keyboard, it'll
// notify the function context to cancel anD return.
func Graceful(c context.Context, fn function) error {
	ctx, stop := signal.NotifyContext(c, os.Interrupt)
	defer stop()

	errorCh := make(chan error, 1)

	go func(chan error) {
		errorCh <- fn(ctx)
	}(errorCh)

	select {
	case <-ctx.Done():
		return nil
	case err := <-errorCh:
		return err
	}
}
