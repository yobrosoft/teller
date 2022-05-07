package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/yobrosoft/teller/pkg/statement"
	"github.com/yobrosoft/teller/pkg/statement/dbs"
)

var (
	root = &cobra.Command{
		Use:   "teller",
		Short: "perform useful things with PDF statements",
		Run:   func(cmd *cobra.Command, _ []string) { cmd.Usage() },
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return setParser()
		},
	}

	flgAddr   string
	flgParser string

	parser statement.Parser
)

func main() {
	root.AddCommand(totalCmd())

	root.Flags().StringVarP(&flgAddr, "address", "a", ":9091", "address of the server")
	root.PersistentFlags().StringVarP(&flgParser, "parser", "p", "", "parser to use for statements")

	cobra.MarkFlagRequired(root.PersistentFlags(), "parser")

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setParser() error {
	switch flgParser {
	case "dbs":
		parser = &dbs.Parser{}
	default:
		return fmt.Errorf("usupported statement parser '%s'", flgParser)
	}
	return nil
}
