package domainerror

import "errors"

var (
	ErrDebitNotAllowed   = errors.New("wallet is inactive or has insufficient balance")
	ErrWalletUnavailable = errors.New("wallet does not exist or is inactive")
)
