/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"converter/common"
	"converter/converters"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts file from one format to another",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		format := ""
		cmd.Flags().StringVar(&format, "format", "", "for instance yml2json")
		if format == "" {
			format = "yml2json"
		}
		var converter common.FormatConverter
		switch format {
		case "yml2json":
			converter = converters.NewYmlToJsonConverter()
		default:
			fmt.Printf("format %s is not supported\n", format)
		}

		sourceFileName, _ := cmd.Flags().GetString("source")
		destinationFileName, _ := cmd.Flags().GetString("destination")
		overwriteExistingFile, _ := cmd.Flags().GetBool("force")

		if sourceFileName == destinationFileName {
			fmt.Printf("source and destination are the same %s\n", sourceFileName)
			return
		}

		// prepare
		if st, err := os.Stat(sourceFileName); err != nil {
			fmt.Printf("could not open source file %s: %v\n", sourceFileName, err)
			return
		} else if st.IsDir() {
			fmt.Printf("source %s is dir \n", sourceFileName)
			return
		}
		if st, err := os.Stat(destinationFileName); err == nil {
			if st.IsDir() {
				fmt.Printf("destnation %s is dir \n", destinationFileName)
				return
			}
			if !overwriteExistingFile {
				fmt.Printf("destination file exists, specify --force to overwite it \n")
				return
			}
		}
		var (
			sourceReader      io.ReadCloser
			destinationWriter io.WriteCloser
			err               error
		)
		if sourceReader, err = os.Open(sourceFileName); err != nil {
			fmt.Printf("could not open source file %s: %v\n", sourceFileName, err)
			return
		}
		defer func() { sourceReader.Close() }()
		if destinationWriter, err = os.Create(destinationFileName); err != nil {
			fmt.Printf("could not open destination file %s: %v\n", destinationFileName, err)
			return
		}
		defer func() { destinationWriter.Close() }()
		err = converter.Convert(destinationWriter, sourceReader)
		if err != nil {
			fmt.Printf("conversion error: %v\n", err)
			return
		}
		fmt.Println("successfully converted\n")
	},
}

func init() {
	convertCmd.Flags().StringP("source", "s", "", "--source=/home/www/example.yml")
	convertCmd.MarkFlagRequired("source")
	convertCmd.Flags().StringP("destination", "d", "", "--destination=/home/www/example.json")
	convertCmd.MarkFlagRequired("destination")
	convertCmd.Flags().BoolP("force", "f", false, "--force")

	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
