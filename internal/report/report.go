package report

import (
	"encoding/json"
	"github.com/azinamouto/kuro/internal/scanner"
	"os"
)

func Write(path string, r scanner.Result) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0600)
}

