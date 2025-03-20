package tests

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// ContextMatcher is a custom matcher for context in mock expectations
type ContextMatcher struct{}

// Matches checks if the given object is a context.Context
func (m ContextMatcher) Matches(x interface{}) bool {
	_, ok := x.(context.Context)
	return ok
}

// String returns a description of the matcher
func (m ContextMatcher) String() string {
	return "is a context.Context"
}

// AnyContext returns a matcher that matches any context
func AnyContext() interface{} {
	return mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})
}
