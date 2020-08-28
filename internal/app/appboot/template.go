package appboot

import (
	"container/list"
	"fmt"
	"github.com/CatchZeng/gutils/array"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CatchZeng/gutils/file"
	gos "github.com/CatchZeng/gutils/os"
	"github.com/appboot/appboot/configs"
	"github.com/appboot/appboot/internal/pkg/logger"
)

// File template file
type File struct {
	Path    string
	Content string
}

// GetTemplates get templates
func GetTemplates() []string {
	var templates []string

	root, err := configs.GetTemplateRoot()
	if err != nil {
		return templates
	}

	templates, _ = file.GetDirListWithFilter(root, func(info os.FileInfo) bool {
		return !strings.HasPrefix(info.Name(), ".")
	})
	return templates
}

var paths = list.New()

// GetFiles get files from template path
func GetFiles(templatePath string) (*list.List, error) {
	var files = list.New()
	if err := filepath.Walk(templatePath, walkFunc); err != nil {
		return files, err
	}

	for i := paths.Front(); i != nil; i = i.Next() {
		p := i.Value.(string)
		bytes, err := ioutil.ReadFile(p)
		if err != nil {
			return files, err
		}
		content := string(bytes)
		f := File{
			Path:    p,
			Content: content,
		}
		files.PushBack(f)
	}
	return files, nil
}

// UpdateTemplate update template with Git
func UpdateTemplate(name string) error {
	return UpdateTemplateWithDownloader(name, NewDownloader())
}

// UpdateTemplateWithDownloader update template
func UpdateTemplateWithDownloader(name string, downloader Downloader) error {
	// get template path
	root, err := configs.GetTemplateRoot()
	if err != nil {
		return err
	}
	templatePath := path.Join(root, name)

	// download templates
	tempDir, err := downloadTemplates(downloader)
	defer os.RemoveAll(tempDir)
	if err != nil {
		if file.Exists(templatePath) {
			logger.LogW(fmt.Sprintf("update template error: %v \nuse old template.", err))
			return nil
		}
		return err
	}

	// check template
	src := path.Join(tempDir, name)
	if !file.Exists(src) {
		logger.LogW(fmt.Sprintf("can not get template %v from %v", name, configs.GetTemplateSource()))
		return nil
	}

	// update template
	if err := updateTemplate(name, root, src); err != nil {
		return err
	}

	if !file.Exists(templatePath) {
		return fmt.Errorf("can not find template from appboot templates, you can add your custom template to %s", root)
	}

	return nil
}

func updateTemplate(name, root, src string) error {
	templatePath := path.Join(root, name)

	// recreate template directory
	mode := file.Mode(templatePath)
	_ = os.RemoveAll(templatePath)
	if err := os.MkdirAll(templatePath, mode); err != nil {
		return err
	}

	// copy template to root from src
	cp := "cp -rf " + src + " " + root
	if err := gos.RunBashCommand(cp); err != nil {
		return err
	}

	return nil
}

// UpdateAllTemplates update all templates with git
func UpdateAllTemplates() error {
	return UpdateAllTemplatesWithDownloader(NewDownloader())
}

// UpdateAllTemplatesWithDownloader update all templates
func UpdateAllTemplatesWithDownloader(downloader Downloader) error {
	// download templates
	tempDir, err := downloadTemplates(downloader)
	defer os.RemoveAll(tempDir)
	if err != nil {
		return err
	}

	// get template root
	root, err := configs.GetTemplateRoot()
	if err != nil {
		return err
	}

	// remove existed templates
	templates := GetTemplates()
	for _, name := range templates {
		list, _ := file.GetDirListWithFilter(tempDir, func(info os.FileInfo) bool {
			return !strings.HasPrefix(info.Name(), ".")
		})
		if array.ContainString(list, name) {
			existed := path.Join(root, name)
			os.RemoveAll(existed)
		}
	}

	// update templates
	cp := "cp -rf " + tempDir + "/*" + " " + root
	if err := gos.RunBashCommand(cp); err != nil {
		return err
	}

	return nil
}

// Warning: caller should clean tempDir after used
// defer os.RemoveAll(tempDir)
func downloadTemplates(downloader Downloader) (string, error) {
	source := configs.GetTemplateSource()
	tempDir, err := ioutil.TempDir(os.TempDir(), "template")
	if err != nil {
		return tempDir, err
	}
	if err := downloader.Run(source, tempDir); err != nil {
		return tempDir, err
	}
	return tempDir, nil
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		paths.PushBack(path)
	}
	return nil
}
