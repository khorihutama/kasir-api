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
	var transactionDate sql.NullTime
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID, created_at", totalAmount).Scan(&transactionID, &transactionDate)
	if err != nil {
		return nil, err
	}

	// base query
	queryTrxDetail := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
	valueTrxDetail := []interface{}{}

	// loop to generate the query string
	for i, detail := range details {
		details[i].TransactionID = transactionID
		queryTrxDetail += fmt.Sprintf("($%d, $%d, $%d, $%d),", i*4+1, i*4+2, i*4+3, i*4+4)
		valueTrxDetail = append(valueTrxDetail, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
	}

	// remove last comma and returning id
	queryTrxDetail = queryTrxDetail[0:len(queryTrxDetail)-1] + " RETURNING id"

	// use query to get the returning id
	rows, err := tx.Query(queryTrxDetail, valueTrxDetail...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// loop to scan the returned id
	i := 0
	for rows.Next() {
		if err := rows.Scan(&details[i].ID); err != nil {
			return nil, err
		}
		i++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
		CreatedAt:   transactionDate.Time,
	}
	return res, nil
}
