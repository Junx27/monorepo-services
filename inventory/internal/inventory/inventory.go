package inventory

import (
	"context"
	"inventory/pkg/database"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Inventory struct {
	UUID          string    `json:"id" db:"uuid"`
	UserID        string    `json:"userId" db:"user_id"`
	ProductID     string    `json:"productId" db:"product_id"`
	TransactionID string    `json:"transactionId" db:"transaction_id"`
	Quantity      int       `json:"quantity" db:"quantity"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type InventoryTransaction struct {
	UserID        string `json:"userId"`
	ProductID     string `json:"productId"`
	TransactionID string `json:"transactionId"`
	Quantity      int    `json:"quantity"`
}

// CreateInventory creates a new inventory record
func StoreInventory(ctx context.Context, req Inventory) error {
	// Query untuk insert data inventory baru, kuantitas selalu dimulai dengan 1
	query := `INSERT INTO inventories (user_id, product_id, transaction_id, quantity, created_at, updated_at) 
			  VALUES (@user_id, @product_id, @transaction_id, 1, NOW(), NOW())`
	args := pgx.NamedArgs{
		"user_id":        req.UserID,
		"product_id":     req.ProductID,
		"transaction_id": req.TransactionID,
	}

	// Eksekusi query untuk menyimpan data ke database
	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println("Error inserting inventory:", err)
		return err
	}

	return nil
}

// GetInventory retrieves a single inventory record by its ID
func GetInventory(ctx context.Context, id string) (Inventory, error) {

	query := `SELECT uuid, user_id, product_id, transaction_id, quantity, created_at, updated_at
			  FROM inventories
			  WHERE uuid = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := database.DB.Query(ctx, query, args)
	if err != nil {
		log.Println("Error executing query:", err)
		return Inventory{}, err
	}
	defer rows.Close()

	inventory, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (Inventory, error) {
		var t Inventory
		err := row.Scan(&t.UUID, &t.UserID, &t.ProductID, &t.TransactionID, &t.Quantity, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return t, err
		}
		return t, nil
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return Inventory{}, pgx.ErrNoRows
		}
		log.Println("Error collecting row:", err)
		return Inventory{}, err
	}

	return inventory, nil
}

// UpdateInventory updates an existing inventory record
func UpdateInventory(ctx context.Context, id string) error {
	// Query untuk menambah kuantitas sebanyak 1
	query := `
		UPDATE inventories
		SET quantity = quantity - 1, updated_at = NOW()
		WHERE transaction_id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	// Execute the query to update the inventory
	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println("Error updating inventory:", err)
		return err
	}

	return nil
}
