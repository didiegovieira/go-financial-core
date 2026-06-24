package entity

import "time"

type PaymentStatus string

const (
	StatusCreated    PaymentStatus = "CREATED"
	StatusProcessing PaymentStatus = "PROCESSING"
	StatusCompleted  PaymentStatus = "COMPLETED"
)

const (
	MethodPix  string = "PIX"
	MethodCard string = "CARD"
)

type Payment struct {
	ID        string
	Amount    float64
	Method    string
	Status    PaymentStatus
	CreatedAt time.Time
}

func (p *Payment) NewPayment(id string, amount float64, method string) *Payment {
	return &Payment{
		ID:        id,
		Amount:    amount,
		Method:    method,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
}
