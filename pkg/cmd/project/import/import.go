package imports

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type ImportOptions struct {
	FilePath string
}

func NewCmdImport(f *cmdutil.Factory) *cobra.Command {

	opts := &ImportOptions{}

	cmd := &cobra.Command{
		Use:   "import [<path>] [flags]",
		Short: "Import a project from a file",
		Example: heredoc.Doc(`
			# interactive mode
			$ gh project import

			# import from a file
			$ gh project import path/to/file.json --project=1
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImport(opts)
		},
	}

	return cmd
}

func runImport(opts *ImportOptions) error {

	return nil
}
