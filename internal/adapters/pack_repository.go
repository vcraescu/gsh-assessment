package adapters

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/vcraescu/gsh-assessment/internal/domain"
)

//go:embed packs.json
var packs json.RawMessage

type PackRepository struct {
	data []domain.Pack
}

func NewPackRepository() (*PackRepository, error) {
	r := &PackRepository{}

	if err := json.Unmarshal(packs, &r.data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return r, nil
}

func (r *PackRepository) FindAll(_ context.Context) ([]domain.Pack, error) {
	out := make([]domain.Pack, 0, len(r.data))

	for _, pack := range r.data {
		out = append(out, pack)
	}

	return out, nil
}
