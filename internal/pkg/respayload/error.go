package respayload

type ErrorCode string

// Follow this rule: first number is the entity number:
// System Level Error = 0
// User = 1
// Tax = 2
// after that, follow the underscore and the sequence number of the error code.
// This to make grouping and debugging error much easier.
const (
	ErrorGeneralValidationError ErrorCode = "0_0001"
	ErrorBindingBodyRequest     ErrorCode = "0_0002"

	ErrorCodeUserCantBeCreated  ErrorCode = "1_0001"
	ErrorCodeUserCantBeFound    ErrorCode = "1_0002"
	ErrorCodeUserWrongPassword  ErrorCode = "1_0003"
	ErrorCodeUserWrongAuthToken ErrorCode = "1_0004"

	ErrorCodeTaxCantBeCreated ErrorCode = "2_0001"
	ErrorCodeTaxDBError       ErrorCode = "2_0002"
)

type Error struct {
	HttpStatusCode int       `json:"http_status_code"` // net/http status error code
	ErrorCode      ErrorCode `json:"error_code"`
	Message        string    `json:"message"`
}
