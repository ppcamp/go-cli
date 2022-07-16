package shutdown_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ppcamp/go-lib/shutdown"
	"github.com/stretchr/testify/assert"
)

func TestGraceful(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("should fail")

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	onFailure := func(_ context.Context) error { return err }
	onSuccess := func(_ context.Context) error { return nil }

	assert.ErrorIs(shutdown.Graceful(ctx, onFailure), err)

	// cancel by some user action Ctrl+C or due to the upper ctx
	cancel()
	assert.Nil(shutdown.Graceful(ctx, onSuccess))
}
