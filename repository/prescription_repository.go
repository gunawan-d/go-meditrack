package repository

import (
	"fmt"

	"meditrack/structs"

	"github.com/jmoiron/sqlx"
)

type PrescriptionRepository struct {
	DB *sqlx.DB
}

func NewPrescriptionRepository(db *sqlx.DB) *PrescriptionRepository {
	return &PrescriptionRepository{DB: db}
}

// Create prescription / reset
func (r *PrescriptionRepository) CreatePrescription(prescription structs.Prescription) error {
	query := `INSERT INTO prescriptions (doctor_id, patient_id, medicine_id, dosage, status, instructions, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW())`
	_, err := r.DB.Exec(query, prescription.DoctorID, prescription.PatientID, prescription.Medicine, prescription.Dosage, prescription.Status, prescription.Instructions)
	return err
}

// Get prescription by patient id
func (r *PrescriptionRepository) GetPrescriptionByPatientID(patientID int) ([]structs.Prescription, error) {
	var prescriptions []structs.Prescription
	query := `SELECT * FROM prescriptions WHERE patient_id = $1`
	err := r.DB.Select(&prescriptions, query, patientID)
	if err != nil {
		return nil, err
	}
	return prescriptions, nil
}

// Update prescription status
func (r *PrescriptionRepository) UpdatePrescriptionStatus(prescriptionID int, status string) error {
	var exists bool
	err := r.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM prescriptions WHERE id = $1)", prescriptionID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("prescription ID %d not found", prescriptionID)
	}

	query := `UPDATE prescriptions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err = r.DB.Exec(query, status, prescriptionID)
	return err
}
