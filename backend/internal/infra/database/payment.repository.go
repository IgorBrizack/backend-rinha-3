package database

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *paymentRepository) GetPayments(ctx context.Context, from, to *time.Time) (dto.PaymentSummaryResponse, error) {
	var out dto.PaymentSummaryResponse

	// Estrutura "achatada" para receber os quatro campos em uma linha
	type flat struct {
		DefaultRequests  int     `gorm:"column:default_requests"`
		DefaultAmount    float64 `gorm:"column:default_amount"`
		FallbackRequests int     `gorm:"column:fallback_requests"`
		FallbackAmount   float64 `gorm:"column:fallback_amount"`
	}

	var row flat

	q := r.db.WithContext(ctx).Model(&payment.Payment{})

	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}

	err := q.Select(`
		COUNT(*) FILTER (WHERE "default" = true)  AS default_requests,
		COALESCE(SUM(amount) FILTER (WHERE "default" = true), 0)  AS default_amount,
		COUNT(*) FILTER (WHERE "default" = false) AS fallback_requests,
		COALESCE(SUM(amount) FILTER (WHERE "default" = false), 0) AS fallback_amount
	`).Scan(&row).Error
	if err != nil {
		return out, err
	}

	// Mapear para seu DTO final
	out.Default.TotalRequests = row.DefaultRequests
	out.Default.TotalAmount = row.DefaultAmount
	out.Fallback.TotalRequests = row.FallbackRequests
	out.Fallback.TotalAmount = row.FallbackAmount

	return out, nil
}

func (r *paymentRepository) PurgePayments(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Exec("TRUNCATE TABLE payments RESTART IDENTITY CASCADE").Error; err != nil {
		return fmt.Errorf("failed to purge payments: %w", err)
	}
	return nil
}

func (r *paymentRepository) CreatePaymentsBatch(ctx context.Context, payments []payment.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	// Usando ON CONFLICT para upsert pelo campo correlation_id
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "correlation_id"}}, // chave única
			DoUpdates: clause.AssignmentColumns([]string{"amount", "default", "created_at"}),
		}).
		Create(&payments).Error; err != nil {
		return fmt.Errorf("failed to insert payments batch: %w", err)
	}
	return nil
}
