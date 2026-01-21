package service

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type MockServices struct {
	GetUser     func(context.Context) (string, error)
	CheckAccess func(context.Context) (bool, error)
	GetContext  func(context.Context) (string, error)
}

type UserService struct{}
type AccessService struct{}
type MemoryService struct{}

func NewMockServices() *MockServices {
	user := &UserService{}
	access := &AccessService{}
	memory := &MemoryService{}

	return &MockServices{
		GetUser:     user.GetUser,
		CheckAccess: access.CheckAccess,
		GetContext:  memory.GetContext,
	}
}

func (u *UserService) GetUser(ctx context.Context) (string, error) {
	//моделирование поведения ошибки от сервиса
	if rand.Intn(10) == 0 {
		return "", errors.New("user service internal error")
	}

	select {
	case <-time.After(10 * time.Millisecond):
		return "user", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (a *AccessService) CheckAccess(ctx context.Context) (bool, error) {
	//моделирование поведения ошибки от сервиса
	if rand.Intn(10) == 0 {
		return false, errors.New("permission service internal error")
	}

	select {
	case <-time.After(50 * time.Millisecond):
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (m *MemoryService) GetContext(ctx context.Context) (string, error) {
	delay := time.Duration(rand.Intn(400)) * time.Millisecond

	select {
	case <-time.After(delay):
		if delay > 500*time.Millisecond {
			return "", errors.New("vector memory request failed")
		}
		return "vector-memory", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
