package controller

import (
	"context"
	"github.com/khaldeezal/subscriptions-service/internal/domain"
)

type serviceSubscription interface {
	Create(ctx context.Context, in *domain.CreateInput) (string, error)
	Get(ctx context.Context, id string) (*domain.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, f *domain.ListFilter) ([]domain.Subscription, error)
	Update(ctx context.Context, id string, in *domain.CreateInput) error
}
