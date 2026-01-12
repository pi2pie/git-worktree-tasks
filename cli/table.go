package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type tableColumn struct {
	Header string
	Style  func(value string) lipgloss.Style
}

func renderTable(cmd *cobra.Command, columns []tableColumn, rows [][]string) {
	if len(columns) == 0 {
		return
	}

	widths := make([]int, len(columns))
	for i, column := range columns {
		widths[i] = lipgloss.Width(column.Header)
	}
	for _, row := range rows {
		for i := range columns {
			if i >= len(row) {
				continue
			}
			if width := lipgloss.Width(row[i]); width > widths[i] {
				widths[i] = width
			}
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), formatTableRow(columns, widths, headers(columns), true))
	for _, row := range rows {
		fmt.Fprintln(cmd.OutOrStdout(), formatTableRow(columns, widths, row, false))
	}
}

func formatTableRow(columns []tableColumn, widths []int, row []string, isHeader bool) string {
	parts := make([]string, len(columns))
	for i, column := range columns {
		cellValue := ""
		if i < len(row) {
			cellValue = row[i]
		}
		cell := cellValue
		padding := widths[i] - lipgloss.Width(cellValue)
		if padding < 0 {
			padding = 0
		}
		cell = cell + strings.Repeat(" ", padding)
		if isHeader {
			cell = ui.HeaderStyle.Render(cell)
		} else if column.Style != nil {
			cell = column.Style(cellValue).Render(cell)
		}
		parts[i] = cell
	}
	return strings.Join(parts, "  ")
}

func headers(columns []tableColumn) []string {
	headerRow := make([]string, len(columns))
	for i, column := range columns {
		headerRow[i] = column.Header
	}
	return headerRow
}
