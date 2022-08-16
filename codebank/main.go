package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/r-avellar/codebank/domain"
	"github.com/r-avellar/codebank/infrastructure/repository"
	"github.com/r-avellar/codebank/usecase"
	"log"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "4789"
	cc.Name = "Renato"
	cc.ExpirationYear = 2025
	cc.ExpirationMonth = 02
	cc.CVV = 456
	cc.Limit = 1000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}
