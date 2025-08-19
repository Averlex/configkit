package configkit

import (
	"encoding/json"
	"fmt"
	"io"
)

// PlainVersionPrinter returns a function that prints version as plain text.
func PlainVersionPrinter(version string) func(io.Writer) error {
	return func(w io.Writer) error {
		if w == nil {
			return fmt.Errorf("nil writer received")
		}
		_, err := fmt.Fprintln(w, version)
		return err
	}
}

// JSONVersionPrinter returns a function that prints version in JSON format.
func JSONVersionPrinter(version, commit, date string) func(io.Writer) error {
	type info struct {
		Version string `json:"version"`
		Commit  string `json:"commit,omitempty"`
		Date    string `json:"date,omitempty"`
	}
	data := info{Version: version, Commit: commit, Date: date}
	return func(w io.Writer) error {
		if w == nil {
			return fmt.Errorf("nil writer received")
		}
		return json.NewEncoder(w).Encode(data)
	}
}
