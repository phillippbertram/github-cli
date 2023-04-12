package shared

import (
	"fmt"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/prompter"
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

func PrintProjectPreview(io *iostreams.IOStreams, project ProjectV2View) error {
	cs := io.ColorScheme()
	table := tableprinter.New(io)
	table.HeaderRow("Number", "Title", "Status", "Closed At", "Public", "Description", "ResourcePath", "Url")

	table.AddField(fmt.Sprintf("%d", project.Number))
	table.AddField(project.Title)
	table.AddField(fmtProjectStatus(cs, project.Closed))
	table.AddField(project.ClosedAt)
	table.AddField(fmt.Sprintf("%t", project.Public))
	table.AddField(project.ShortDescription)
	table.AddField(project.ResourcePath)
	table.AddField(project.URL)
	table.EndRow()

	err := table.Render()
	return err
}

func SelectProject(p prompter.Prompter, cs *iostreams.ColorScheme, projects []api.ProjectV2) (int, error) {
	candidates := []string{}
	for _, project := range projects {
		candidates = append(candidates, fmt.Sprintf("%s (%s)", project.Title, fmtProjectStatus(cs, project.Closed)))
	}

	selected, err := p.Select("Select project", "", candidates)

	return projects[selected].Number, err
}
