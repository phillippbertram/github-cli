package shared

import (
	"fmt"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/iostreams"
)

func fmtProjectStatus(cs *iostreams.ColorScheme, closed bool) string {
	if closed {
		return cs.Red("Closed")
	}
	return cs.Green("Open")
}

func PrintProjects(io *iostreams.IOStreams, projects []api.ProjectV2) error {
	cs := io.ColorScheme()
	table := tableprinter.New(io)
	table.HeaderRow("Number", "Title", "Status")

	for _, p := range projects {
		table.AddField(fmt.Sprintf("%d", p.Number))
		table.AddField(p.Title)
		table.AddField(fmtProjectStatus(cs, p.Closed))
		table.EndRow()
	}

	err := table.Render()
	return err
}
