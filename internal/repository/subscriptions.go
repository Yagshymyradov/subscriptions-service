package repository

import (
	"context"

	"github.com/Yagshymyradov/subscriptions-service/internal/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	List(ctx context.Context, userID string) ([]models.Subscription, error)
	Update(ctx context.Context, s *models.Subscription) error
	Delete(ctx context.Context, id int) error
	TotalCost(ctx context.Context, userID string, month, year int, serviceFilter string) (int, error)
}
