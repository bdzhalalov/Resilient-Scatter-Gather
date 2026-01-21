package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bdzhalalov/resilient-scatter-gather/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func userServiceMock(_ context.Context) (string, error) { return "user-1", nil }

func accessServiceMock(_ context.Context) (bool, error) { return true, nil }

func memoryServiceMock(delay time.Duration) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		select {
		case <-time.After(delay):
			return "vector-memory", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

var logger = logrus.New()

func TestChatSummaryHandlerAllServicesOK(t *testing.T) {
	svc := &service.MockServices{
		GetUser:     userServiceMock,
		CheckAccess: accessServiceMock,
		GetContext:  memoryServiceMock(50 * time.Millisecond),
	}

	h := New(logger, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/summary", nil)
	rec := httptest.NewRecorder()

	h.GetSummary(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)

	fmt.Println(resp)

	if resp["user"] != "user-1" {
		t.Fatal("user missing")
	}

	if resp["access"] != true {
		t.Fatal("access missing")
	}

	if resp["memory"] != "vector-memory" {
		t.Fatal("response from memory service expected")
	}
}

func TestChatSummaryHandlerMemoryServiceTimeout(t *testing.T) {
	svc := &service.MockServices{
		GetUser:     userServiceMock,
		CheckAccess: accessServiceMock,
		GetContext:  memoryServiceMock(500 * time.Millisecond),
	}

	h := New(logger, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/summary", nil)
	rec := httptest.NewRecorder()

	h.GetSummary(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)

	if _, ok := resp["context"]; ok {
		t.Fatal("context should be absent")
	}
}

func TestChatSummaryHandlerUserServiceError(t *testing.T) {
	svc := &service.MockServices{
		GetUser:     func(ctx context.Context) (string, error) { return "", errors.New("internal error") },
		CheckAccess: accessServiceMock,
		GetContext:  memoryServiceMock(10 * time.Millisecond),
	}

	h := New(logger, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/summary", nil)
	rec := httptest.NewRecorder()

	h.GetSummary(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestChatSummaryHandlerAccessServiceError(t *testing.T) {
	svc := &service.MockServices{
		GetUser:     userServiceMock,
		CheckAccess: func(ctx context.Context) (bool, error) { return false, errors.New("internal error") },
		GetContext:  memoryServiceMock(10 * time.Millisecond),
	}

	h := New(logger, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/summary", nil)
	rec := httptest.NewRecorder()

	h.GetSummary(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestChatSummary_PermissionTimeout(t *testing.T) {
	svc := &service.MockServices{
		GetUser: userServiceMock,
		CheckAccess: func(ctx context.Context) (bool, error) {
			select {
			case <-time.After(300 * time.Millisecond):
				return true, nil
			case <-ctx.Done():
				return false, ctx.Err()
			}
		},
		GetContext: memoryServiceMock(10 * time.Millisecond),
	}

	h := New(logger, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/chat/summary", nil)
	rec := httptest.NewRecorder()

	h.GetSummary(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
