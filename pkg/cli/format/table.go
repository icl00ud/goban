// Package format provides output formatters for CLI output.
package format

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Table prints tabular data to stdout.
func Table(headers []string, rows [][]string) {
	tw := tablewriter.NewWriter(os.Stdout)
	tw.SetHeader(headers)
	tw.SetAutoWrapText(false)
	tw.SetAutoFormatHeaders(true)
	tw.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tw.SetAlignment(tablewriter.ALIGN_LEFT)
	tw.SetBorder(true)
	tw.AppendBulk(rows)
	tw.Render()
}

// Msg prints a simple formatted message.
func Msg(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
