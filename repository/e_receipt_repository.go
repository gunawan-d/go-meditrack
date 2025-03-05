package repository

import (
	"meditrack/structs"

	"github.com/jmoiron/sqlx"
)

type EReceiptRepository struct {
	DB *sqlx.DB
}

func NewEReceiptRepository(db *sqlx.DB) *EReceiptRepository {
	return &EReceiptRepository{DB: db}
}

func (r *EReceiptRepository) GetEReceiptByTransactionID(transactionID int) (*structs.EReceipt, error) {
	var eReceipt structs.EReceipt
	query := `
        SELECT 
            t.id AS transaction_id,
            u.name AS patient_name,
            m.name AS medicine_name,
            p.dosage,
            t.total_price,
            t.status,
            t.created_at AS issued_at
        FROM transactions t
        JOIN users u ON t.patient_id = u.id
        JOIN medicines m ON t.medicine_id = m.id
        JOIN prescriptions p ON t.prescription_id = p.id
        WHERE t.id = $1
    `
	err := r.DB.Get(&eReceipt, query, transactionID)
	if err != nil {
		return nil, err
	}
	return &eReceipt, nil
}
