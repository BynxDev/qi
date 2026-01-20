package qi

import "fmt"

// ErrorCode represents an API error code.
type ErrorCode int

const (
	ErrorCodeOrderAlreadyExists                  ErrorCode = 1
	ErrorCodeOrderNotFound                       ErrorCode = 2
	ErrorCodeOrderAlreadyCancelled               ErrorCode = 3
	ErrorCodeNoCompatibleServicesFound           ErrorCode = 4
	ErrorCodeCanNotProcessRequest                ErrorCode = 5
	ErrorCodeRequisitesNotFound                  ErrorCode = 6
	ErrorCodeRequisitesAlreadyExists             ErrorCode = 7
	ErrorCodeCanNotCreateNewRequisites           ErrorCode = 8
	ErrorCodeTerminalNotFoundException           ErrorCode = 9
	ErrorCodePaymentAlreadyExists                ErrorCode = 10
	ErrorCodeMaxNumberOfPaymentsForOrderExceeded ErrorCode = 11
	ErrorCodePaymentNotFound                     ErrorCode = 12
	ErrorCodeUnknownStrategy                     ErrorCode = 13
	ErrorCodeProcessingImpossible                ErrorCode = 14
	ErrorCodeCanNotCancelPayment                 ErrorCode = 15
	ErrorCodeCanNotConfirmPayment                ErrorCode = 16
	ErrorCodeCanNotFinishAuthentication          ErrorCode = 17
	ErrorCodeRefundsNotAllowed                   ErrorCode = 18
	ErrorCodePaymentParamsNotFound               ErrorCode = 19
	ErrorCodeRefundError                         ErrorCode = 20
	ErrorCodeValidationError                     ErrorCode = 21
	ErrorCodeIncorrectPaymentState               ErrorCode = 22
	ErrorCodeInternalSystemError                 ErrorCode = 23
	ErrorCodeExternalSystemError                 ErrorCode = 24
	ErrorCodeInvalidPaymentFormDomain            ErrorCode = 26
	ErrorCodeBadCredentials                      ErrorCode = 27
	ErrorCodeLimitViolation                      ErrorCode = 28
	ErrorCodeTransferNotFound                    ErrorCode = 29
	ErrorCodeIncorrectTransferState              ErrorCode = 30
	ErrorCodeTokenNotFound                       ErrorCode = 31
	ErrorCodeTokenProcessNotAllowed              ErrorCode = 32
	ErrorCodeCanNotCancelTransfer                ErrorCode = 33
	ErrorCodeTransferAlreadyExists               ErrorCode = 34
	ErrorCodeInvalidTokenType                    ErrorCode = 35
)

// ErrorMessage represents an API error message.
type ErrorMessage string

