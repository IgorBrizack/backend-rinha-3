package database

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) payment.Repository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, p payment.Payment) error {
	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) GetPayments(ctx context.Context, from, to *time.Time) ([]payment.Payment, error) {
	var payments []payment.Payment

	query := r.db.WithContext(ctx).Model(&payment.Payment{})
	if from != nil {
		query = query.Where("created_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("created_at <= ?", *to)
	}

	if err := query.Order("created_at ASC").Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}

	return payments, nil
}

func (r *paymentRepository) PurgePayments(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Exec("TRUNCATE TABLE payments RESTART IDENTITY CASCADE").Error; err != nil {
		return fmt.Errorf("failed to purge payments: %w", err)
	}
	return nil
}
