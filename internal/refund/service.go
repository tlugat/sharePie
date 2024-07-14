package refund

import (
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Service struct {
	Repository      types.IRefundRepository
	UserRepository  types.IUserRepository
	EventRepository types.IEventRepository
}

func NewService(
	repository types.IRefundRepository,
	userRepository types.IUserRepository,
	eventRepository types.IEventRepository,
) types.IRefundService {
	return &Service{
		Repository:      repository,
		UserRepository:  userRepository,
		EventRepository: eventRepository,
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

func (service *Service) Create(input types.CreateRefundInput, user models.User) (models.Refund, error) {
	fromUser, err := service.UserRepository.FindOneById(input.FromUserID)
	if err != nil {
		return models.Refund{}, err
	}

	toUser, err := service.UserRepository.FindOneById(input.ToUserID)
	if err != nil {
		return models.Refund{}, err
	}

	event, err := service.EventRepository.FindOne(input.EventID)
	if err != nil {
		return models.Refund{}, err
	}

	refund := models.Refund{
		FromUserID: input.FromUserID,
		From:       fromUser,
		ToUserID:   input.ToUserID,
		To:         toUser,
		Amount:     input.Amount,
		EventID:    input.EventID,
		Event:      event,
		AuthorID:   user.ID,
		Author:     user,
		Date:       input.Date,
	}

	return service.Repository.Create(refund)
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

	return service.Repository.Update(refund)
}
func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}
