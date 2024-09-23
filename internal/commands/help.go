package commands

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

func HelpFunc(summaries map[string]string) CommandFunc {
	return func(_ []string) error {
		keys := []string{}

		for key := range maps.All(summaries) {
			keys = append(keys, key)
		}

		slices.Sort(keys)

		var builder strings.Builder

		builder.WriteString("\nCommands:\n")

		tableWriter := tabwriter.NewWriter(&builder, 0, 8, 0, '\t', 0)

		for _, key := range slices.All(keys) {
			fmt.Fprintf(tableWriter, "\n%s\t%s", key, summaries[key])
		}

		tableWriter.Flush()

		builder.WriteString("\n\n")

		fmt.Fprint(os.Stdout, builder.String())

		return nil
	}
}
