package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/pkg/cmd/project/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	HttpClient   func() (*http.Client, error)
	BaseRepo     func() (ghrepo.Interface, error)
	IO           *iostreams.IOStreams
	Organization string // default is "user projects"
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {

	opts := &ListOptions{
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
		IO:         f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "",
		Long:  "",
		Example: heredoc.Doc(`
			# List projects for the logged in user
			$ gh project list

			# Interactive organization selection
			$ gh project list --org

			# List projects for a specific organization
			$ gh project list --org myorg
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Organization, "org", "o", "", "List projects for an organization")

	return cmd
}

func runList(opts *ListOptions) error {

	httpClient, err := opts.HttpClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	// TODO: interactive org selection
	projects, err := shared.ListProjects(httpClient, repo, opts.Organization)

	// TODO: handle auth error and suggest to run `gh auth login --scopes project`
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("No projects found")
		return nil
	}

	shared.PrintProjects(opts.IO, projects)

	return nil
}
