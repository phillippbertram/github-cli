package view

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type ViewOptions struct {
	FilePath string
}

func NewCmdView(f *cmdutil.Factory) *cobra.Command {

	opts := &ViewOptions{}

	cmd := &cobra.Command{
		Use:   "view <project-name>",
		Short: "",
		Long:  "",
		Example: heredoc.Doc(`
			# interactive mode
			$ gh project view my-project
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView(opts)
		},
	}

	return cmd
}

func runView(opts *ViewOptions) error {
	return nil
}
