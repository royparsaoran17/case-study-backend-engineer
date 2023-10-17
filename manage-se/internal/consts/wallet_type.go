package consts

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	TypeDebit  = Type("debit")
	TypeCredit = Type("credit")
)
