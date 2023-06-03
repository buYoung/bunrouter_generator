package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed template
var templateFiles embed.FS

var (
	templateFilesNames = []string{"controller.go.tmpl", "dto.go.tmpl", "service.go.tmpl"}
)

type Config struct {
	Names         []string
	Dir           string
	templateFiles *template.Template
}
type templateData struct {
	Name        string
	LowerName   string
	PackageName string
}

func main() {
	var cmd = &cobra.Command{
		Use:   "bunrouter_generator",
		Short: "This application generates a CRUD for golang's bunrouter.",
		Run: func(cmd *cobra.Command, args []string) {
			config := &Config{
				Names: viper.GetStringSlice("names"),
				Dir:   viper.GetString("dir"),
			}
			generateCRUD(config)
		},
	}

	cmd.PersistentFlags().StringSliceP("names", "n", []string{}, "This option takes a list of folders, package names for the CRUD to be generated.")
	err := cmd.MarkFlagRequired("names")
	if err != nil {
		customErrorMessages(err)
		return
	}
	viper.BindPFlag("names", cmd.PersistentFlags().Lookup("names"))

	cmd.PersistentFlags().String("dir", "internals/api", "This option sets the default location to store the CRUD files that will be generated.\n")
	viper.BindPFlag("dir", cmd.PersistentFlags().Lookup("dir"))

	cmd.Execute()
}

func generateCRUD(config *Config) {
	config.templateFiles = readTemplates()

	withDefault(config)
}

func withDefault(config *Config) {
	fmt.Println("Generating Default")

	for _, name := range config.Names {

		name = strings.TrimSpace(name)
		lowerName := strings.ToLower(name)
		upperName := stringToUpperCaseWithfirst(name)

		folderCreate(filepath.Join(config.Dir, name))

		for _, templateName := range templateFilesNames {
			fileName := templateName[:len(templateName)-len(".tmpl")]
			generateTemplate := config.templateFiles.Lookup(templateName)
			joinFilePath := filepath.Join(config.Dir, name, fileName)
			stat, err := os.Stat(joinFilePath)
			if err == nil {
				fmt.Printf("File %s already exists, skipping...\n", stat.Name())
				continue
			}

			file, err := os.Create(joinFilePath)
			if err != nil {
				customErrorMessages(err)
			}

			err = generateTemplate.Execute(
				file,
				templateData{Name: upperName, LowerName: lowerName, PackageName: lowerName})
			if err != nil {
				customErrorMessages(err)
			}

			file.Close()
		}
	}
}

func readTemplates() *template.Template {
	templateFiles, err := template.ParseFS(templateFiles, "template/*")
	if err != nil {
		customErrorMessages(err)
	}
	return templateFiles
}

func folderCreate(dir string) {
	stat, err := os.Stat(dir)
	if err == nil {
		if !stat.IsDir() {
			customErrorMessages("The path provided is not a directory")
		}

		return
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		customErrorMessages(err)
	}
}

func stringToUpperCaseWithfirst(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}

func customErrorMessages(err interface{}) {
	fmt.Println(err)

	os.Exit(1)
}
