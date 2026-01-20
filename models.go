package qi

import (
	"encoding/json"
	"strings"
	"time"
)

// Time is a custom time type that handles the QiCard API's time format.
// The API returns timestamps without timezone suffix (e.g., "2026-01-20T11:57:31").
type Time struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Try RFC3339 first
	parsed, err := time.Parse(time.RFC3339, s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Fall back to the API's format without timezone
	parsed, err = time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Time.Format(time.RFC3339))
}

// NewTime creates a new Time from a time.Time value.
func NewTime(t time.Time) Time {
	return Time{Time: t}
}

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	PaymentStatusCreated                   PaymentStatus = "CREATED"
	PaymentStatusFormShowed                PaymentStatus = "FORM_SHOWED"
	PaymentStatusThreeDSMethodCallRequired PaymentStatus = "THREE_DS_METHOD_CALL_REQUIRED"
	PaymentStatusAuthenticationRequired    PaymentStatus = "AUTHENTICATION_REQUIRED"
	PaymentStatusAuthenticationStarted     PaymentStatus = "AUTHENTICATION_STARTED"
	PaymentStatusAuthenticationFailed      PaymentStatus = "AUTHENTICATION_FAILED"
	PaymentStatusAuthenticated             PaymentStatus = "AUTHENTICATED"
	PaymentStatusInitialized               PaymentStatus = "INITIALIZED"
	PaymentStatusStarted                   PaymentStatus = "STARTED"
	PaymentStatusSuccess                   PaymentStatus = "SUCCESS"
	PaymentStatusFailed                    PaymentStatus = "FAILED"
	PaymentStatusError                     PaymentStatus = "ERROR"
	PaymentStatusExpired                   PaymentStatus = "EXPIRED"
)

// RefundStatus represents the status of a refund.
type RefundStatus string

const (
	RefundStatusSuccess    RefundStatus = "SUCCESS"
	RefundStatusFailed     RefundStatus = "FAILED"
	RefundStatusProcessing RefundStatus = "PROCESSING"
)

// PaymentSystem represents the payment system type.
type PaymentSystem string

const (
	PaymentSystemVisa       PaymentSystem = "VISA"
	PaymentSystemMasterCard PaymentSystem = "MASTER_CARD"
)

// PaymentType represents the type of payment.
type PaymentType string

const (
	PaymentTypePaymentToken PaymentType = "PAYMENT_TOKEN"
)

// PaymentTokenType represents the type of payment token.
type PaymentTokenType string

const (
	PaymentTokenTypeAuth     PaymentTokenType = "AUTH"
	PaymentTokenTypeNonRecur PaymentTokenType = "NON_RECUR"
	PaymentTokenTypeUnauth   PaymentTokenType = "UNAUTH"
)

// ItemPaymentMethod represents the payment method for an item.
type ItemPaymentMethod string

const (
	ItemPaymentMethodFullPayment    ItemPaymentMethod = "FULL_PAYMENT"
	ItemPaymentMethodFullPrepayment ItemPaymentMethod = "FULL_PREPAYMENT"
	ItemPaymentMethodPrepayment     ItemPaymentMethod = "PREPAYMENT"
	ItemPaymentMethodAdvance        ItemPaymentMethod = "ADVANCE"
	ItemPaymentMethodPartialPayment ItemPaymentMethod = "PARTIAL_PAYMENT"
	ItemPaymentMethodCredit         ItemPaymentMethod = "CREDIT"
	ItemPaymentMethodCreditPayment  ItemPaymentMethod = "CREDIT_PAYMENT"
)

// ItemPaymentObject represents the payment object type for an item.
type ItemPaymentObject string

const (
	ItemPaymentObjectCommodity            ItemPaymentObject = "COMMODITY"
	ItemPaymentObjectExcise               ItemPaymentObject = "EXCISE"
	ItemPaymentObjectJob                  ItemPaymentObject = "JOB"
	ItemPaymentObjectService              ItemPaymentObject = "SERVICE"
	ItemPaymentObjectGamblingBet          ItemPaymentObject = "GAMBLING_BET"
	ItemPaymentObjectGamblingPrize        ItemPaymentObject = "GAMBLING_PRIZE"
	ItemPaymentObjectLottery              ItemPaymentObject = "LOTTERY"
	ItemPaymentObjectLotteryPrize         ItemPaymentObject = "LOTTERY_PRIZE"
	ItemPaymentObjectIntellectualActivity ItemPaymentObject = "INTELLECTUAL_ACTIVITY"
	ItemPaymentObjectPayment              ItemPaymentObject = "PAYMENT"
	ItemPaymentObjectAgentCommission      ItemPaymentObject = "AGENT_COMMISSION"
	ItemPaymentObjectComposite            ItemPaymentObject = "COMPOSITE"
	ItemPaymentObjectAnother              ItemPaymentObject = "ANOTHER"
)

