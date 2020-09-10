package cmd

import (
	"errors"
	"github.com/appboot/appboot/internal/pkg/database"
	"github.com/appboot/appboot/internal/pkg/path"
	"os"

	"github.com/appboot/appboot/internal/app/appboot/generator"
	"github.com/appboot/appboot/internal/pkg/logger"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// this command generate struct and basic database func
// workflow as following:
// - [appboot create] project template
// - design mysql data struct and create db and table
// - [appboot model] generate the struct and db func

const (
	defaultMysqlHost    = "127.0.0.1"
	defaultMysqlPort    = "3306"
	defaultTemplatePath = "~/.appboot/templates/.Model"
	)

var generate = &cobra.Command{
	Use:   "model",
	Short: "generate struct and database func",
	Long:  `generate struct and database func`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		app := generator.Generator{}
		// Mysql Host
		defaultHost := getDefaultMysqlHost()
		host, err := promptDefault("mysql_host", "mysql host cannot be empty", defaultHost)
		if err != nil {
			logger.LogE(err)
			return
		}
		app.Host = host
		// Port
		defaultPort := getDefaultMysqlPort()
		port, err := promptDefault("mysql_port", "port cannot be empty", defaultPort)
		if err != nil {
			logger.LogE(err)
			return
		}
		app.Port = port

		// db
		db, err := prompt("mysql_db_name", "db cannot be empty")
		if err != nil {
			logger.LogE(err)
			return
		}
		// 设置数据库名称
		app.DB = db
		database.SetDbName(db)

		// user
		user, err := prompt("mysql_user", "user cannot be empty")
		if err != nil {
			logger.LogE(err)
			return
		}
		app.User = user

		// password
		password, err := prompt("mysql_password", "db cannot be empty")
		if err != nil {
			logger.LogE(err)
			return
		}
		app.Pwd = password

		// 设置数据库配置
		database.SetUp(app.User, app.Pwd, app.Host, app.Port)
		defer database.Close()

		tables, err := generator.ListTableNameAndComment()
		if err != nil {
			logger.LogE(err)
			return
		}

		if len(tables) == 0 {
			logger.LogE(errors.New("database is not exist"))
			return
		}

		const All = "All"
		tables = append(tables, All)
		selectedTable, err := promptTables(tables)
		if err != nil {
			logger.LogE(err)
			return
		}
		app.SelectedTable = selectedTable

		// outPath
		outPath, err := prompt("output_path", "path cannot be empty")
		if err != nil {
			logger.LogE(err)
			return
		}

		app.Path = path.HandlerHomeDirAndWorkDir(outPath)

		app.TemplatePath = defaultTemplatePath

		if err := generator.Run(app); err != nil {
			return
		}
	},
}

func promptDefault(label string, alert string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: emptyValidate(alert),
		Default:  defaultValue,
	}
	return prompt.Run()
}


func promptTables(tables []string) (string, error) {
	prompt := promptui.Select{
		Label: "select table",
		Items: tables,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return result, err
	}
	return result, nil
}

func getDefaultMysqlHost() string {
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		return host
	}
	return defaultMysqlHost
}

func getDefaultMysqlPort() string {
	if port := os.Getenv("MYSQL_PORT"); port != "" {
		return port
	}
	return defaultMysqlPort
}

func init() {
	rootCmd.AddCommand(generate)
}
