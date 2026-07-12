package wallet

type CreditInput struct {
	Phone  string `json:"phone" binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}

type CreditOutput struct {
	Balance int64 `json:"balance"`
}
