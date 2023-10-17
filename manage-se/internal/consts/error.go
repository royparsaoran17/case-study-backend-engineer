package consts

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrWalletAlreadyExist = Error("wallet already exist")
	ErrWalletNotFound     = Error("wallet not found")

	ErrTransactionNotFound     = Error("transaction not found")
	ErrTransactionAlreadyExist = Error("transaction already exist")

	ErrInvalidUUID = Error("UUID is not in its proper form")

	ErrUserUnauthorized       = Error("unauthorized")
	ErrWalletAlreadyEnabled   = Error("wallet already enabled")
	ErrWalletDisabled         = Error("wallet is disabled")
	ErrForbidden              = Error("forbidden")
	ErrUnauthorized           = Error("unauthorized")
	ErrBearerTokenNotProvided = Error("token not provided")
)
