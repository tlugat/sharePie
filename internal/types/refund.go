package types

import (
	"sharePie-api/internal/models"
	"time"
)

type IRefundRepository interface {
	Find() ([]models.Refund, error)
	FindOne(id uint) (models.Refund, error)
	FindByEventId(eventId uint) ([]models.Refund, error)
	Create(refund models.Refund) (models.Refund, error)
	Update(refund models.Refund) (models.Refund, error)
	Delete(id uint) error
}

type IRefundService interface {
	Find() ([]models.Refund, error)
	FindOne(id uint) (models.Refund, error)
	FindByEventId(eventId uint) ([]models.Refund, error)
	Create(input CreateRefundInput, user models.User, eventId uint) (models.Refund, error)
	Update(id uint, input UpdateRefundInput) (models.Refund, error)
	Delete(id uint) error
}

type CreateRefundInput struct {
	FromUserID uint      `json:"fromUserId" binding:"required"`
	ToUserID   uint      `json:"toUserId" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Date       time.Time `json:"date" time_format:"2006-01-02T15:04:05Z07:00" binding:"required"`
}

type UpdateRefundInput struct {
	Amount     float64   `json:"amount"`
	FromUserID uint      `json:"fromUserId"`
	ToUserID   uint      `json:"toUserId"`
	Date       time.Time `json:"date" time_format:"2006-01-02T15:04:05Z07:00"`
	ID         uint      `json:"id" binding:"required"`
}
