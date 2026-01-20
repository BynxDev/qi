# QiCard Payment Gateway Go Client

A Go client library for the QiCard Payment Gateway API.

## Installation

```bash
go get github.com/bynx.dev/qi
```

## Usage

### Creating a Client

```go
import "github.com/bynx.dev/qi"

// Create a client with terminal ID and basic auth
client := qi.NewClient("your-terminal-id",
    qi.WithBasicAuth("username", "password"),
)

// Or with signature-based authentication
client := qi.NewClient("your-terminal-id",
    qi.WithSignature("your-signature"),
)

// With custom base URL (for sandbox/testing)
client := qi.NewClient("your-terminal-id",
    qi.WithBaseURL("https://sandbox.qi.iq/api/v1"),
    qi.WithBasicAuth("username", "password"),
)
```

### Creating a Payment

```go
payment, err := client.CreatePayment(context.Background(), &qi.CreatePaymentRequest{
    RequestID:       "unique-request-id",
    Amount:          100.50,
    Currency:        "IQD",
    FinishPaymentURL: "https://yoursite.com/payment/complete",
    NotificationURL:  "https://yoursite.com/webhooks/payment",
    CustomerInfo: &qi.CustomerInfo{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
        Phone:     "009647xxxxxxxxx",
    },
})
if err != nil {
    log.Fatal(err)
}

// Redirect user to payment.FormURL
fmt.Println("Payment form URL:", payment.FormURL)
```

### Getting Payment Status

```go
// By payment ID
status, err := client.GetPaymentStatus(context.Background(), "payment-id")

// By request ID
status, err := client.GetPaymentStatusByRequest(context.Background(), "request-id")

if status.Status == qi.PaymentStatusSuccess {
    fmt.Println("Payment successful!")
}
```

### Canceling a Payment

```go
// By payment ID
resp, err := client.CancelPayment(context.Background(), "payment-id", &qi.CancelPaymentRequest{
    RequestID: "cancel-request-id",
    Amount:    50.00, // Optional: partial cancel
})

// By request ID
resp, err := client.CancelPaymentByRequest(context.Background(), "request-id", &qi.CancelPaymentRequest{
    RequestID: "cancel-request-id",
})
```

### Refunding a Payment

```go
// By payment ID
refund, err := client.RefundPayment(context.Background(), "payment-id", &qi.CreateRefundRequest{
    RequestID: "refund-request-id",
    Amount:    25.00, // Optional: partial refund
    Message:   "Customer requested refund",
})

// By request ID
refund, err := client.RefundPaymentByRequest(context.Background(), "request-id", &qi.CreateRefundRequest{
    RequestID: "refund-request-id",
})
```

### Error Handling

```go
payment, err := client.CreatePayment(ctx, req)
if err != nil {
    if apiErr, ok := err.(*qi.APIError); ok {
        if apiErr.IsValidationError() {
            fmt.Println("Validation error:", apiErr.Err.Error.Message)
        } else if apiErr.IsNotFound() {
            fmt.Println("Not found")
        } else if apiErr.IsAuthenticationError() {
            fmt.Println("Authentication failed")
        }
    }
    log.Fatal(err)
}
```

## Payment Statuses

| Status                    | Description                                 |
| ------------------------- | ------------------------------------------- |
| `CREATED`                 | Payment record created, awaiting processing |
| `FORM_SHOWED`             | Payment form displayed                      |
| `AUTHENTICATION_REQUIRED` | Requires payer authentication (3DS)         |
| `AUTHENTICATION_STARTED`  | Authentication procedure started            |
| `AUTHENTICATION_FAILED`   | Authentication failed (terminal)            |
| `AUTHENTICATED`           | Authentication completed                    |
| `INITIALIZED`             | Payment initialized                         |
| `STARTED`                 | Financial transaction processing started    |
| `SUCCESS`                 | Payment completed successfully (terminal)   |
| `FAILED`                  | Payment rejected (terminal)                 |
| `ERROR`                   | Payment ended with error (terminal)         |
| `EXPIRED`                 | Payment expired (terminal)                  |

## License

MIT
