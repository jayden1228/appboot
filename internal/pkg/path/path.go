package path

import (
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// HandleHomedir replace ~ with homeDir
func HandleHomedir(filePath string) string {
	const homeDir = "~"
	if strings.HasPrefix(filePath, homeDir) {
		home, err := homedir.Dir()
		if err != nil {
			return filePath
		}
		result := strings.Replace(filePath, homeDir, home, 1)
		return result
	}
	return filePath
}

// HandlerWorkDir replace . with workDir
func HandlerWorkDir(filePath string) string {
	const pwdDir = "."
	if strings.HasPrefix(filePath, pwdDir) {
		pwd, err := os.Getwd()
		if err != nil {
			return filePath
		}
		result := strings.Replace(filePath, pwdDir, pwd, 1)
		return result
	}
	return filePath
}

// HandlerHomeDirAndWorkDir relace ~ with homeDir and . with workDir
func HandlerHomeDirAndWorkDir(filepath string) string {
	return HandleHomedir(HandlerWorkDir(filepath))
}