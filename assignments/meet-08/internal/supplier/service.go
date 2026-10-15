package supplier

import (
	"context"
	"encoding/json"
)

// Service coordinates supplier use cases.
type Service struct{ Repository Repository }

func (s Service) List(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.List(ctx)
}
