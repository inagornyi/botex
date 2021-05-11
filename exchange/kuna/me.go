package kuna

type Me struct {
	Email     string    `json:"email"`
	Activated bool      `json:"activated"`
	Accounts  []Account `json:"accounts"`
}
