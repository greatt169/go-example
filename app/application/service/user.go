package service

import (
	"fmt"
	helpers "github.com/AeroAgency/golang-helpers-lib"
)

type User struct {
	errorFormatter *helpers.ErrorFormatter
	jwt            *helpers.Jwt
	meta           *helpers.Meta
}

// Конструктор
func NewUser() *User {
	errorFormatter := &helpers.ErrorFormatter{}
	return &User{
		errorFormatter: errorFormatter,
		jwt:            &helpers.Jwt{},
		meta:           &helpers.Meta{},
	}
}

// Возвращает имя пользователея для отображенич
// Формат: Иванов И.И.
func (u User) GetNameByToken(tokenString string) string {
	claims, _ := u.jwt.ParseUnverified(tokenString)
	givenName := fmt.Sprintf("%s", claims["given_name"])
	middleName := fmt.Sprintf("%s", claims["middle_name"])
	familyName := fmt.Sprintf("%s", claims["family_name"])
	givenNameInitial := string([]rune(givenName)[0])
	middleNameInitial := string([]rune(middleName)[0])
	name := familyName + " " + givenNameInitial + "." + middleNameInitial
	return name
}
