package models

import (
	"time"
)

type Patient struct {
	ID             string    `json:"id" dynamodb:"id"`
	FirstName      string    `json:"first_name" dynamodb:"first_name"`
	LastName       string    `json:"last_name" dynamodb:"last_name"`
	Email          string    `json:"email" dynamodb:"email"`
	Phone          string    `json:"phone" dynamodb:"phone"`
	DateOfBirth    string    `json:"date_of_birth" dynamodb:"date_of_birth"`
	Address        string    `json:"address" dynamodb:"address"`
	InsuranceInfo  string    `json:"insurance_info" dynamodb:"insurance_info"`
	MedicalHistory string    `json:"medical_history" dynamodb:"medical_history"`
	CreatedAt      time.Time `json:"created_at" dynamodb:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" dynamodb:"updated_at"`
}

type CreatePatientRequest struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Phone          string `json:"phone" validate:"required"`
	DateOfBirth    string `json:"date_of_birth" validate:"required"`
	Address        string `json:"address"`
	InsuranceInfo  string `json:"insurance_info"`
	MedicalHistory string `json:"medical_history"`
}

type UpdatePatientRequest struct {
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	DateOfBirth    *string `json:"date_of_birth,omitempty"`
	Address        *string `json:"address,omitempty"`
	InsuranceInfo  *string `json:"insurance_info,omitempty"`
	MedicalHistory *string `json:"medical_history,omitempty"`
}
