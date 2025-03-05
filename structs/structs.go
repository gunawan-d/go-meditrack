package structs

import "time"

// User represents a user in the system
type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

// Category represents a medicine category
type Category struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Medicine represents a medicine item
type Medicine struct {
	ID             int       `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	CategoryID     int       `db:"category_id" json:"category_id"`
	Stock          int       `db:"stock" json:"stock"`
	Price          float64   `db:"price" json:"price"`
	ExpirationDate string    `db:"expiration_date" json:"expiration_date"` 
	CreatedAt      time.Time `db:"created_at" json:"-"`
}

// Prescription represents a medical prescription
type Prescription struct {
    ID          int       `db:"id" json:"id"`
    DoctorID    int       `db:"doctor_id" json:"doctor_id"`
    PatientID   int       `db:"patient_id" json:"patient_id"`
    Medicine  int       `db:"medicine_id" json:"medicine_id"`  
    Dosage      string    `db:"dosage" json:"dosage"`
    Instructions string   `db:"instructions" json:"instructions"`
    Status      string    `db:"status" json:"status"` 
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdateAt    time.Time `db:"updated_at" json:"updated_at"`
}


// Transaction represents a transaction for medicine purchase
type Transaction struct {
	ID            int       `db:"id"`
	PatientID     int       `db:"patient_id"`
	PrescriptionID int      `db:"prescription_id"`
	MedicineID    int       `db:"medicine_id"`
	TotalPrice    float64   `db:"total_price"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
}

// EReceipt represents an electronic receipt
type EReceipt struct {
    TransactionID  int       `db:"transaction_id"`
    PatientName    string    `db:"patient_name"`
    MedicineName   string    `db:"medicine_name"`
    Dosage         string    `db:"dosage"`
    TotalPrice     float64   `db:"total_price"`
    Status         string    `db:"status"`
    IssuedAt       time.Time `db:"issued_at"`
}
