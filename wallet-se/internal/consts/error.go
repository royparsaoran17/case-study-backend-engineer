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

	ErrInvalidUUID    = Error("UUID is not in its proper form")
	ErrWalletDisabled = Error("wallet is disabled")
)
