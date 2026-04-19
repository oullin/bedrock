package prompts

// Table displays tabular data with headers and rows.
// If headers is empty, rows are displayed without a header row.
func Table(headers []string, rows [][]string) {
	w := getWriter()
	w.Write(getTheme().TableRenderer(headers, rows))
}
