package model

import (
	"encoding/json"
	"io"

	"github.com/pborman/uuid"
	"github.com/sichacvah/portable_chat/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	ROLE_ADMIN   = "admin"
	USER_ONLINE  = "online"
	USER_OFFLINE = "offline"
)

type User struct {
	Id                   string `json:"id"`
	Password             string `json:password,omitempty`
	PasswordConfirmation string `json:password_confirmation,omitempty`
	Token                string `json:token,omitempty`
	Login                string `json:"login"`
	Name                 string `json:"name"`
	Surname              string `json:"surname"`
	Patronymic           string `json:"patronymic"`
	PersonelNumber       string `json:"personel_number"`
}

// ToJson convert a User to a json string
func (u *User) ToJson() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (u *User) Sanitize() {
	u.Password = ""
	u.PasswordConfirmation = ""
	u.Token = ""
}

func (u *User) setHashedPassword() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		panic(err)
	}
	u.Password = string(hashedPassword)
}

func (u *User) PreSave() {
	u.Id = uuid.New()
	u.setHashedPassword()
	u.PasswordConfirmation = ""
}

func (u *User) SetToken() error {
	authBackend := utils.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(u.Id)
	if err != nil {
		panic(err)
	}
	u.Token = string(token)
	return nil
}

// UserFromJson will decode the imput and return a user
func UserFromJson(data io.Reader) *User {
	decoder := json.NewDecoder(data)

	var user User
	err := decoder.Decode(&user)
	if err == nil {
		return &user
	} else {
		return nil
	}
}

func UserMapToJson(u map[string]*User) string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func UserMapFromJson(data io.Reader) map[string]*User {
	decoder := json.NewDecoder(data)
	var users map[string]*User
	err := decoder.Decode(&users)
	if err == nil {
		return users
	} else {
		return nil
	}
}
