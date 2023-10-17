package consts

type Status string

func (e Status) String() string {
	return string(e)
}

const (
	StatusTransactionSuccess = Status("success")
	StatusTransactionFailed  = Status("failed")
	StatusWalletDisabled     = Status("disabled")
	StatusWalletEnabled      = Status("enabled")
)
