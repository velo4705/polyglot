package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velo4705/polyglot/internal/language"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported languages",
	RunE:  listLanguages,
}

func listLanguages(cmd *cobra.Command, args []string) error {
	fmt.Println("Supported Languages:")
	fmt.Println("-------------------")

	for _, lang := range language.GetAllLanguages() {
		handler := language.GetHandler(lang)
		if handler != nil {
			fmt.Printf("%-15s %s\n", lang.Name(), handler.Extensions())
		}
	}

	return nil
}
