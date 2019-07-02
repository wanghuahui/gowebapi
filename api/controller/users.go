package controller

import (
	"fmt"
	"gowebapi/api/model"
	"gowebapi/api/shared/database"
	"gowebapi/api/shared/passhash"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type (
	user struct {
		ID        uint32
		FirstName string
		LastName  string
		Email     string
		StatusID  uint8
		CreatedAt string
		UpdatedAt time.Time
		Token     token
	}
	token struct {
		AccessToken string
		TokenType   string
		ExpiresTime uint32
	}
)

// UserByEmail gets user information from email
func UserByEmail(email string) (model.User, error) {
	var err error

	result := model.User{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, first_name, last_name, email, password, status_id, created_at FROM user WHERE email = ? LIMIT 1", email)
	default:
		err = model.ErrCode
	}
	// fmt.Println(err.Error())
	return result, model.StandardizeError(err)
}

// UserCreate creates user
func UserCreate(c echo.Context) error {
	var err error

	u := &model.User{}
	if err = c.Bind(u); err != nil {
		return err
	}

	password, _ := passhash.HashString(u.Password)
	// Get database result
	_, err = UserByEmail(u.Email)

	if err == model.ErrNoResult { // If success (no user exists with that email)
		switch database.ReadConfig().Type {
		case database.TypeMySQL:
			_, err = database.SQL.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", u.FirstName,
				u.LastName, u.Email, password)
		default:
			err = model.ErrCode
		}
	}

	// Get database result
	result, _ := UserByEmail(u.Email)
	fmt.Println(result.FirstName, result.LastName)
	user := &user{
		ID:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
		StatusID:  result.StatusID,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		Token: token{
			AccessToken: "access_token",
			TokenType:   "Bearer",
			ExpiresTime: 3600,
		},
	}
	return c.JSON(http.StatusCreated, user)
}
