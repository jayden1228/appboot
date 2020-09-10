package generator

// Generator data struct
type Generator struct {
	User          string
	Pwd           string
	Host          string
	Port          string
	DB            string
	Path          string
	TemplatePath string
	SelectedTable  string
}

func Run(app Generator) error {
	return CreateDBEntity(app.Path, app.SelectedTable, app.TemplatePath)
}
