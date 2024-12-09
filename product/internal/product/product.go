package product

import (
	"context"
	"log"
	"product/pkg/database"
	"time"

	"github.com/jackc/pgx/v5"
)

type Product struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     int       `json:"price" db:"price"`
	Stock     int       `json:"stock" db:"stock"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func GetProduct(ctx context.Context, id string) (Product, error) {
	query := "SELECT id, name, price, stock, created_at FROM products WHERE id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := database.DB.Query(ctx, query, args)
	if err != nil {
		log.Println(err.Error())
		return Product{}, err
	}

	defer rows.Close()

	product, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (Product, error) {
		var p Product
		err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)
		if err != nil {
			log.Println(err.Error())
			return p, err
		}
		return p, nil
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return Product{}, pgx.ErrNoRows
		}

		log.Println(err.Error())
		return Product{}, err
	}

	return product, nil
}

func StoreProduct(ctx context.Context, req Product) error {
	query := "INSERT INTO products (name, price, stock) VALUES (@name, @price, @stock)"
	args := pgx.NamedArgs{
		"name":  req.Name,
		"price": req.Price,
		"stock": req.Stock,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func UpdateProduct(ctx context.Context, product Product) error {
	query := "UPDATE products SET name=@name, price=@price, stock=@stock WHERE id = @id"
	args := pgx.NamedArgs{
		"id":    product.ID,
		"name":  product.Name,
		"price": product.Price,
		"stock": product.Stock,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func IncreaseStock(ctx context.Context, productID string) error {
	// Increase stock by 1
	query := "UPDATE products SET stock = stock + 1 WHERE id = @id"
	args := pgx.NamedArgs{
		"id": productID,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func ReduceStock(ctx context.Context, productID string) error {
	// Reduce stock by 1
	query := "UPDATE products SET stock = stock - 1 WHERE id = @id AND stock > 0"
	args := pgx.NamedArgs{
		"id": productID,
	}

	_, err := database.DB.Exec(ctx, query, args)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