// ItemTax represents the VAT rate for an item.
type ItemTax string

const (
	ItemTaxNone   ItemTax = "NONE"
	ItemTaxVAT0   ItemTax = "VAT0"
	ItemTaxVAT10  ItemTax = "VAT10"
	ItemTaxVAT20  ItemTax = "VAT20"
	ItemTaxVAT110 ItemTax = "VAT110"
	ItemTaxVAT120 ItemTax = "VAT120"
)

// CreatePaymentRequest represents a request to create a payment.
type CreatePaymentRequest struct {
	RequestID        string            `json:"requestId"`
	Amount           float64           `json:"amount,omitempty"`
	Currency         string            `json:"currency,omitempty"`
	Locale           string            `json:"locale,omitempty"`
	FinishPaymentURL string            `json:"finishPaymentUrl,omitempty"`
	NotificationURL  string            `json:"notificationUrl,omitempty"`
	CustomerInfo     *CustomerInfo     `json:"customerInfo,omitempty"`
	BrowserInfo      *BrowserInfo      `json:"browserInfo,omitempty"`
	AdditionalInfo   map[string]string `json:"additionalInfo,omitempty"`
}

// Payment represents payment details returned from the API.
type Payment struct {
	RequestID      string            `json:"requestId"`
	PaymentID      string            `json:"paymentId"`
	Status         PaymentStatus     `json:"status"`
	Canceled       bool              `json:"canceled,omitempty"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	CreationDate   Time              `json:"creationDate"`
	FormURL        string            `json:"formUrl,omitempty"`
	AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
}

// PaymentStatusResponse represents the response when getting payment status.
type PaymentStatusResponse struct {
	RequestID       string            `json:"requestId"`
	PaymentID       string            `json:"paymentId"`
	Status          PaymentStatus     `json:"status"`
	Canceled        bool              `json:"canceled,omitempty"`
	Amount          float64           `json:"amount"`
	ConfirmedAmount float64           `json:"confirmedAmount,omitempty"`
	Currency        string            `json:"currency"`
	PaymentType     string            `json:"paymentType,omitempty"`
	CreationDate    Time              `json:"creationDate"`
	Details         *PaymentDetails   `json:"details,omitempty"`
	AdditionalInfo  map[string]string `json:"additionalInfo,omitempty"`
}

// PaymentDetails contains detailed information about a payment.
type PaymentDetails struct {
	ResultCode    string                 `json:"resultCode,omitempty"`
	RRN           string                 `json:"rrn,omitempty"`
	AuthID        string                 `json:"authId,omitempty"`
	AuthDate      *Time                  `json:"authDate,omitempty"`
	MaskedPan     string                 `json:"maskedPan,omitempty"`
	PaymentSystem PaymentSystem          `json:"paymentSystem,omitempty"`
	CustomDetails map[string]interface{} `json:"customDetails,omitempty"`
}

// CancelPaymentRequest represents a request to cancel a payment.
type CancelPaymentRequest struct {
	RequestID string  `json:"requestId,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
}

