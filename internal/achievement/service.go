package achievement

import (
	"errors"
	"fmt"
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
		return nil, errors.New(fmt.Sprintf("failed to find achievements: %v", err))
	}
	return achievements, nil
}

func (service *Service) FindOne(id uint) (models.Achievement, error) {
	achievement, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Achievement{}, errors.New(fmt.Sprintf("failed to find achievement with id %d: %v", id, err))
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

	newAchievement, err := service.Repository.Create(achievement)
	if err != nil {
		return models.Achievement{}, errors.New(fmt.Sprintf("failed to create achievement: %v", err))
	}
	return newAchievement, nil
}

func (service *Service) Update(id uint, input types.UpdateAchievementInput) (models.Achievement, error) {
	achievement, err := service.Repository.FindOne(id)
	if err != nil {
		return models.Achievement{}, errors.New(fmt.Sprintf("failed to find achievement with id %d: %v", id, err))
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

	updatedAchievement, err := service.Repository.Update(achievement)
	if err != nil {
		return models.Achievement{}, errors.New(fmt.Sprintf("failed to update achievement with id %d: %v", id, err))
	}
	return updatedAchievement, nil
}

func (service *Service) Delete(id uint) error {
	err := service.Repository.Delete(id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete achievement with id %d: %v", id, err))
	}
	return nil
}
