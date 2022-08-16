package usecase

import (
	"github.com/r-avellar/codebank/domain"
	"github.com/r-avellar/codebank/dto"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDTO dto.Transaction) (domain.Transaction, error) {
	creditCard := u.hydrateCreditCard(transactionDTO)
	ccBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}
	creditCard.ID = ccBalanceAndLimit.ID
	creditCard.Limit = ccBalanceAndLimit.Limit
	creditCard.Balance = ccBalanceAndLimit.Balance

	t := u.newTransaction(transactionDTO, ccBalanceAndLimit)
	t.ProcessAndValidate(creditCard)

	err = u.TransactionRepository.SaveTransaction(*t, *creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}
	return *t, nil
}

func (u UseCaseTransaction) hydrateCreditCard(transactionDTO dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDTO.Name
	creditCard.Number = transactionDTO.Number
	creditCard.ExpirationMonth = transactionDTO.ExpirationMonth
	creditCard.ExpirationYear = transactionDTO.ExpirationYear
	creditCard.CVV = transactionDTO.CVV
	return creditCard
}

func (u UseCaseTransaction) newTransaction(transaction dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()
	t.Amount = transaction.Amount
	t.Store = transaction.Store
	t.Description = transaction.Description
	t.CreatedAt = transaction.CreatedAt
	return t
}
