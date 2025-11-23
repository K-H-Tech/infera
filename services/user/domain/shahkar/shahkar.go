package shahkar

import (
	"context"

	"zarinpal-platform/core/trace"
)

type Shahkar interface {
	IsShahkarValid(ctx context.Context, mobile string, nationalCode string) (bool, error)
}

type shahkar struct {
}

func NewShahkar() Shahkar {
	return &shahkar{}
}
func (s *shahkar) IsShahkarValid(ctx context.Context, mobile string, nationalCode string) (bool, error) {
	_, span := trace.GetTracer().Start(ctx, "Shahkar.IsShahkarValid")
	defer span.End()

	// Implement the logic to check Shahkar validity

	return true, nil
}
