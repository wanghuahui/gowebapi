package controller

import (
	"time"

	"gowebapi/api/model"
	"gowebapi/api/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// UserByEmail gets user information from email
func UserByEmail(email string) (model.User, error) {
	var err error

	result := model.User{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, password, status_id, first_name FROM user WHERE email = ? LIMIT 1", email)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
			err = c.Find(bson.M{"email": email}).One(&result)
		} else {
			err = model.ErrUnavailable
		}
	case database.TypeBolt:
		err = database.View("user", email, &result)
		if err != nil {
			err = model.ErrNoResult
		}
	default:
		err = model.ErrCode
	}

	return result, model.StandardizeError(err)
}

// UserCreate creates user
func UserCreate(firstName, lastName, email, password string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", firstName,
			lastName, email, password)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

			user := &model.User{
				ObjectID:  bson.NewObjectId(),
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Password:  password,
				StatusID:  1,
				CreatedAt: now,
				UpdatedAt: now,
				Deleted:   0,
			}
			err = c.Insert(user)
		} else {
			err = model.ErrUnavailable
		}
	case database.TypeBolt:
		user := &model.User{
			ObjectID:  bson.NewObjectId(),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
			StatusID:  1,
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}

		err = database.Update("user", user.Email, &user)
	default:
		err = model.ErrCode
	}

	return model.StandardizeError(err)
}

// func accessible(c echo.Context) error {
// 	return c.String(http.StatusOK, "Accessible")
// }
