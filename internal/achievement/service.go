package achievement

import (
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
)

type Service struct {
	Repository types.IAchievementRepository
}

func NewService(
	repository types.IAchievementRepository,
) types.IAchievementService {
	return &Service{
		Repository: repository,
	}
}

func (service *Service) Find() ([]models.Achievement, error) {
	achievements, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	return achievements, nil
}

func (service *Service) FindOne(id uint) (models.Achievement, error) {
	achievement, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Achievement{}, err
	}

	return achievement, nil
}

func (service *Service) Create(input types.CreateAchievementInput) (models.Achievement, error) {
	achievement := models.Achievement{
		Name:        input.Name,
		Description: input.Description,
		Points:      input.Points,
		Condition:   input.Condition,
	}

	return service.Repository.Create(achievement)
}

func (service *Service) Update(id uint, input types.UpdateAchievementInput) (models.Achievement, error) {
	achievement, err := service.Repository.FindOne(id)

	if err != nil {
		return models.Achievement{}, err
	}

	if input.Name != "" {
		achievement.Name = input.Name
	}
	if input.Description != "" {
		achievement.Description = input.Description
	}
	if input.Points != 0 {
		achievement.Points = input.Points
	}
	if input.Condition != "" {
		achievement.Condition = input.Condition
	}

	return service.Repository.Update(achievement)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}
