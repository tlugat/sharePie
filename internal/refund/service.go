package refund

import (
	"errors"
	"fmt"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Service struct {
	Repository     types.IRefundRepository
	UserRepository types.IUserRepository
	EventService   types.IEventService
}

func NewService(
	repository types.IRefundRepository,
	userRepository types.IUserRepository,
	eventService types.IEventService,
) types.IRefundService {
	return &Service{
		Repository:     repository,
		UserRepository: userRepository,
		EventService:   eventService,
	}
}

func (service *Service) Find() ([]models.Refund, error) {
	refunds, err := service.Repository.Find()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to find refunds: %v", err))
	}

	return refunds, nil
}

func (service *Service) FindOne(id uint) (models.Refund, error) {
	refund, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to find refund with id %d: %v", id, err))
	}

	return refund, nil
}

func (service *Service) FindByEventId(eventId uint) ([]models.Refund, error) {
	refunds, err := service.Repository.FindByEventId(eventId)
	if err != nil {
		return nil, err
	}

	return refunds, nil
}

func (service *Service) Create(input types.CreateRefundInput, user models.User, eventId uint) (models.Refund, error) {
	fromUser, err := service.UserRepository.FindOneById(input.FromUserID)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", input.FromUserID, err))
	}

	toUser, err := service.UserRepository.FindOneById(input.ToUserID)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", input.ToUserID, err))
	}

	event, err := service.EventService.FindOne(eventId)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to find event with id %d: %v", eventId, err))
	}

	newRefund := models.Refund{
		FromUserID: input.FromUserID,
		From:       fromUser,
		ToUserID:   input.ToUserID,
		To:         toUser,
		Amount:     input.Amount,
		EventID:    eventId,
		Event:      event,
		AuthorID:   user.ID,
		Author:     user,
		Date:       input.Date,
	}

	refund, err := service.Repository.Create(newRefund)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to create refund: %v", err))
	}

	balances, err := service.EventService.CreateBalances(event)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to create balances for event with id %d: %v", event.ID, err))
	}

	if _, err := service.EventService.CreateTransactions(event, balances); err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to create transactions for event with id %d: %v", event.ID, err))
	}

	return refund, nil
}

func (service *Service) Update(id uint, input types.UpdateRefundInput) (models.Refund, error) {
	refund, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to find refund with id %d: %v", id, err))
	}

	if input.Amount != 0 {
		refund.Amount = input.Amount
	}
	if input.FromUserID != 0 {
		fromUser, err := service.UserRepository.FindOneById(input.FromUserID)
		if err != nil {
			return models.Refund{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", input.FromUserID, err))
		}
		refund.FromUserID = input.FromUserID
		refund.From = fromUser
	}
	if input.ToUserID != 0 {
		toUser, err := service.UserRepository.FindOneById(input.ToUserID)
		if err != nil {
			return models.Refund{}, errors.New(fmt.Sprintf("failed to find user with id %d: %v", input.ToUserID, err))
		}
		refund.ToUserID = input.ToUserID
		refund.To = toUser
	}
	if !input.Date.IsZero() {
		refund.Date = input.Date
	}

	event, err := service.EventService.FindOne(refund.EventID)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to find event with id %d: %v", refund.EventID, err))
	}

	updatedRefund, err := service.Repository.Update(refund)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to update refund with id %d: %v", id, err))
	}

	balances, err := service.EventService.CreateBalances(event)
	if err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to create balances for event with id %d: %v", event.ID, err))
	}

	if _, err := service.EventService.CreateTransactions(event, balances); err != nil {
		return models.Refund{}, errors.New(fmt.Sprintf("failed to create transactions for event with id %d: %v", event.ID, err))
	}

	return updatedRefund, nil
}
func (service *Service) Delete(id uint) error {
	err := service.Repository.Delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete refund with id %d: %v", id, err))
	}
	return nil
}
