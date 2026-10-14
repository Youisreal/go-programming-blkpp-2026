package product

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"unicode/utf8"
)

type Service struct{ Repository Repository }

func normalize(p *Input) error {
	p.ProductName = strings.TrimSpace(p.ProductName)
	if p.ProductName == "" || utf8.RuneCountInString(p.ProductName) > 40 || strings.ContainsRune(p.ProductName, '\x00') ||
		utf8.RuneCountInString(p.QuantityPerUnit) > 20 || strings.ContainsRune(p.QuantityPerUnit, '\x00') ||
		p.SupplierID <= 0 || p.CategoryID <= 0 || p.UnitPrice < 0 || p.UnitPrice > 9999999999.99 ||
		math.IsNaN(p.UnitPrice) || math.IsInf(p.UnitPrice, 0) ||
		p.UnitsInStock < 0 || p.UnitsInStock > 32767 || p.UnitsOnOrder < 0 || p.UnitsOnOrder > 32767 ||
		p.ReorderLevel < 0 || p.ReorderLevel > 32767 || (p.Discontinued != 0 && p.Discontinued != 1) {
		return ErrInvalidInput
	}
	return nil
}

func one(rows json.RawMessage, err error) (json.RawMessage, error) {
	if err != nil {
		return nil, err
	}
	var result []json.RawMessage
	if err := json.Unmarshal(rows, &result); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrNotFound
	}
	return result[0], nil
}

func (s Service) List(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.List(ctx)
}
func (s Service) Get(ctx context.Context, id int) (json.RawMessage, error) {
	return one(s.Repository.Get(ctx, id))
}
func (s Service) Create(ctx context.Context, p Input) (json.RawMessage, error) {
	if err := normalize(&p); err != nil {
		return nil, err
	}
	return one(s.Repository.Create(ctx, p))
}
func (s Service) Update(ctx context.Context, id int, p Input) (json.RawMessage, error) {
	if err := normalize(&p); err != nil {
		return nil, err
	}
	return one(s.Repository.Update(ctx, id, p))
}
func (s Service) Delete(ctx context.Context, id int) (json.RawMessage, error) {
	return one(s.Repository.Delete(ctx, id))
}
func (s Service) Categories(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.Categories(ctx)
}
func (s Service) Suppliers(ctx context.Context) (json.RawMessage, error) {
	return s.Repository.Suppliers(ctx)
}
