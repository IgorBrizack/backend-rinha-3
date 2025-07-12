package database

import (
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) payment.Repository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePayment(p payment.Payment) error {
	return r.db.Create(&p).Error
}

func (r *paymentRepository) GetPayments(from, to *time.Time) ([]payment.Payment, error) {
	var payments []payment.Payment

	query := r.db

	if from != nil && to != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *from, *to)
	} else if from != nil {
		query = query.Where("created_at >= ?", *from)
	} else if to != nil {
		query = query.Where("created_at <= ?", *to)
	}

	err := query.Find(&payments).Error
	return payments, err
}
