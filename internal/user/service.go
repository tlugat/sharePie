package user

import (
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
	"sharePie-api/pkg/utils"
)

type Service struct {
	Repository       types.IUserRepository
	AvatarRepository types.IAvatarRepository
}

func NewService(repository types.IUserRepository, avatarRepository types.IAvatarRepository) types.IUserService {
	return &Service{Repository: repository, AvatarRepository: avatarRepository}
}

func (service *Service) Find() ([]models.User, error) {
	return service.Repository.Find()
}

func (service *Service) FindOneById(id uint) (models.User, error) {
	return service.Repository.FindOneById(id)
}

func (service *Service) FindOneByEmail(email string) (models.User, error) {
	return service.Repository.FindOneByEmail(email)
}

func (service *Service) Create(input types.CreateUserInput) (models.User, error) {
	defaultAvatar, err := service.AvatarRepository.FindOne(25)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Role:     utils.UserRole,
		AvatarID: defaultAvatar.ID,
		Avatar:   defaultAvatar,
	}
	return service.Repository.Create(user)
}

func (service *Service) Update(id uint, input types.UpdateUserInput) (models.User, error) {
	user, err := service.Repository.FindOneById(id)

	if err != nil {
		return models.User{}, err
	}

	if input.Username != "" {
		user.Username = input.Username
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Avatar != 0 {
		avatar, err := service.AvatarRepository.FindOne(input.Avatar)
		if err != nil {
			return models.User{}, err
		}
		user.AvatarID = input.Avatar
		user.Avatar = avatar
	}

	return service.Repository.Update(user)
}

func (service *Service) UpdateFirebaseToken(id uint, input types.UpdateUserFirebaseTokenInput) (models.User, error) {
	user, err := service.Repository.FindOneById(id)

	if err != nil {
		return models.User{}, err
	}

	if input.FirebaseToken != "" {
		user.FirebaseToken = &input.FirebaseToken
	}

	return service.Repository.Update(user)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}
