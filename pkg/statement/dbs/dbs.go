package dbs

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"

	"github.com/yobrosoft/teller/pkg/statement"
)

var (
	transDateReg  = regexp.MustCompile(`^\d{2}\s[A-Z]{3}$`)
	dateReg       = regexp.MustCompile(`^\d{2}\s[A-Z][a-z]{2}\s\d{4}$`)
	amountReg     = regexp.MustCompile(`^\d+\.\d+$`)
	startingToken = "NEW TRANSACTIONS"
)

// Parser is a DBS bank statement parser.
type Parser struct{}

// Parse implements statement.Parser
func (p *Parser) Parse(path string) (*statement.Statement, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var (
		trans []*statement.Transaction
		text  []string
		date  time.Time
	)

	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		rows, err := p.GetTextByRow()
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			for _, word := range row.Content {
				// word is actually a text blob in the pdf
				s := strings.TrimSpace(word.S)
				if len(s) > 0 {
					text = append(text, s)

					if date.IsZero() && dateReg.MatchString(s) {
						date, err = time.Parse("02 Jan 2006", s)
						if err != nil {
							return nil, fmt.Errorf("error parsing statement date (%s): %w", s, err)
						}
					}
				}
			}
		}
	}

	if date.IsZero() {
		return nil, fmt.Errorf("no statement date found")
	}

	i := 0
	for {
		if i >= len(text) {
			break
		}

		t := text[i]
		if transDateReg.MatchString(t) {
			// parse time
			tm, err := time.Parse("02 Jan 2006", fmt.Sprintf("%s %d", capitalCase(t), date.Year()))
			if err != nil {
				return nil, fmt.Errorf("error parsing date (%s): %w", t, err)
			}

			i++

			var (
				amount float64
				desc   string
			)

			// grab all desc blobs
			for i < len(text) {
				candidate := strings.Replace(text[i], ",", "", 1)
				if amountReg.MatchString(candidate) {
					amount, err = strconv.ParseFloat(candidate, 64)
					if err != nil {
						return nil, fmt.Errorf("error parsing transaction ammount of %s on %v", amount, tm)
					}
					break
				}
				desc += text[i]
				i++
			}

			// this indicates a payment to the card
			if text[i+1] == "CR" {
				amount = -amount
				i++
			}

			trans = append(trans, &statement.Transaction{
				Time:        tm,
				Description: desc,
				Amount:      amount,
			})
			continue
		}
		i++
	}

	sort.SliceStable(trans, func(i, j int) bool {
		return trans[i].Time.Before(trans[j].Time)
	})

	// if the transactions span 2 years the beginning will be later
	// than the last ones so send them back a year
	end := len(trans)
	for i := len(trans) - 1; i >= 0; i-- {
		if trans[i].Time.Before(trans[0].Time.AddDate(0, 1, 1)) {
			break
		}
		trans[i].Time.AddDate(-1, 0, 0)
		end = i
	}

	if end < len(trans) {
		trans = append(trans[end:], trans[0:end]...)
	}

	s := &statement.Statement{
		Bank:         "DBS Bank",
		Transactions: trans,
		Start:        trans[0].Time,
		End:          trans[len(trans)-1].Time,
	}

	return s, nil
}

func capitalCase(s string) string {
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

func index(src []string, target string) int {
	for i, s := range src {
		if s == target {
			return i
		}
	}
	return -1
}
