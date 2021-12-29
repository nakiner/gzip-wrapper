// Code generated by gitlab.eldorado.ru/golang/go-kit-service-generator  DO NOT EDIT.
package logging

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx)
	assert.Equal(t, fallbackLogger, logger)
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	expected := log.NewNopLogger()
	ctx = WithContext(ctx, expected)
	actual := FromContext(ctx)
	assert.Equal(t, expected, actual)
}
