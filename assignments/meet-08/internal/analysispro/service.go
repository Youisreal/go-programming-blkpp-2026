package analysispro

import (
	"context"
	"encoding/json"
)

type Service struct{ Repository Repository }

func (s Service) Analysispro(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.Analysispro(ctx)
}
