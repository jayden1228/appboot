package model

import "github.com/appboot/appboot/internal/pkg/database"

// Generator data struct
type Generator struct {
	User          string
	Pwd           string
	Host          string
	Port          string
	DB            string
	Path          string
	TemplatePath string
}

func Run(app Generator) error {
	database.SetDbName(app.DB)
	database.SetUp(app.User, app.Pwd, app.Host, app.Port)
	defer database.Close()
	return CreateDBEntity(app.Path, app.TemplatePath)
}
