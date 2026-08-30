package services

import (
	"concurrency-simulator/services/shared/topic_messages"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AccountService struct {
	log    *zap.Logger
	driver *sql.DB
}

func (ac *AccountService) Execute(payment topic_messages.Payment) bool {
	tx, err := ac.driver.Begin()
	if err != nil {
		ac.log.Error("failed to ping database", zap.Error(err))
		return false
	}

	is_account_created := ac.accountExists(tx, payment.Email)

	if is_account_created {
		ac.log.Info("Account already exists", zap.String("email", payment.Email))
		return is_account_created
	}

	ac.log.Info("Account not found, creating account", zap.String("email", payment.Email))
	is_account_created, err = ac.createAccount(tx, payment)

	if err := tx.Commit(); err != nil {
		ac.log.Error("failed to commit transaction", zap.Error(err))
		return false
	}

	return is_account_created
}

func (ac *AccountService) accountExists(tx *sql.Tx, email string) bool {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM account
			WHERE email = $1
		)
	`

	var exists bool
	err := tx.QueryRow(query, email).Scan(&exists)
	if err != nil {
		ac.log.Error("failed to check account exists", zap.Error(err), zap.String("email", email))
		return false
	}

	ac.log.Info("QUERY COMPLETED", zap.String("email", email), zap.Bool("exists", exists))
	return exists
}

func (ac *AccountService) createAccount(tx *sql.Tx, data topic_messages.Payment) (bool, error) {
	account_creation_query := `
		INSERT INTO account (first_name, last_name, email, created_at) 
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	var account_id string
	var created_at time.Time
	account_creation_error := tx.QueryRow(
		account_creation_query,
		data.FirstName, data.LastName, data.Email,
	).Scan(&account_id, &created_at)

	if account_creation_error != nil {
		ac.log.Error("failed to create account", zap.Error(account_creation_error), zap.String("email", data.Email))
		tx.Rollback()
		return false, account_creation_error
	}

	account_outbox_event_query := `
		INSERT INTO outbox (id, domain_id, payload, entity, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	event_id, _ := uuid.NewV7()
	type EventData struct {
		Status    string    `json:"status"`
		FistName  string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}

	meta := EventData{
		Status:    "CREATED",
		FistName:  data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		CreatedAt: created_at,
	}

	event_data, _ := json.Marshal(meta)

	_, outbox_event_error := tx.Exec(account_outbox_event_query, event_id, account_id, event_data, "ACCOUNT", "PENDING")

	if outbox_event_error != nil {
		ac.log.Error("failed to create account event", zap.Error(outbox_event_error), zap.String("email", data.Email))
		tx.Rollback()
		return false, outbox_event_error
	}

	return true, nil
}
