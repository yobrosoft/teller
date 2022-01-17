package statement

import (
	"time"
)

type Transaction struct {
	Time        time.Time
	Description string
	Amount      float64
}

type Statement struct {
	Bank         string
	Start        time.Time
	End          time.Time
	Transactions []*Transaction
}
