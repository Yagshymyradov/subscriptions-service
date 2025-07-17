package service

import (
	"context"

	"github.com/Yagshymyradov/subscriptions-service/internal/models"
	"github.com/Yagshymyradov/subscriptions-service/internal/repository"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(r repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: r}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) Get(ctx context.Context, id int) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) List(ctx context.Context, userID string) ([]models.Subscription, error) {
	return s.repo.List(ctx, userID)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) TotalCost(ctx context.Context, userID string, month, year int, serviceFilter string) (int, error) {
	return s.repo.TotalCost(ctx, userID, month, year, serviceFilter)
}
