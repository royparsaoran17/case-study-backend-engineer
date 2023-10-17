package consts

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrRoleAlreadyExist = Error("role already exist")
	ErrRoleNotFound     = Error("role not found")

	ErrWrongPassword     = Error("wrong password")
	ErrPhoneAlreadyExist = Error("phone already exist")
	ErrCustomerNotFound  = Error("customer not found")

	ErrInvalidUUID = Error("UUID is not in its proper form")
)
