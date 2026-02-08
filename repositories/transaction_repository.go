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

func (repo *TransactionRepository) Report(start_date string, end_date string) (*models.Report, error) {
	var report models.Report

	dateFilter := ""
	args := []interface{}{}

	// check if start_date and end_date is filled
	if start_date != "" && end_date != "" {
		dateFilter = "WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2"
		args = append(args, start_date, end_date)
	} else {
		dateFilter = "WHERE DATE(t.created_at) = CURRENT_DATE"
	}

	// get the total revenue and total transasksi
	queryTotal := fmt.Sprintf(`SELECT COALESCE(SUM(total_amount), 0) as total_revenue, 
		COUNT(*) as total_transaksi 
		FROM transactions t 
		%s`, dateFilter)

	err := repo.db.QueryRow(queryTotal, args...).Scan(&report.TotalRevenue, &report.TotalTransaction)
	if err != nil {
		return nil, err
	}

	// get the most sold product
	queryTerjual := fmt.Sprintf(`SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty_terjual 
			FROM transaction_details td 
			JOIN transactions t ON td.transaction_id = t.id 
			JOIN products p ON td.product_id = p.id 
			%s 
			GROUP BY p.id, p.name LIMIT 1`, dateFilter)

	err = repo.db.QueryRow(queryTerjual, args...).Scan(&report.ProdukTerlaris.Name, &report.ProdukTerlaris.QtySold)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// set data to 0 when no rows
	if err == sql.ErrNoRows {
		report.ProdukTerlaris.Name = ""
		report.ProdukTerlaris.QtySold = 0
	}

	return &report, nil
}