const (
	ErrorMessageOrderAlreadyExists                  ErrorMessage = "ORDER_ALREADY_EXISTS"
	ErrorMessageOrderNotFound                       ErrorMessage = "ORDER_NOT_FOUND"
	ErrorMessageOrderAlreadyCancelled               ErrorMessage = "ORDER_ALREADY_CANCELLED"
	ErrorMessageNoCompatibleServicesFound           ErrorMessage = "NO_COMPATIBLE_SERVICES_FOUND"
	ErrorMessageCanNotProcessRequest                ErrorMessage = "CAN_NOT_PROCESS_REQUEST"
	ErrorMessageRequisitesNotFound                  ErrorMessage = "REQUISITES_NOT_FOUND"
	ErrorMessageRequisitesAlreadyExists             ErrorMessage = "REQUISITES_ALREADY_EXISTS"
	ErrorMessageCanNotCreateNewRequisites           ErrorMessage = "CAN_NOT_CREATE_NEW_REQUISITES"
	ErrorMessageTerminalNotFoundException           ErrorMessage = "TERMINAL_NOT_FOUND_EXCEPTION"
	ErrorMessagePaymentAlreadyExists                ErrorMessage = "PAYMENT_ALREADY_EXISTS"
	ErrorMessageMaxNumberOfPaymentsForOrderExceeded ErrorMessage = "MAX_NUMBER_OF_PAYMENTS_FOR_ORDER_EXCEEDED"
	ErrorMessagePaymentNotFound                     ErrorMessage = "PAYMENT_NOT_FOUND"
	ErrorMessageUnknownStrategy                     ErrorMessage = "UNKNOWN_STRATEGY"
	ErrorMessageProcessingImpossible                ErrorMessage = "PROCESSING_IMPOSSIBLE"
	ErrorMessageCanNotCancelPayment                 ErrorMessage = "CAN_NOT_CANCEL_PAYMENT"
	ErrorMessageCanNotConfirmPayment                ErrorMessage = "CAN_NOT_CONFIRM_PAYMENT"
	ErrorMessageCanNotFinishAuthentication          ErrorMessage = "CAN_NOT_FINISH_AUTHENTICATION"
	ErrorMessageRefundsNotAllowed                   ErrorMessage = "REFUNDS_NOT_ALLOWED"
	ErrorMessagePaymentParamsNotFound               ErrorMessage = "PAYMENT_PARAMS_NOT_FOUND"
	ErrorMessageRefundError                         ErrorMessage = "REFUND_ERROR"
	ErrorMessageValidationError                     ErrorMessage = "VALIDATION_ERROR"
	ErrorMessageIncorrectPaymentState               ErrorMessage = "INCORRECT_PAYMENT_STATE"
	ErrorMessageInternalSystemError                 ErrorMessage = "INTERNAL_SYSTEM_ERROR"
	ErrorMessageExternalSystemError                 ErrorMessage = "EXTERNAL_SYSTEM_ERROR"
	ErrorMessageInvalidPaymentFormDomain            ErrorMessage = "INVALID_PAYMENT_FORM_DOMAIN"
	ErrorMessageBadCredentials                      ErrorMessage = "BAD_CREDENTIALS"
	ErrorMessageLimitViolation                      ErrorMessage = "LIMIT_VIOLATION"
	ErrorMessageTransferNotFound                    ErrorMessage = "TRANSFER_NOT_FOUND"
	ErrorMessageIncorrectTransferState              ErrorMessage = "INCORRECT_TRANSFER_STATE"
	ErrorMessageTokenNotFound                       ErrorMessage = "TOKEN_NOT_FOUND"
	ErrorMessageTokenProcessNotAllowed              ErrorMessage = "TOKEN_PROCESS_NOT_ALLOWED"
	ErrorMessageCanNotCancelTransfer                ErrorMessage = "CAN_NOT_CANCEL_TRANSFER"
	ErrorMessageTransferAlreadyExists               ErrorMessage = "TRANSFER_ALREADY_EXISTS"
	ErrorMessageInvalidTokenType                    ErrorMessage = "INVALID_TOKEN_TYPE"
)

// Error represents an API error response.
type Error struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails contains the error code and message.
type ErrorDetails struct {
	Code    ErrorCode    `json:"code"`
	Message ErrorMessage `json:"message"`
}

// APIError represents an error returned by the API.
type APIError struct {
	StatusCode int
	Message    string
	Err        *Error
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("API error (status %d): code=%d, message=%s",
			e.StatusCode, e.Err.Error.Code, e.Err.Error.Message)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a not found error.
func (e *APIError) IsNotFound() bool {
	if e.Err != nil {
		return e.Err.Error.Code == ErrorCodePaymentNotFound ||
			e.Err.Error.Code == ErrorCodeOrderNotFound ||
			e.Err.Error.Code == ErrorCodeTransferNotFound ||
			e.Err.Error.Code == ErrorCodeTokenNotFound
	}
	return e.StatusCode == 404
}

// IsValidationError returns true if the error is a validation error.
func (e *APIError) IsValidationError() bool {
	if e.Err != nil {
		return e.Err.Error.Code == ErrorCodeValidationError
	}
	return e.StatusCode == 400
}

// IsAuthenticationError returns true if the error is an authentication error.
func (e *APIError) IsAuthenticationError() bool {
	if e.Err != nil {
		return e.Err.Error.Code == ErrorCodeBadCredentials
	}
	return e.StatusCode == 401
}
