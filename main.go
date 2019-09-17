package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

var (
	schemaFile string
)

func init() {
	cmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "The schema file which will be used for validation.")
}

var cmd = &cobra.Command{
	Use:   "yamlvalidate -s <schema> [files...]",
	Short: "",
	Long:  "",

	Args: func(cmd *cobra.Command, args []string) error {
		for _, a := range args {
			_, err := os.Stat(a)
			if err != nil {
				return fmt.Errorf("unable to open file %s - %s", a, err)
			}
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		_, err := os.Stat(schemaFile)
		if err != nil {
			log.Fatal("error opening schemaFile file: ", err)
		}

		schema := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaFile))

		ok := true
		var errs []string
		var failedFiles []string
		for _, f := range args {
			fb, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatalf("error reading file to validate: %s - %s", f, err)
			}

			jb, err := yaml.YAMLToJSON(fb)
			if err != nil {
				log.Fatalf("error converting yaml to json: %s - %s", f, err)
			}

			doc := gojsonschema.NewStringLoader(string(jb))
			res, err := gojsonschema.Validate(schema, doc)
			if err != nil {
				log.Fatalf("error validating file: %s - %s", f, err)
			}

			if !res.Valid() {
				ok = false
				failedFiles = append(failedFiles, f)
			}

			for _, e := range res.Errors() {
				e := fmt.Sprintf("error validating file %s - %s %s", f, e.Field(), e.Description())
				log.Println(e)
				errs = append(errs, e)
			}
		}

		if !ok {
			log.Println("Validation failed")
			log.Println("Please fix the following files: ")
			for _, f := range failedFiles {
				log.Println(f)
			}

			os.Exit(1)
		}

		log.Println("All files validated successfully!")
	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
