package dq

const (
	ErrUnexpected             = "Unexpected error, please try again."
	ErrFieldRequired          = "%s is required."
	ErrFieldGTEAndLTE         = "%s must be greater than %d characters and less than %d."
	ErrFieldMustBeArrayOfType = "%s must be an array of %s."
	ErrInvalidJSONBody        = "Invalid JSON body."
	ErrInvalidToken           = "Invalid Token."
	ErrInternal               = "Internal error."
	ErrUnauthorized           = "Unauthorized."
	ErrNotFound               = "%s not found."
	ErrEmailTaken             = "Email is taken."
	UserNotFound              = "A user with these credentials could not be found."
)
