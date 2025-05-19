package storage

import (
	"context"
	"github.com/serjnester/weather-subscription-service/domain/models"
)

type SubscriptionStorage interface {
	Create(ctx context.Context, sub models.Subscription) error
	IsAlreadySubscribed(ctx context.Context, email, city string) (bool, error)
	GetByToken(ctx context.Context, token string) (models.Subscription, error)
	Confirm(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
}

var _ SubscriptionStorage = (*SubImpl)(nil)

type SubImpl struct {
	q *Queries
}

func NewSubscriptionStorage(q *Queries) *SubImpl {
	return &SubImpl{q: q}
}

func (s *SubImpl) Create(ctx context.Context, sub models.Subscription) error {
	_, err := s.q.CreateSubscription(ctx, CreateSubscriptionParams{
		Email:     sub.Email,
		City:      sub.City,
		Token:     sub.Token,
		Frequency: sub.Frequency,
	})
	return err
}

func (s *SubImpl) IsAlreadySubscribed(ctx context.Context, email, city string) (bool, error) {
	count, err := s.q.IsAlreadySubscribed(ctx, IsAlreadySubscribedParams{
		Email: email,
		City:  city,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SubImpl) GetByToken(ctx context.Context, token string) (models.Subscription, error) {
	dbSub, err := s.q.GetSubscriptionByToken(ctx, token)
	if err != nil {
		return models.Subscription{}, err
	}

	return models.Subscription{
		Email:     dbSub.Email,
		City:      dbSub.City,
		Token:     dbSub.Token,
		Confirmed: dbSub.Confirmed,
	}, nil
}

func (s *SubImpl) Confirm(ctx context.Context, token string) error {
	return s.q.ConfirmSubscription(ctx, token)
}

func (s *SubImpl) Unsubscribe(ctx context.Context, token string) error {
	return s.q.Unsubscribe(ctx, token)
}
