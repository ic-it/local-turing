package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/r3labs/diff/v3"

	"github.com/fatih/color"
)

func RenderResult(testName, expected, got, input string, err error, execTime time.Duration, strip bool) (string, error) {
	writer := &strings.Builder{}
	writer.WriteString(fmt.Sprintf("Test Case %s\n", testName))
	writer.WriteString(fmt.Sprintf("Execution Time: %v\n", execTime))
	if err != nil {
		writer.WriteString(color.RedString("TEST FAILED (%s)\n", err))
		return writer.String(), nil
	}
	expectedS := strings.Split(expected, "\n")
	gotS := strings.Split(got, "\n")
	if strip {
		for i, line := range expectedS {
			expectedS[i] = strings.TrimSpace(line)
		}
		for i, line := range gotS {
			gotS[i] = strings.TrimSpace(line)
		}
	}
	expected = strings.Join(expectedS, "\n")
	got = strings.Join(gotS, "\n")
	chl, err := diff.Diff(
		gotS,
		expectedS,
		diff.SliceOrdering(true),
		diff.AllowTypeMismatch(true),
	)
	if err != nil {
		return writer.String(), err
	}
	if len(chl) == 0 {
		writer.WriteString(color.GreenString("TEST PASSED\n"))
		return writer.String(), nil
	}
	writer.WriteString(color.RedString("TEST FAILED\n"))
	diffs := make([]string, len(expectedS))
	for _, ch := range chl {
		to := fmt.Sprintf("%v", ch.To)
		from := fmt.Sprintf("%v", ch.From)
		d := ""
		switch ch.Type {
		case "create":
			d = fmt.Sprintf("+ %s", color.GreenString(to))
		case "delete":
			d = fmt.Sprintf("- %s", color.RedString(from))
		case "update":
			d = fmt.Sprintf("%s -> %s", color.RedString(from), color.GreenString(to))
		default:
			panic("unknown change type")
		}
		for _, lines := range ch.Path {
			line, err := strconv.Atoi(lines)
			if err != nil {
				return writer.String(), err
			}
			if len(diffs) <= line {
				diffs = append(diffs, d)
			} else {
				diffs[line] = d
			}
		}
	}
	diffRes := strings.Join(diffs, "\n")
	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.AppendHeader(table.Row{"Got", "Expected", "Diff", "Inputs"})
	t.AppendRow([]interface{}{got, expected, diffRes, input})
	t.Render()
	return writer.String(), nil
}
