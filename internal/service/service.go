package service

import (
	"context"

	"github.com/khaldeezal/subscriptions-service/internal/domain"
)

type subscriptionService struct {
	subscriptionRepository repositorySubscription
}

func NewSubscriptionService(subscriptionRepository repositorySubscription) *subscriptionService {
	return &subscriptionService{
		subscriptionRepository: subscriptionRepository,
	}
}

func (s subscriptionService) Create(ctx context.Context, in *domain.CreateInput) (string, error) {
	return s.subscriptionRepository.Create(ctx, in)
}

func (s subscriptionService) Get(ctx context.Context, id string) (*domain.Subscription, error) {
	return s.subscriptionRepository.Get(ctx, id)
}

func (s subscriptionService) Delete(ctx context.Context, id string) error {
	return s.subscriptionRepository.Delete(ctx, id)
}

func (s subscriptionService) List(ctx context.Context, f *domain.ListFilter) ([]domain.Subscription, error) {
	return s.subscriptionRepository.List(ctx, f)
}

func (s subscriptionService) Update(ctx context.Context, id string, in *domain.CreateInput) error {
	return s.subscriptionRepository.Update(ctx, id, in)
}
