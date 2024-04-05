package program

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type TableCellAlignment string

const (
	TableCellAlignmentLeft  TableCellAlignment = "left"
	TableCellAlignmentRight TableCellAlignment = "right"
)

type Table struct {
	Columns []TableColumn
	Rows    [][]interface{}
}

type TableColumn struct {
	Label     string
	Alignment TableCellAlignment
}

func NewTable() *Table {
	return &Table{
		Rows: make([][]interface{}, 0),
	}
}

func (t *Table) AddColumn(c TableColumn) {
	t.Columns = append(t.Columns, c)
}

func (t *Table) AddRow(row ...interface{}) {
	t.Rows = append(t.Rows, row)
}

func (t *Table) Print() {
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	rows := t.Render()
	widths := t.columnWidths(rows)

	if isTerminal {
		for i, c := range t.Columns {
			if i > 0 {
				fmt.Fprintf(os.Stderr, "  ")
			}

			fmtString := "%-*s"
			if c.Alignment == TableCellAlignmentRight {
				fmtString = "%*s"
			}

			label := fmt.Sprintf(fmtString, widths[i], strings.ToUpper(c.Label))
			fmt.Fprintf(os.Stderr, label)
		}

		fmt.Fprintln(os.Stderr)
	}

	for _, row := range rows {
		for j, s := range row {
			c := t.Columns[j]

			if j > 0 {
				fmt.Printf("  ")
			}

			fmtString := "%-*s"
			if c.Alignment == TableCellAlignmentRight {
				fmtString = "%*s"
			}

			fmt.Printf(fmtString, widths[j], s)
		}

		fmt.Println("")
	}
}

func (t *Table) Render() [][]string {
	rows := make([][]string, len(t.Rows))

	for i, row := range t.Rows {
		rows[i] = make([]string, len(row))

		for j, value := range row {
			rows[i][j] = t.RenderValue(value)
		}
	}

	return rows
}

func (t *Table) RenderValue(value interface{}) (s string) {
	switch v := value.(type) {
	case time.Time:
		s = v.Format(time.RFC3339)
	case *time.Time:
		if v != nil {
			s = t.RenderValue(*v)
		}
	default:
		s = fmt.Sprintf("%v", value)
	}

	return
}

func (t *Table) columnWidths(rows [][]string) []int {
	widths := make([]int, len(t.Columns))

	for i, c := range t.Columns {
		widths[i] = len(c.Label)
	}

	for _, row := range rows {
		for j, value := range row {
			if len(value) > widths[j] {
				widths[j] = len(value)
			}
		}
	}

	return widths
}
