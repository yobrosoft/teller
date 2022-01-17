package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:   "teller",
		Short: "perform useful things with PDF statements",
		Run:   func(cmd *cobra.Command, _ []string) { cmd.Usage() },
	}

	flgAddr string
)

func main() {
	root.AddCommand(totalCmd())

	root.Flags().StringVarP(&flgAddr, "address", "a", ":9091", "address of the server")
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
