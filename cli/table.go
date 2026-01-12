package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/dev-pi2pie/git-worktree-tasks/ui"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

type tableColumn struct {
	Header   string
	Style    func(value string) lipgloss.Style
	MinWidth int
	Flexible bool
	MaxWidth int
	Truncate bool
}

func renderTable(cmd *cobra.Command, columns []tableColumn, rows [][]string, grid bool) {
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
			cellWidth := lipgloss.Width(row[i])
			if cellWidth > widths[i] {
				widths[i] = cellWidth
			}
		}
	}

	for i, column := range columns {
		if column.MaxWidth > 0 && widths[i] > column.MaxWidth {
			widths[i] = column.MaxWidth
		}
	}

	terminalWidth := terminalWidth()
	widths = fitWidths(terminalWidth, columns, widths, grid)

	if grid {
		fmt.Fprintln(cmd.OutOrStdout(), formatTableDivider(widths))
	}
	fmt.Fprintln(cmd.OutOrStdout(), formatTableRow(columns, widths, headers(columns), true, grid))
	if grid {
		fmt.Fprintln(cmd.OutOrStdout(), formatTableDivider(widths))
	}
	for _, row := range rows {
		fmt.Fprintln(cmd.OutOrStdout(), formatTableRow(columns, widths, row, false, grid))
	}
	if grid {
		fmt.Fprintln(cmd.OutOrStdout(), formatTableDivider(widths))
	}
}

func formatTableRow(columns []tableColumn, widths []int, row []string, isHeader bool, grid bool) string {
	parts := make([]string, len(columns))
	for i, column := range columns {
		cellValue := ""
		if i < len(row) {
			cellValue = row[i]
		}
		cellWidth := widths[i]
		cellValue = clampCell(cellValue, cellWidth, column.Truncate || column.Flexible)
		cell := cellValue
		padding := cellWidth - runewidth.StringWidth(cellValue)
		if padding < 0 {
			padding = 0
		}
		cell = cell + strings.Repeat(" ", padding)
		if isHeader {
			cell = ui.HeaderStyle.Render(cell)
		} else if column.Style != nil {
			cell = column.Style(cellValue).Render(cell)
		}
		if grid {
			parts[i] = " " + cell + " "
		} else {
			parts[i] = cell
		}
	}
	if grid {
		return "|" + strings.Join(parts, "|") + "|"
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

func clampCell(value string, width int, allowTruncate bool) string {
	if width <= 0 {
		return ""
	}
	if !allowTruncate {
		return value
	}
	if runewidth.StringWidth(value) <= width {
		return value
	}
	tail := "..."
	if width <= len(tail) {
		return runewidth.Truncate(value, width, "")
	}
	return runewidth.Truncate(value, width, tail)
}

func fitWidths(terminalWidth int, columns []tableColumn, widths []int, grid bool) []int {
	if terminalWidth <= 0 || len(widths) == 0 {
		return widths
	}
	total := tableWidth(widths, grid)
	if total <= terminalWidth {
		return widths
	}

	minWidths := make([]int, len(columns))
	for i, column := range columns {
		minWidth := column.MinWidth
		if minWidth <= 0 {
			minWidth = 3
		}
		if minWidth > widths[i] {
			minWidth = widths[i]
		}
		minWidths[i] = minWidth
	}

	for total > terminalWidth {
		reduced := false
		for i, column := range columns {
			if !column.Flexible {
				continue
			}
			if widths[i] > minWidths[i] {
				widths[i]--
				total--
				reduced = true
				if total <= terminalWidth {
					break
				}
			}
		}
		if !reduced {
			break
		}
	}

	return widths
}

func tableWidth(widths []int, grid bool) int {
	total := 0
	for _, width := range widths {
		total += width
	}
	if grid {
		if len(widths) > 0 {
			total += (3 * len(widths)) + 1
		}
		return total
	}
	if len(widths) > 1 {
		total += 2 * (len(widths) - 1)
	}
	return total
}

func formatTableDivider(widths []int) string {
	parts := make([]string, len(widths))
	for i, width := range widths {
		parts[i] = strings.Repeat("-", width+2)
	}
	return "+" + strings.Join(parts, "+") + "+"
}

func terminalWidth() int {
	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || width <= 0 {
		return 120
	}
	return width
}
