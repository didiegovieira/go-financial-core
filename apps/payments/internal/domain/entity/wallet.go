package entity

type Wallet struct {
	ID string
}

func (w *Wallet) NewWallet(id string) *Wallet {
	return &Wallet{
		ID: id,
	}
}
