package project

import (
	"github.com/MakeNowJust/heredoc"
	importCmd "github.com/cli/cli/v2/pkg/cmd/project/import"
	listCmd "github.com/cli/cli/v2/pkg/cmd/project/list"
	vieCmd "github.com/cli/cli/v2/pkg/cmd/project/view"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdProject(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "",
		Long:  "",
		Example: heredoc.Doc(`
			This command needs additional permissions.
			In order to use this command, you must authenticate with the following scopes.
			$ gh auth login --scopes "project"

			# List projects for the logged in user
			$ gh project list

			# Interactive organization selection
			$ gh project list --org

			# List projects for a specific organization
			$ gh project list --org myorg

			# Import a project from a CSV file
			$ gh project import --file myproject.csv

			# View a project
			$ gh project view 1
		`),
	}

	cmd.AddCommand(listCmd.NewCmdList(f))
	cmd.AddCommand(importCmd.NewCmdImport(f))
	cmd.AddCommand(vieCmd.NewCmdView(f))

	return cmd
}
