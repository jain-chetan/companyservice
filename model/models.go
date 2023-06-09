package models

import "github.com/google/uuid"

type Company struct {
	ID          uuid.UUID   `json:"id,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Employees   int         `json:"employees"`
	Registered  bool        `json:"registered"`
	Type        CompanyType `json:"type"`
}

type CompanyType string

const (
	Corporation        CompanyType = "Corporation"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

//DBConfig has information required to connect to DB
type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

//Response has the message and code
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//CreateResponse - response structure for create where ID is sent as response
type CreateResponse struct {
	ID      uuid.UUID `json:"id"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
