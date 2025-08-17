package repository

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"iris-api/internal/models"
)

type PatientRepositoryInterface interface {
	CreatePatient(req models.CreatePatientRequest) (*models.Patient, error)
	GetPatient(id string) (*models.Patient, error)
	UpdatePatient(id string, req models.UpdatePatientRequest) (*models.Patient, error)
	DeletePatient(id string) error
	ListPatients(limit int, lastKey string) ([]models.Patient, string, error)
}

type PatientRepository struct {
	DynamoDB  *dynamodb.DynamoDB
	TableName string
}

func NewPatientRepository(db *dynamodb.DynamoDB, tableName string) *PatientRepository {
	return &PatientRepository{
		DynamoDB:  db,
		TableName: tableName,
	}
}

func (r *PatientRepository) CreatePatient(req models.CreatePatientRequest) (*models.Patient, error) {
	patient := &models.Patient{
		ID:             uuid.New().String(),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DateOfBirth:    req.DateOfBirth,
		Address:        req.Address,
		InsuranceInfo:  req.InsuranceInfo,
		MedicalHistory: req.MedicalHistory,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	item, err := dynamodbattribute.MarshalMap(patient)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	}

	_, err = r.DynamoDB.PutItem(input)
	if err != nil {
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepository) GetPatient(id string) (*models.Patient, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := r.DynamoDB.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var patient models.Patient
	err = dynamodbattribute.UnmarshalMap(result.Item, &patient)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *PatientRepository) UpdatePatient(id string, req models.UpdatePatientRequest) (*models.Patient, error) {
	updateExpression := "SET updated_at = :updated_at"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":updated_at": {
			S: aws.String(time.Now().Format(time.RFC3339)),
		},
	}

	if req.FirstName != nil {
		updateExpression += ", first_name = :first_name"
		expressionAttributeValues[":first_name"] = &dynamodb.AttributeValue{S: req.FirstName}
	}
	if req.LastName != nil {
		updateExpression += ", last_name = :last_name"
		expressionAttributeValues[":last_name"] = &dynamodb.AttributeValue{S: req.LastName}
	}
	if req.Email != nil {
		updateExpression += ", email = :email"
		expressionAttributeValues[":email"] = &dynamodb.AttributeValue{S: req.Email}
	}
	if req.Phone != nil {
		updateExpression += ", phone = :phone"
		expressionAttributeValues[":phone"] = &dynamodb.AttributeValue{S: req.Phone}
	}
	if req.DateOfBirth != nil {
		updateExpression += ", date_of_birth = :date_of_birth"
		expressionAttributeValues[":date_of_birth"] = &dynamodb.AttributeValue{S: req.DateOfBirth}
	}
	if req.Address != nil {
		updateExpression += ", address = :address"
		expressionAttributeValues[":address"] = &dynamodb.AttributeValue{S: req.Address}
	}
	if req.InsuranceInfo != nil {
		updateExpression += ", insurance_info = :insurance_info"
		expressionAttributeValues[":insurance_info"] = &dynamodb.AttributeValue{S: req.InsuranceInfo}
	}
	if req.MedicalHistory != nil {
		updateExpression += ", medical_history = :medical_history"
		expressionAttributeValues[":medical_history"] = &dynamodb.AttributeValue{S: req.MedicalHistory}
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("ALL_NEW"),
	}

	result, err := r.DynamoDB.UpdateItem(input)
	if err != nil {
		return nil, err
	}

	var patient models.Patient
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &patient)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *PatientRepository) DeletePatient(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := r.DynamoDB.DeleteItem(input)
	return err
}

func (r *PatientRepository) ListPatients(limit int, lastKey string) ([]models.Patient, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.TableName),
	}

	if limit > 0 {
		input.Limit = aws.Int64(int64(limit))
	}

	if lastKey != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(lastKey),
			},
		}
	}

	result, err := r.DynamoDB.Scan(input)
	if err != nil {
		return nil, "", err
	}

	var patients []models.Patient
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &patients)
	if err != nil {
		return nil, "", err
	}

	var nextKey string
	if result.LastEvaluatedKey != nil {
		if key, exists := result.LastEvaluatedKey["id"]; exists && key.S != nil {
			nextKey = *key.S
		}
	}

	return patients, nextKey, nil
}
