package repository

import (
	"database/sql"
	"fmt"
	"time"

	"meditrack/structs"
)

type TransactionRepository struct {
	DB *sql.DB
}

func (r *TransactionRepository) GetPrescriptionByID(d int) (any, error) {
	panic("unimplemented")
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

// CreateTransaction Create transaction
func (r *TransactionRepository) CreateTransaction(patientID, prescriptionID, medicineID int, totalPrice float64) (int, error) {
	query := `INSERT INTO transactions (patient_id, prescription_id, medicine_id, total_price, status, created_at) 
			  VALUES ($1, $2, $3, $4, 'pending', $5) RETURNING id`
	var transactionID int
	err := r.DB.QueryRow(query, patientID, prescriptionID, medicineID, totalPrice, time.Now()).Scan(&transactionID)
	if err != nil {
		return 0, err
	}
	return transactionID, nil
}

// UpdateTransactionStatus
func (r *TransactionRepository) UpdateTransactionStatus(transactionID int, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, status, transactionID)
	return err
}

// GetTransactionByID Get detail transaksi by ID
func (r *TransactionRepository) GetTransactionByID(transactionID int) (*structs.Transaction, error) {
	query := `SELECT id, patient_id, prescription_id, medicine_id, total_price, status, created_at FROM transactions WHERE id = $1`
	var transaction structs.Transaction
	err := r.DB.QueryRow(query, transactionID).Scan(
		&transaction.ID, &transaction.PatientID, &transaction.PrescriptionID, &transaction.MedicineID,
		&transaction.TotalPrice, &transaction.Status, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ReduceMedicineStock reduce medicine by transaction
func (repo *TransactionRepository) ReduceMedicineStock(transactionID int) error {
	// Ambil medicine_id dari transaksi
	var medicineID int
	query := `SELECT medicine_id FROM transactions WHERE id = $1`
	err := repo.DB.QueryRow(query, transactionID).Scan(&medicineID)
	if err != nil {
		return err
	}

	// Kurangi stok di tabel medicines (1 unit per transaksi)
	updateQuery := `UPDATE medicines SET stock = stock - 1 WHERE id = $1 AND stock > 0`
	result, err := repo.DB.Exec(updateQuery, medicineID)
	if err != nil {
		return err
	}

	// Cek apakah stok berhasil dikurangi
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("failed to update stock, medicine might be out of stock")
	}

	return nil
}

func (repo *TransactionRepository) CompleteTransaction(transactionID int) error {
	query := `UPDATE transactions SET status = 'completed' WHERE id = $1`
	_, err := repo.DB.Exec(query, transactionID)
	return err
}

func (repo *TransactionRepository) GetMedicineByID(medicineID int) (*structs.Medicine, error) {
	query := `SELECT id, name, category_id, stock, price, expiration_date FROM medicines WHERE id = $1`

	var medicine structs.Medicine
	err := repo.DB.QueryRow(query, medicineID).Scan(
		&medicine.ID, &medicine.Name, &medicine.CategoryID,
		&medicine.Stock, &medicine.Price, &medicine.ExpirationDate,
	)

	if err != nil {
		return nil, err
	}

	return &medicine, nil
}

func (r *TransactionRepository) CompleteTransactionOnce(transactionID int) error {
	var currentStatus string

	// Get status transaction now
	err := r.DB.QueryRow("SELECT status FROM transactions WHERE id = $1", transactionID).Scan(&currentStatus)
	if err != nil {
		return err
	}

	// if transaction completed
	if currentStatus == "completed" {
		return nil
	}

	// If not  completed
	query := "UPDATE transactions SET status = 'completed' WHERE id = $1"
	_, err = r.DB.Exec(query, transactionID)
	return err
}
