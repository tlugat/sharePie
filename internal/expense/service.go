package expense

import (
	"errors"
	"fmt"
	models2 "sharePie-api/internal/models"
	"sharePie-api/internal/participant"
	"sharePie-api/internal/payer"
	"sharePie-api/internal/types"
	"sharePie-api/pkg/config/thirdparty/cloudinary"
)

type Service struct {
	Repository            types.IExpenseRepository
	TagRepository         types.ITagRepository
	UserRepository        types.IUserRepository
	ParticipantRepository participant.IParticipantRepository
	PayerRepository       payer.IPayerRepository
	EventRepository       types.IEventRepository
}

func NewService(
	repository types.IExpenseRepository,
	tagRepository types.ITagRepository,
	userRepository types.IUserRepository,
	participantRepository participant.IParticipantRepository,
	payerRepository payer.IPayerRepository,
	eventRepository types.IEventRepository,
) types.IExpenseService {
	return &Service{
		Repository:            repository,
		TagRepository:         tagRepository,
		UserRepository:        userRepository,
		ParticipantRepository: participantRepository,
		PayerRepository:       payerRepository,
		EventRepository:       eventRepository,
	}
}

func (service *Service) Find() ([]models2.Expense, error) {
	expenses, err := service.Repository.Find()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (service *Service) FindOne(id uint) (models2.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models2.Expense{}, err
	}

	return expense, nil
}

