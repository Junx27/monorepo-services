package transaction

import (
	"context"
	"log"
	"time"
	"transaction/pkg/database"

	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	UUID      string    `json:"id" db:"uuid"`
	Status    string    `json:"status" db:"status"`
	UserID    string    `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func StoreTransaction(ctx context.Context, req Transaction) error {
	query := "INSERT INTO transactions (status, user_id) VALUES (@status, @user_id)"
	args := pgx.NamedArgs{
		"status":  req.Status,
		"user_id": req.UserID,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (h *Handler) GetTransactionFromDB(ctx context.Context, id string) (Transaction, error) {
	query := `SELECT uuid, status, user_id, created_at, updated_at 
              FROM transactions 
              WHERE uuid = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := database.DB.Query(ctx, query, args)
	if err != nil {
		log.Println("Error executing query:", err)
		return Transaction{}, err
	}
	defer rows.Close()

	transaction, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (Transaction, error) {
		var t Transaction
		err := row.Scan(&t.UUID, &t.Status, &t.UserID, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return t, err
		}
		return t, nil
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return Transaction{}, pgx.ErrNoRows
		}
		log.Println("Error collecting row:", err)
		return Transaction{}, err
	}

	return transaction, nil
}
