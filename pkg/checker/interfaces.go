package checker

import (
	"context"
)

// Checker defines how to check a service's status
type Checker interface {
	Check(ctx context.Context) (bool, string)
}
