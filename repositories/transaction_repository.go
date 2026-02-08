package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// initiate subtotal => jumlah total transaksi
	totalAmount := 0
	// initiate model tx detail
	details := make([]models.TransactionDetail, 0)
	// loop item
	for _, item := range items {
		var productName string
		var productID, price, stock int

		err := tx.QueryRow("SELECT id, name, price, stock FROM products where id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product ID %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		// calculate current total = quantity * price
		// sum into subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// deduct stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id=$2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		// insert into trxDetail
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert trx
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	fmt.Println(details)
	err = fmt.Errorf("cek")

	// insert trxDetail
	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}
	return res, nil
}
