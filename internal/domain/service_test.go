package domain_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/gsh-assessment/internal/domain"
	"testing"
)

var _ domain.PackRepository = (*PackRepository)(nil)

type PackRepository struct {
	mock.Mock
}

func (r *PackRepository) FindAll(ctx context.Context) ([]domain.Pack, error) {
	args := r.Called(ctx)
	out, _ := args.Get(0).([]domain.Pack)

	return out, args.Error(1)
}

func TestOrderService_Create(t *testing.T) {
	t.Parallel()

	type fields struct {
		repository *PackRepository
	}

	type args struct {
		quantity int
	}

	packs := []domain.Pack{
		{
			Size: 250,
		},
		{
			Size: 500,
		},
		{
			Size: 1000,
		},
		{
			Size: 2000,
		},
		{
			Size: 5000,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		on      func(t *testing.T, f fields)
		want    domain.Order
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "quantity is zero",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.Error(t, err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "quantity is less than zero",
			args: args{
				quantity: -2,
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.Error(t, err, domain.ErrInvalidArgument)
			},
		},
		{
			name: "one item",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 1,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     250,
					},
				},
			},
		},
		{
			name: "250 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 250,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     250,
					},
				},
			},
		},
		{
			name: "251 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 251,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     500,
					},
				},
			},
		},
		{
			name: "500 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 500,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     500,
					},
				},
			},
		},
		{
			name: "501 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 501,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     500,
					},
					{
						Quantity: 1,
						Pack:     250,
					},
				},
			},
		},
		{
			name: "750 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 750,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     500,
					},
					{
						Quantity: 1,
						Pack:     250,
					},
				},
			},
		},
		{
			name: "751 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 751,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     1000,
					},
				},
			},
		},
		{
			name: "1001 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 1001,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     1000,
					},
					{
						Quantity: 1,
						Pack:     250,
					},
				},
			},
		},
		{
			name: "2455 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 2455,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 1,
						Pack:     2000,
					},
					{
						Quantity: 1,
						Pack:     500,
					},
				},
			},
		},
		{
			name: "12000 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 12000,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 2,
						Pack:     5000,
					},
					{
						Quantity: 1,
						Pack:     2000,
					},
				},
			},
		},
		{
			name: "12001 items",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 12001,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 2,
						Pack:     5000,
					},
					{
						Quantity: 1,
						Pack:     2000,
					},
					{
						Quantity: 1,
						Pack:     250,
					},
				},
			},
		},
		{
			name: "very large number",
			fields: fields{
				repository: &PackRepository{},
			},
			args: args{
				quantity: 1000000,
			},
			on: func(t *testing.T, f fields) {
				f.repository.
					On("FindAll", mock.Anything).
					Return(packs, nil)
			},
			want: domain.Order{
				Rows: []domain.OrderRow{
					{
						Quantity: 200,
						Pack:     5000,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.on != nil {
				tt.on(t, tt.fields)
			}

			svc := domain.NewOrderService(tt.fields.repository)
			got, err := svc.Create(context.Background(), tt.args.quantity)

			if tt.wantErr != nil {
				tt.wantErr(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
