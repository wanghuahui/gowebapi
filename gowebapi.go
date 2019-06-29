package main

import (
	"encoding/json"
	"gowebapi/api/route"
	"gowebapi/api/shared/database"
	"gowebapi/api/shared/jsonconfig"
	"gowebapi/api/shared/recaptcha"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// configuration contains the application settings
type configuration struct {
	Database  database.Info  `json:"Database"`
	Recaptcha recaptcha.Info `json:"Recaptcha"`
}

// config the settings variable
var config = &configuration{}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Connect to database
	database.Connect(config.Database)
	databases := database.ReadConfig()

	// Echo instance
	e := echo.New()

	route.Load(e)
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		dbname := databases.Type
		return c.String(http.StatusOK, "Hello, "+string(dbname)+"!\n")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
