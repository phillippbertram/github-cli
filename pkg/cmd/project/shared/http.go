package shared

import (
	"net/http"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghrepo"
)

func ListProjects(client *http.Client, repo ghrepo.Interface, org string) ([]api.ProjectV2, error) {
	apiClient := api.NewClientFromHTTP(client)

	if len(org) > 0 {
		owner := ghrepo.New(org, repo.RepoName())
		return api.OrganizationProjectsV2(apiClient, owner)
	}
	return api.CurrentUserProjectsV2(apiClient, repo.RepoHost())

}
