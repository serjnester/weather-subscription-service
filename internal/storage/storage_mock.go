package storage

import (
	"context"
	"github.com/serjnester/weather-subscription-service/domain/models"
)

type MockStorage struct {
	CreateFn              func(ctx context.Context, sub models.Subscription) error
	IsAlreadySubscribedFn func(ctx context.Context, email, city string) (bool, error)
	GetByTokenFn          func(ctx context.Context, token string) (models.Subscription, error)
	ConfirmFn             func(ctx context.Context, token string) error
	UnsubscribeFn         func(ctx context.Context, token string) error
}

func (m *MockStorage) Create(ctx context.Context, sub models.Subscription) error {
	return m.CreateFn(ctx, sub)
}

func (m *MockStorage) IsAlreadySubscribed(ctx context.Context, email, city string) (bool, error) {
	return m.IsAlreadySubscribedFn(ctx, email, city)
}

func (m *MockStorage) GetByToken(ctx context.Context, token string) (models.Subscription, error) {
	return m.GetByTokenFn(ctx, token)
}

func (m *MockStorage) Confirm(ctx context.Context, token string) error {
	return m.ConfirmFn(ctx, token)
}

func (m *MockStorage) Unsubscribe(ctx context.Context, token string) error {
	return m.UnsubscribeFn(ctx, token)
}
