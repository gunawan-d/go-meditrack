package repository

import (
	"database/sql"
	"errors"
	"log"
	"meditrack/structs"
	"time"

	"github.com/jmoiron/sqlx"
)

type MedicineRepository struct {
	DB *sqlx.DB
}

func NewMedicineRepository(db *sqlx.DB) *MedicineRepository {
	return &MedicineRepository{DB: db}
}

// Add Medicine
func (r *MedicineRepository) CreateMedicine(medicine structs.Medicine) error {
	var categoryID int

	// Check category_id Available
	err := r.DB.QueryRow("SELECT id FROM categories WHERE id = $1", medicine.CategoryID).Scan(&categoryID)

	if err == sql.ErrNoRows { // Jika tidak ada, buat kategori baru
		log.Println("Category not found, creating new category...")
		categoryName := "Default Category" // Bisa diubah sesuai kebutuhan
		err = r.DB.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id", categoryName).Scan(&categoryID)
		if err != nil {
			log.Println("Error creating category:", err)
			return err
		}
	} else if err != nil {
		log.Println("Error checking category:", err)
		return err
	}

	// Parse expiration_date
	parseDate, err := time.Parse("2006-01-02", medicine.ExpirationDate)
	if err != nil {
		log.Println("Invalid expiration date format:", err)
		return err
	}

	// Insert medicine
	query := `INSERT INTO medicines (name, category_id, price, stock, expiration_date) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err = r.DB.Exec(query, medicine.Name, categoryID, medicine.Price, medicine.Stock, parseDate)
	return err
}

// Get All medicines
func (r *MedicineRepository) GetAllMedicines() ([]structs.Medicine, error) {
	var medicines []structs.Medicine
	query := `SELECT * FROM medicines`
	err := r.DB.Select(&medicines, query)
	return medicines, err
}

// Update stock medicine
func (r *MedicineRepository) UpdateMedicineStock(id int, stock int) error {
	query := `UPDATE medicines SET stock = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, stock, id)
	return err
}

func (repo *MedicineRepository) ReduceMedicineStock(medicineID int) error {
	query := "UPDATE medicines SET stock = stock - 1 WHERE id = ? AND stock > 0"
	result, err := repo.DB.Exec(query, medicineID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("insufficient stock")
	}

	return nil
}

func (r *MedicineRepository) DeleteMedicine(medicineID int) error {
	_, err := r.DB.Exec("DELETE FROM medicines WHERE id = $1", medicineID)
	return err
}
