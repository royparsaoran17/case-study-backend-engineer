package jwt

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrTokenMalformed = Error("token malformed")
	ErrInvalidToken   = Error("invalid token")
	ErrTokenExpired   = Error("token expired")
)

func NewError(err error) Error {
	return err.(Error)
}
