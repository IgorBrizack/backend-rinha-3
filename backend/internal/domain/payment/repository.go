package payment

import (
	"time"
)

type Repository interface {
	CreatePayment(payment Payment) error
	GetPayments(from, to *time.Time) ([]Payment, error)
}
