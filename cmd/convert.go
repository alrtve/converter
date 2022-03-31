package cmd

import (
	"converter/common"
	"converter/converters"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts file from one format to another",
	Long:  `Converts file from one format to another`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		if format == "" {
			format = "yml2json"
		}
		var (
			converter      common.FormatConverter
			sourceExt      string
			destinationExt string
		)
		switch format {
		case "yml2json":
			prettyPrint, _ := cmd.Flags().GetBool("prettyprint")
			converter = converters.NewYmlToJsonConverter().WithPrettyPrint(prettyPrint)
			sourceExt = "yml"
			destinationExt = "json"
		default:
			fmt.Printf("Error: format %s is not supported\n", format)
		}

		sourceName, _ := cmd.Flags().GetString("source")
		destinationName, _ := cmd.Flags().GetString("destination")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		// prepare
		var (
			sourceFi       os.FileInfo
			destinationFi  os.FileInfo
			sourceErr      error
			destinationErr error
		)
		if sourceFi, sourceErr = os.Stat(sourceName); sourceErr != nil {
			fmt.Printf("Error: could not open file %s: %v\n", sourceName, sourceErr)
			return
		}
		if destinationFi, destinationErr = os.Stat(destinationName); destinationErr != nil && !errors.Is(destinationErr, os.ErrNotExist) {
			fmt.Printf("Error: could not processin destination %s: %v\n", destinationName, destinationErr)
		}

		// both are files
		if !sourceFi.IsDir() && (errors.Is(destinationErr, os.ErrNotExist) || !destinationFi.IsDir()) {
			if sourceName == destinationName {
				fmt.Printf("Error: source and destination are the same %s\n", sourceName)
				return
			}

			if destinationErr == nil && !overwrite {
				fmt.Printf("Error: destination file exists, specify --overwrite paramter to overwrite it\n")
				return
			}
			if err := convertFile(sourceName, destinationName, converter, overwrite); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		} else if sourceFi.IsDir() && (errors.Is(destinationErr, os.ErrNotExist) || destinationFi.IsDir()) {
			// both are dirs
			overwrite = true
			if err := convertDir(sourceName, destinationName, sourceExt, destinationExt, converter, overwrite); err != nil {
				fmt.Printf("Error: could not convert %s -> %s: %v\n", sourceName, destinationName, err)
				return
			}
		} else {
			fmt.Printf("Error: both source and destination must be either files or dirs\n")
			return
		}

		fmt.Printf("Successfully converted\n")
	},
}

func convertFile(sourceName, destinationName string, converter common.FormatConverter, overwrite bool) (err error) {
	if _, err := os.Stat(destinationName); err == nil && !overwrite {
		return nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	var (
		sourceReader      io.ReadCloser
		destinationWriter io.WriteCloser
	)
	if sourceReader, err = os.Open(sourceName); err != nil {
		return fmt.Errorf("could not open source file %s: %w", sourceName, err)
	}
	defer func() { sourceReader.Close() }()
	if destinationWriter, err = os.Create(destinationName); err != nil {
		return fmt.Errorf("could not open destination file %s: %w", destinationName, err)
	}
	defer func() { destinationWriter.Close() }()
	if err := converter.Convert(destinationWriter, sourceReader); err != nil {
		return fmt.Errorf("conversion error %s -> %s: %w", sourceName, destinationName, err)
	}
	return nil
}

func convertDir(sourceDir, destinationDir, sourceExt, destinationExt string, converter common.FormatConverter, overwrite bool) error {
	sourceFis, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(destinationDir); err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(destinationDir, os.ModePerm); err != nil {
			return fmt.Errorf("could not create destination dir %s: %v\n", destinationDir, err)
		}
	} else if err != nil {
		return err
	}

	for _, stat := range sourceFis {
		if stat.IsDir() {
			if err := convertDir(path.Join(sourceDir, stat.Name()), path.Join(destinationDir, stat.Name()), sourceExt, destinationExt, converter, overwrite); err != nil {
				return err
			}
		} else {
			parts := strings.Split(stat.Name(), ".")
			ext := parts[len(parts)-1]
			if ext != sourceExt {
				continue
			}
			parts[len(parts)-1] = destinationExt
			if err := convertFile(path.Join(sourceDir, stat.Name()), path.Join(destinationDir, strings.Join(parts, ".")), converter, overwrite); err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	convertCmd.Flags().StringP("format", "f", "yml2json", "Conversion direction in from2to format (i.e. yml2json)")
	convertCmd.Flags().StringP("source", "s", "", "Source file or dir")
	convertCmd.Flags().StringP("destination", "d", "", "Destination file or dir")
	convertCmd.Flags().BoolP("prettyprint", "", false, "Prettyprint output")
	convertCmd.Flags().BoolP("overwrite", "o", false, "Overwrite destination file if exists one. For directory source and destination the flag is always set")

	convertCmd.MarkFlagRequired("source")
	convertCmd.MarkFlagRequired("destination")

	rootCmd.AddCommand(convertCmd)
}
