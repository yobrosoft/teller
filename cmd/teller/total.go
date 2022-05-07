package main

import (
	"fmt"

	"github.com/spf13/cobra"
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
	stmts, err := statement.ParseFiles(parser, args...)
	if err != nil {
		return fmt.Errorf("failed to parse statements: %w", err)
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