// PaymentCancelResponse represents the response when canceling a payment.
type PaymentCancelResponse struct {
	RequestID      string            `json:"requestId"`
	PaymentID      string            `json:"paymentId"`
	Status         PaymentStatus     `json:"status"`
	Canceled       bool              `json:"canceled"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	CreationDate   Time              `json:"creationDate"`
	Cancels        []Cancel          `json:"cancels,omitempty"`
	AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
}

// Cancel represents cancellation details.
type Cancel struct {
	RequestID    string  `json:"requestId,omitempty"`
	Created      Time    `json:"created"`
	Successfully bool    `json:"successfully"`
	Amount       float64 `json:"amount"`
}

// CreateRefundRequest represents a request to create a refund.
type CreateRefundRequest struct {
	RequestID string           `json:"requestId,omitempty"`
	Amount    float64          `json:"amount,omitempty"`
	Message   string           `json:"message,omitempty"`
	ExtParams *RefundExtParams `json:"extParams,omitempty"`
}

// RefundExtParams contains additional refund parameters.
type RefundExtParams struct {
	Phone              string `json:"phone,omitempty"`
	RecipientBankID    string `json:"recipientBankId,omitempty"`
	ProcessRefundAsOCT bool   `json:"processRefundAsOct,omitempty"`
}

// Refund represents refund details returned from the API.
type Refund struct {
	RefundID     string          `json:"refundId"`
	RequestID    string          `json:"requestId,omitempty"`
	PaymentID    string          `json:"paymentId"`
	Amount       float64         `json:"amount"`
	Currency     string          `json:"currency"`
	CreationDate Time            `json:"creationDate"`
	Message      string          `json:"message,omitempty"`
	Details      *PaymentDetails `json:"details,omitempty"`
	Status       RefundStatus    `json:"status"`
	Canceled     bool            `json:"canceled,omitempty"`
	Cancels      []Cancel        `json:"cancels,omitempty"`
}

// CustomerInfo contains customer details.
type CustomerInfo struct {
	FirstName                    string `json:"firstName,omitempty"`
	MiddleName                   string `json:"middleName,omitempty"`
	LastName                     string `json:"lastName,omitempty"`
	Phone                        string `json:"phone,omitempty"`
	Email                        string `json:"email,omitempty"`
	AccountID                    string `json:"accountId,omitempty"`
	AccountNumber                string `json:"accountNumber,omitempty"`
	Address                      string `json:"address,omitempty"`
	City                         string `json:"city,omitempty"`
	ProvinceCode                 string `json:"provinceCode,omitempty"`
	CountryCode                  string `json:"countryCode,omitempty"`
	PostalCode                   string `json:"postalCode,omitempty"`
	BirthDate                    string `json:"birthDate,omitempty"`
	IdentificationType           string `json:"identificationType,omitempty"`
	IdentificationNumber         string `json:"identificationNumber,omitempty"`
	IdentificationCountryCode    string `json:"identificationCountryCode,omitempty"`
	IdentificationExpirationDate string `json:"identificationExpirationDate,omitempty"`
	Nationality                  string `json:"nationality,omitempty"`
	CountryOfBirth               string `json:"countryOfBirth,omitempty"`
	FundSource                   string `json:"fundSource,omitempty"`
	ParticipantID                string `json:"participantId,omitempty"`
	AdditionalMessage            string `json:"additionalMessage,omitempty"`
	TransactionReason            string `json:"transactionReason,omitempty"`
	ClaimCode                    string `json:"claimCode,omitempty"`
}

// BrowserInfo contains browser details for 3DS authentication.
type BrowserInfo struct {
	BrowserAcceptHeader string `json:"browserAcceptHeader"`
	BrowserIP           string `json:"browserIp"`
	BrowserJavaEnabled  bool   `json:"browserJavaEnabled"`
	BrowserLanguage     string `json:"browserLanguage"`
	BrowserColorDepth   string `json:"browserColorDepth"`
	BrowserScreenWidth  string `json:"browserScreenWidth"`
	BrowserScreenHeight string `json:"browserScreenHeight"`
	BrowserTZ           string `json:"browserTZ"`
	BrowserUserAgent    string `json:"browserUserAgent"`
}

// AuthenticateInfo contains data for 3DS authentication.
type AuthenticateInfo struct {
	URL            string              `json:"url"`
	CardholderInfo string              `json:"cardholderInfo,omitempty"`
	Params         *AuthenticateParams `json:"params,omitempty"`
}

// AuthenticateParams contains parameters for 3DS authentication.
type AuthenticateParams struct {
	PaReq   string `json:"paReq,omitempty"`
	MD      string `json:"md,omitempty"`
	TermURL string `json:"termUrl,omitempty"`
	CReq    string `json:"creq,omitempty"`
}

// PaymentData contains payment details for token-based payments.
type PaymentData struct {
	PaymentType  PaymentType `json:"paymentType"`
	PaymentToken string      `json:"paymentToken,omitempty"`
}

// ItemsInfo contains purchase information.
type ItemsInfo struct {
	Description string        `json:"description,omitempty"`
	Items       []PaymentItem `json:"items,omitempty"`
}

// PaymentItem represents an item in a purchase.
type PaymentItem struct {
	Name          string            `json:"name,omitempty"`
	Price         float64           `json:"price,omitempty"`
	Quantity      float64           `json:"quantity,omitempty"`
	Amount        float64           `json:"amount,omitempty"`
	PaymentMethod ItemPaymentMethod `json:"paymentMethod,omitempty"`
	PaymentObject ItemPaymentObject `json:"paymentObject,omitempty"`
	Tax           ItemTax           `json:"tax,omitempty"`
}