func (service *Service) FindByEventId(id uint) ([]models2.Expense, error) {
	expenses, err := service.Repository.FindByEventId(id)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (service *Service) Create(input types.CreateExpenseInput, user models2.User) (models2.Expense, error) {
	expense := models2.Expense{
		Name:        input.Name,
		Description: input.Description,
		TagID:       input.Tag,
		Amount:      input.Amount,
		AuthorID:    user.ID,
		EventID:     input.Event,
	}

	participants, err := service.handleParticipants(input.Participants, expense)
	if err != nil {
		return models2.Expense{}, err
	}
	expense.Participants = participants

	payers, err := service.handlePayers(input.Payers, expense)
	if err != nil {
		return models2.Expense{}, err
	}
	expense.Payers = payers

	if input.Image != "" {
		image, err := cloudinary.UploadImage(input.Image, "Events")
		if err != nil {
			return models2.Expense{}, err
		}
		expense.Image = image
	}

	if _, err := service.Repository.Create(expense); err != nil {
		return models2.Expense{}, err
	}

	event, err := service.EventRepository.FindOne(input.Event)
	if err != nil {
		return models2.Expense{}, err
	}

	balances, err := service.createBalances(event)
	if err != nil {
		return models2.Expense{}, err
	}

	if _, err := service.createTransactions(event, balances); err != nil {
		return models2.Expense{}, err
	}

	return expense, nil
}

func (service *Service) Update(id uint, input types.UpdateExpenseInput) (models2.Expense, error) {
	expense, err := service.Repository.FindOne(id)

	if err != nil {
		return models2.Expense{}, err
	}

	if input.Name != "" {
		expense.Name = input.Name
	}

	if input.Description != "" {
		expense.Description = input.Description
	}

	if input.Tag != 0 {
		expense.TagID = input.Tag
	}

	if input.Image != "" {
		image, err := cloudinary.UploadImage(input.Image, "Events")
		if err != nil {
			return models2.Expense{}, err
		}
		expense.Image = image
	}

	if input.Amount != 0 {
		expense.Amount = input.Amount
	}

	if input.Participants != nil {
		updatedParticipants, err := service.handleParticipants(input.Participants, expense)
		if err != nil {
			return models2.Expense{}, err
		}
		err = service.ParticipantRepository.DeleteByExpenseId(expense.ID)
		if err != nil {
			return models2.Expense{}, err
		}
		expense.Participants = updatedParticipants
	}

	if input.Payers != nil {
		updatedPayers, err := service.handlePayers(input.Payers, expense)
		if err != nil {
			return models2.Expense{}, err
		}
		err = service.PayerRepository.DeleteByExpenseID(expense.ID)
		if err != nil {
			return models2.Expense{}, err
		}
		expense.Payers = updatedPayers
	}

	event, err := service.EventRepository.FindOne(id)
	if err != nil {
		return models2.Expense{}, err
	}

	balances, err := service.createBalances(event)
	if err != nil {
		return models2.Expense{}, err
	}

	if _, err := service.createTransactions(event, balances); err != nil {
		return models2.Expense{}, err
	}

	return service.Repository.Update(expense)
}

func (service *Service) Delete(id uint) error {
	return service.Repository.Delete(id)
}

func (service *Service) handleParticipants(participantsInput []types.ParticipantInput, expense models2.Expense) ([]models2.Participant, error) {
	var participants []models2.Participant
	userIDSet := make(map[uint]bool)

	for _, p := range participantsInput {
		if _, exists := userIDSet[p.Id]; exists {
			return nil, fmt.Errorf("duplicate user ID found: %d", p.Id)
		}
		participantUser, err := service.UserRepository.FindOneById(p.Id)
		if err != nil {
			return nil, err
		}
		userIDSet[p.Id] = true
		participants = append(participants, models2.Participant{
			UserID: participantUser.ID,
			User:   participantUser,
			Amount: p.Amount,
		})
	}

	err := service.validateParticipants(participants, expense.EventID)
	if err != nil {
		return nil, err
	}

	err = validateDueAmounts(expense.Amount, participants)
	if err != nil {
		return nil, err

	}

	return participants, nil
}

func (service *Service) handlePayers(payersInput []types.PayerInput, expense models2.Expense) ([]models2.Payer, error) {
	var payers []models2.Payer
	userIDSet := make(map[uint]bool)
	for _, p := range payersInput {
		if _, exists := userIDSet[p.Id]; exists {
			return nil, fmt.Errorf("duplicate user ID found: %d", p.Id)
		}
		payerUser, err := service.UserRepository.FindOneById(p.Id)
		if err != nil {
			return nil, err
		}
		userIDSet[p.Id] = true
		payers = append(payers, models2.Payer{
			UserID: payerUser.ID,
			User:   payerUser,
			Amount: p.Amount,
		})
	}
	err := service.validatePayers(payers, expense.EventID)
	if err != nil {
		return nil, err
	}
	err = validatePaidAmounts(expense.Amount, payers)
	if err != nil {
		return nil, err
	}
	return payers, nil
}

func validatePaidAmounts(amount float64, payers []models2.Payer) error {
	var totalPaid = 0.00

	for _, p := range payers {
		totalPaid += p.Amount
	}

	if totalPaid != amount {
		return errors.New("total paid amount does not match the expense amount")
	}

	return nil
}

func validateDueAmounts(amount float64, participants []models2.Participant) error {
	var totalDue float64

	for _, p := range participants {
		totalDue += p.Amount
	}

	if totalDue != amount {
		return errors.New("total due amount does not match the expense amount")
	}

	return nil
}

func (service *Service) validateParticipants(participants []models2.Participant, eventID uint) error {
	eventUsers, err := service.EventRepository.FindUsers(eventID)
	if err != nil {
		return err
	}

	userMap := make(map[uint]bool)
	for _, u := range eventUsers {
		userMap[u.ID] = true
	}

	for _, p := range participants {
		if _, exists := userMap[p.UserID]; !exists {
			return errors.New("one or more participants are not associated with this event")
		}
	}

	return nil
}

func (service *Service) validatePayers(payers []models2.Payer, eventID uint) error {
	eventUsers, err := service.EventRepository.FindUsers(eventID)
	if err != nil {
		return err
	}

	userMap := make(map[uint]bool)
	for _, u := range eventUsers {
		userMap[u.ID] = true
	}

	for _, p := range payers {
		if _, exists := userMap[p.UserID]; !exists {
			return errors.New("un ou plusieurs payeurs ne sont pas associés à cet événement")
		}
	}

	return nil
}

func (service *Service) createBalances(event models2.Event) ([]models2.Balance, error) {
	expenses, err := service.Repository.FindByEventId(event.ID)
	if err != nil {
		return nil, err
	}

	eventUsers, err := service.EventRepository.FindUsers(event.ID)

	if err != nil {
		return nil, err
	}

	userBalances := make(map[uint]float64)
	userDetails := make(map[uint]models2.User)

	for _, eventUser := range eventUsers {
		userBalances[eventUser.ID] = 0
		userDetails[eventUser.ID] = eventUser
	}

	var total float64

	for _, eventExpense := range expenses {
		total += eventExpense.Amount
		for _, expenseParticipant := range eventExpense.Participants {
			userBalances[expenseParticipant.UserID] -= expenseParticipant.Amount
		}
		for _, expensePayer := range eventExpense.Payers {
			userBalances[expensePayer.UserID] += expensePayer.Amount
		}
	}

	var balances []models2.Balance
	for userID, userBalance := range userBalances {
		balance := models2.Balance{
			UserID:  userID,
			User:    userDetails[userID],
			Amount:  userBalance,
			EventID: event.ID,
		}
		balances = append(balances, balance)
	}
	err = service.EventRepository.DeleteBalances(event)
	if err != nil {
		return nil, err
	}
	err = service.EventRepository.CreateBalances(balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func (service *Service) createTransactions(event models2.Event, balances []models2.Balance) ([]models2.Transaction, error) {

	var creditors []models2.Balance
	var debtors []models2.Balance

	for _, balance := range balances {
		if balance.Amount > 0 {
			creditors = append(creditors, balance)
		} else if balance.Amount < 0 {
			debtors = append(debtors, balance)
		}
	}

	var transactions []models2.Transaction

	for len(creditors) > 0 && len(debtors) > 0 {
		creditor := creditors[0]
		debtor := debtors[0]

		amount := min(creditor.Amount, -debtor.Amount)

		transaction := models2.Transaction{
			FromUserID: debtor.UserID,
			From:       debtor.User,
			ToUserID:   creditor.UserID,
			To:         creditor.User,
			EventID:    event.ID,
			Completed:  false,
			Amount:     amount,
		}
		transactions = append(transactions, transaction)

		creditor.Amount -= amount
		debtor.Amount += amount

		if creditor.Amount == 0 {
			creditors = creditors[1:]
		} else {
			creditors[0] = creditor
		}

		if debtor.Amount == 0 {
			debtors = debtors[1:]
		} else {
			debtors[0] = debtor
		}
	}

	err := service.EventRepository.DeleteTransactions(event)
	if err != nil {
		return nil, err
	}
	err = service.EventRepository.CreateTransactions(transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
