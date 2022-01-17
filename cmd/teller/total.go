package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yobrosoft/teller/pkg/dbs"
	"github.com/yobrosoft/teller/pkg/statement"
)

func totalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total STATEMENTS",
		Short: "output the total amount from all transactions in the input statements",
		RunE:  runTotal,
	}
	return cmd
}

func runTotal(cmd *cobra.Command, args []string) error {
	var (
		stmts []*statement.Statement
		files []string
	)

	for _, a := range args {
		info, err := os.Stat(a)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", a, err)
		}

		if !info.IsDir() {
			files = append(files, a)
			continue
		}

		fl, err := ioutil.ReadDir(a)
		if err != nil {
			return err
		}

		for _, f := range fl {
			files = append(files, filepath.Join(a, f.Name()))
		}
	}

	for _, f := range files {
		if filepath.Ext(f) != ".pdf" {
			continue
		}
		s, err := dbs.ParseStatement(f)
		if err != nil {
			return fmt.Errorf("failed to parse statement %s: %w", f, err)
		}
		stmts = append(stmts, s)
	}

	var total, carry float64
	for _, s := range stmts {
		totalSpent := spent(s)
		totalPaid := paid(s)
		total += totalSpent
		carry += totalSpent - totalPaid

		fmt.Printf("%v\n", s.End.Format("Jan 2006"))
		fmt.Printf("%v - %v\t%.2f\n", s.Start.Format("02 Jan"), s.End.Format("02 Jan"), totalSpent)
		fmt.Printf("payed\t%.2f\n", totalPaid)
		fmt.Printf("carry\t%.2f\n", carry)
		fmt.Printf("\n")
	}

	fmt.Printf("--------------------\nTotal: %.2f\n", total)
	return nil
}

func paid(s *statement.Statement) float64 {
	total := 0.0
	for _, t := range s.Transactions {
		if t.Amount < 0 {
			total += t.Amount * -1
		}
	}
	return total
}

func spent(s *statement.Statement) float64 {
	total := 0.0
	for _, t := range s.Transactions {
		if t.Amount > 0 {
			total += t.Amount
		}
	}
	return total
}
