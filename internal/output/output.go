package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Printer handles formatted output to stdout.
type Printer struct {
	jsonMode bool
}

// NewPrinter creates a new Printer.
func NewPrinter(jsonMode bool) *Printer {
	return &Printer{jsonMode: jsonMode}
}

// IsJSON returns true if the printer is in JSON mode.
func (p *Printer) IsJSON() bool {
	return p.jsonMode
}

// Table prints a table with the given headers and rows.
func (p *Printer) Table(headers []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("  ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}

// JSON marshals v and prints it to stdout.
func (p *Printer) JSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// KV prints aligned key: value pairs.
func (p *Printer) KV(pairs [][2]string) {
	maxLen := 0
	for _, pair := range pairs {
		if len(pair[0]) > maxLen {
			maxLen = len(pair[0])
		}
	}
	for _, pair := range pairs {
		fmt.Printf("%-*s  %s\n", maxLen, pair[0]+":", pair[1])
	}
}

// Line prints a plain line to stdout.
func (p *Printer) Line(msg string) {
	fmt.Println(msg)
}

// Success prints a success message to stdout.
func (p *Printer) Success(msg string) {
	fmt.Println("✓ " + msg)
}
