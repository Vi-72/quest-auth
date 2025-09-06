package assertions

import (
	"context"
	"time"

	"github.com/stretchr/testify/assert"
)

// For auth service we currently don't persist domain events in DB.
// This file provides placeholders for symmetry and future expansion.

type E2EAssertions struct{ assert *assert.Assertions }

func NewE2EAssertions(a *assert.Assertions) *E2EAssertions { return &E2EAssertions{assert: a} }

// VerifyNoUnexpectedDelay ensures operations finish within expected time window.
func (a *E2EAssertions) VerifyNoUnexpectedDelay(ctx context.Context, started time.Time) {
	a.assert.WithinDuration(time.Now(), started, 2*time.Second)
}
