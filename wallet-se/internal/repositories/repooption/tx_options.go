package repooption

import (
	"github.com/jmoiron/sqlx"
)

type TxOptions struct {
	Tx              *sqlx.Tx
	NotCommitInRepo bool
}

type TxOption func(*TxOptions)

func WithTx(tx *sqlx.Tx) TxOption {
	return func(options *TxOptions) {
		options.Tx = tx
		options.NotCommitInRepo = true
	}
}
