package qi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BynxDev/qi"
)

func TestCreatePayment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/payment" {
			t.Errorf("expected /payment, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Terminal-Id") != "test-terminal" {
			t.Errorf("expected X-Terminal-Id header")
		}

		response := qi.Payment{
			RequestID:    "test-request-id",
			PaymentID:    "test-payment-id",
			Status:       qi.PaymentStatusCreated,
			Amount:       100.50,
			Currency:     "IQD",
			CreationDate: time.Now(),
			FormURL:      "https://example.com/pay",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := qi.NewClient("test-terminal", qi.WithBaseURL(server.URL))

	payment, err := client.CreatePayment(context.Background(), &qi.CreatePaymentRequest{
		RequestID: "test-request-id",
		Amount:    100.50,
		Currency:  "IQD",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payment.PaymentID != "test-payment-id" {
		t.Errorf("expected payment ID test-payment-id, got %s", payment.PaymentID)
	}

	if payment.Status != qi.PaymentStatusCreated {
		t.Errorf("expected status CREATED, got %s", payment.Status)
	}
}

func TestGetPaymentStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/payment/test-payment-id/status" {
			t.Errorf("expected /payment/test-payment-id/status, got %s", r.URL.Path)
		}

		response := qi.PaymentStatusResponse{
			RequestID:    "test-request-id",
			PaymentID:    "test-payment-id",
			Status:       qi.PaymentStatusSuccess,
			Amount:       100.50,
			Currency:     "IQD",
			CreationDate: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := qi.NewClient("test-terminal", qi.WithBaseURL(server.URL))

	status, err := client.GetPaymentStatus(context.Background(), "test-payment-id")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Status != qi.PaymentStatusSuccess {
		t.Errorf("expected status SUCCESS, got %s", status.Status)
	}
}

func TestCancelPayment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/payment/test-payment-id/cancel" {
			t.Errorf("expected /payment/test-payment-id/cancel, got %s", r.URL.Path)
		}

		response := qi.PaymentCancelResponse{
			RequestID:    "cancel-request-id",
			PaymentID:    "test-payment-id",
			Status:       qi.PaymentStatusCreated,
			Canceled:     true,
			Amount:       100.50,
			Currency:     "IQD",
			CreationDate: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := qi.NewClient("test-terminal", qi.WithBaseURL(server.URL))

	resp, err := client.CancelPayment(context.Background(), "test-payment-id", &qi.CancelPaymentRequest{
		RequestID: "cancel-request-id",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Canceled {
		t.Error("expected payment to be canceled")
	}
}

func TestRefundPayment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/payment/test-payment-id/refund" {
			t.Errorf("expected /payment/test-payment-id/refund, got %s", r.URL.Path)
		}

		response := qi.Refund{
			RefundID:     "test-refund-id",
			PaymentID:    "test-payment-id",
			Amount:       50.25,
			Currency:     "IQD",
			CreationDate: time.Now(),
			Status:       qi.RefundStatusSuccess,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := qi.NewClient("test-terminal", qi.WithBaseURL(server.URL))

	refund, err := client.RefundPayment(context.Background(), "test-payment-id", &qi.CreateRefundRequest{
		RequestID: "refund-request-id",
		Amount:    50.25,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if refund.RefundID != "test-refund-id" {
		t.Errorf("expected refund ID test-refund-id, got %s", refund.RefundID)
	}

	if refund.Status != qi.RefundStatusSuccess {
		t.Errorf("expected status SUCCESS, got %s", refund.Status)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := qi.Error{
			Error: qi.ErrorDetails{
				Code:    qi.ErrorCodeValidationError,
				Message: qi.ErrorMessageValidationError,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := qi.NewClient("test-terminal", qi.WithBaseURL(server.URL))

	_, err := client.CreatePayment(context.Background(), &qi.CreatePaymentRequest{
		RequestID: "test-request-id",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*qi.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if !apiErr.IsValidationError() {
		t.Error("expected validation error")
	}
}
