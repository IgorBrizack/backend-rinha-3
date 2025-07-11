package payment

type Repository interface {
	CreatePayment(payment Payment) error
	GetSummary() error
}
