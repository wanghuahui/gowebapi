package controller

import (
	"gowebapi/api/model"
	"gowebapi/api/shared/database"
	"gowebapi/api/shared/passhash"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	user struct {
		ID        uint32 `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		StatusID  uint8  `json:"status_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Token     token  `json:"token"`
	}
	token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   uint32 `json:"expires_in"`
	}
)

// userByEmail gets user information from email
func userByEmail(email string) (model.User, error) {
	var err error

	result := model.User{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, first_name, last_name, email, status_id, created_at FROM user WHERE email = ? LIMIT 1", email)
	default:
		err = model.ErrCode
	}
	// fmt.Println(err.Error())
	return result, model.StandardizeError(err)
}

// userByID gets user information from id
func userByID(id uint32) (model.User, error) {
	var err error

	result := model.User{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, first_name, last_name, email, status_id, created_at FROM user WHERE id = ? LIMIT 1", id)
	default:
		err = model.ErrCode
	}
	return result, model.StandardizeError(err)
}

func tokenCreate(id uint32) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	// claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
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
	_, err = userByEmail(u.Email)

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
	result, _ := userByEmail(u.Email)
	t, err := tokenCreate(result.ID)
	if err != nil {
		return err
	}
	user := &user{
		ID:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
		StatusID:  result.StatusID,
		CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
		Token: token{
			AccessToken: t,
			TokenType:   "Bearer",
			ExpiresIn:   uint32(3600 * 72),
		},
	}
	var res string
	res = fmt.Sprintf(`{"id":%d,"first_name":"%v","last_name":"%v","email":"%v","status_id":%d,"created_at":"%v",
		"token":{"access_token":%v,"token_type":%v,"expires_in":%d}}`,
		result.ID, result.FirstName, result.LastName, result.Email, result.StatusID, result.CreatedAt.Format("2006-01-02 15:04:05"),
		"access_token", "Bearer", 3600)
	return c.JSONBlob(http.StatusCreated, []byte(res))
}

// UserShow shows user
func UserShow(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(uint32)
	// Get database result
	result, err := userByID(id)
	if err != nil {
		return err
	}
}
