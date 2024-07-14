package refund

import (
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
		return nil, err
	}

	return refunds, nil
}

func (service *Service) FindOne(id uint) (models.Refund, error) {
	refund, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Refund{}, err
	}

	return refund, nil
}

func (service *Service) Create(input types.CreateRefundInput, user models.User, eventId uint) (models.Refund, error) {
	fromUser, err := service.UserRepository.FindOneById(input.FromUserID)
	if err != nil {
		return models.Refund{}, err
	}

	toUser, err := service.UserRepository.FindOneById(input.ToUserID)
	if err != nil {
		return models.Refund{}, err
	}

	event, err := service.EventService.FindOne(eventId)
	if err != nil {
		return models.Refund{}, err
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
		return models.Refund{}, err
	}

	balances, err := service.EventService.CreateBalances(event)
	if err != nil {
		return models.Refund{}, err
	}

	if _, err := service.EventService.CreateTransactions(event, balances); err != nil {
		return models.Refund{}, err
	}

	return refund, nil
}

func (service *Service) Update(id uint, input types.UpdateRefundInput) (models.Refund, error) {
	refund, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Refund{}, err
	}

	if input.Amount != 0 {
		refund.Amount = input.Amount
	}
	if input.FromUserID != 0 {
		fromUser, err := service.UserRepository.FindOneById(input.FromUserID)
		if err != nil {
			return models.Refund{}, err
		}
		refund.FromUserID = input.FromUserID
		refund.From = fromUser
	}
	if input.ToUserID != 0 {
		toUser, err := service.UserRepository.FindOneById(input.ToUserID)
		if err != nil {
			return models.Refund{}, err
		}
		refund.ToUserID = input.ToUserID
		refund.To = toUser
	}
	if !input.Date.IsZero() {
		refund.Date = input.Date
	}

	event, err := service.EventService.FindOne(refund.EventID)
	if err != nil {
		return models.Refund{}, err
	}

	updatedRefund, err := service.Repository.Update(refund)
	if err != nil {
		return models.Refund{}, err
	}

	balances, err := service.EventService.CreateBalances(event)
	if err != nil {
		return models.Refund{}, err
	}

	if _, err := service.EventService.CreateTransactions(event, balances); err != nil {
		return models.Refund{}, err
	}

	return updatedRefund, nil
}
func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}
