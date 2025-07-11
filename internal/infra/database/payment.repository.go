package database

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
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

func (r *paymentRepository) GetSummary() error {
	return r.db.Find(&dto.PaymentSummaryResponse{}).Error
}
