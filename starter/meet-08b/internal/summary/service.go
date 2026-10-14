package summary

import (
	"context"
	"encoding/json"
)

type Service struct{ Repository Repository }

func (s Service) Summary(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.Summary(ctx)
}
