package analysiscust

import (
	"context"
	"encoding/json"
)

type Service struct{ Repository Repository }

func (s Service) Analysiscust(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.Analysiscust(ctx)
}
