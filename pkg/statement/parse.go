package statement

import (
	"io/fs"
	"path/filepath"
)

// Parser is capable of parsing specific bank statement files.
type Parser interface {
	// Parse parses a given bank statement file and returns a Statement struct.
	Parse(path string) (*Statement, error)
}

// ParseFiles takes a list of files or directories and returns parsed statements from the files.
func ParseFiles(parser Parser, files ...string) ([]*Statement, error) {
	stmts := []*Statement{}
	for _, f := range files {
		err := filepath.Walk(f, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			s, err := parser.Parse(path)
			if err != nil {
				return err
			}
			stmts = append(stmts, s)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return stmts, nil
}
