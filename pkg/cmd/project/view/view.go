package view

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/browser"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/cli/cli/v2/internal/prompter"
	"github.com/cli/cli/v2/internal/text"
	"github.com/cli/cli/v2/pkg/cmd/project/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (ghrepo.Interface, error)
	Prompter   prompter.Prompter
	Browser    browser.Browser

	Interactive  bool
	Organization string
	ProjectId    *int
	Web          bool
}

func NewCmdView(f *cmdutil.Factory) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Prompter:   f.Prompter,
		Browser:    f.Browser,
	}

	cmd := &cobra.Command{
		Use:   "view [<project-number>]",
		Short: "View a summary of a project",
		Example: heredoc.Doc(`
			# Interactively select a project to view
			$ gh project view
			
			# View a specific project
			$ gh project view 1

			# View a specific project in the browser
			$ gh project view 1 --web
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.BaseRepo = f.BaseRepo

			if len(args) > 0 {
				projIDStr := args[0]
				projID, err := strconv.Atoi(projIDStr)
				if err != nil {
					return fmt.Errorf("invalid project number: %q", projIDStr)
				}
				opts.ProjectId = &projID
			}

			opts.Interactive = opts.ProjectId == nil

			if opts.Interactive && !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf("must provide `project-number` when not running interactively")
			}

			return runView(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Organization, "org", "o", "", "List projects for an organization")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open run in the browser")

	return cmd
}

func runView(opts *ViewOptions) error {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	defer opts.IO.StopProgressIndicator()

	projectId := opts.ProjectId

	if opts.Interactive {
		// get all projects
		opts.IO.StartProgressIndicator()
		projects, err := shared.ListAllProjects(httpClient, repo, opts.Organization)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}

		if len(projects) == 0 {
			return fmt.Errorf("no projects found")
		}

		// TODO: add information about the source of the project (repo, org, user)
		selectedProject, err := shared.SelectProject(opts.Prompter, opts.IO.ColorScheme(), projects)
		if err != nil {
			return err
		}
		projectId = &selectedProject
	}

	opts.IO.StartProgressIndicator()
	project, err := shared.GetProject(httpClient, repo, opts.Organization, *projectId)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return fmt.Errorf("failed to get run: %w", err)
	}

	if opts.Web {
		url := project.URL
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", text.DisplayURL(url))
		}

		return opts.Browser.Browse(url)
	}

	shared.PrintProjectPreview(opts.IO, *project)
	return nil
}
