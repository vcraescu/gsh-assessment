package domain

import (
	"context"
	"fmt"
	"math"
	"sort"
)

type PackRepository interface {
	FindAll(ctx context.Context) ([]Pack, error)
}

type OrderService struct {
	repository PackRepository
}

func NewOrderService(repository PackRepository) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Create(ctx context.Context, quantity int) (Order, error) {
	out := Order{}

	if quantity <= 0 {
		return out, fmt.Errorf("quantity must be greater than zero; got %v: %w", quantity, ErrInvalidArgument)
	}

	packs, err := s.repository.FindAll(ctx)
	if err != nil {
		return out, fmt.Errorf("findAll: %w", err)
	}

	if len(packs) == 0 {
		return out, nil
	}

	// sort packs descending
	sort.Slice(packs, func(i, j int) bool {
		return packs[i].Size > packs[j].Size
	})

	minPack := packs[len(packs)-1]
	quantity = int(math.Ceil(float64(quantity)/float64(minPack.Size))) * minPack.Size

	for _, pack := range packs {
		if packQuantity := quantity / pack.Size; packQuantity > 0 {
			quantity %= pack.Size

			out.Rows = append(out.Rows, OrderRow{
				Quantity: packQuantity,
				Pack:     pack.Size,
			})
		}
	}

	// sort rows descending to ensure predictable output
	sort.Slice(out.Rows, func(i, j int) bool {
		return out.Rows[i].Pack > out.Rows[j].Pack
	})

	return out, nil
}
