package models

// User结构体
type Account struct {
	ID          int
	Username    string `json:"username"`
	Password    string `json:"password"`
	Repwd       string `json:"repwd"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	PubDatetime string
	UpdDatetime string
}
