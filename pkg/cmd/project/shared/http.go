package shared

import (
	"fmt"
	"net/http"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/ghrepo"
	"github.com/shurcooL/githubv4"
)

func ListProjects(client *http.Client, repo ghrepo.Interface, org string) ([]api.ProjectV2, error) {
	apiClient := api.NewClientFromHTTP(client)

	if len(org) > 0 {
		owner := ghrepo.New(org, repo.RepoName())
		return api.OrganizationProjectsV2(apiClient, owner)
	}

	return api.CurrentUserProjectsV2(apiClient, repo.RepoHost())
}

func CurrentUserProjects(client *http.Client, repo ghrepo.Interface) ([]api.ProjectV2, error) {
	apiClient := api.NewClientFromHTTP(client)
	return CurrentUserProjectsV2(apiClient, repo.RepoHost(), nil)
}

// TODO: move this to api
// if query is nil, all projects are returned (open and closed)
// if query is not nil, only projects matching the query are returned ("is:open" -> open and "is:closed" -> closed)
func CurrentUserProjectsV2(client *api.Client, hostname string, query *string) ([]api.ProjectV2, error) {
	type responseData struct {
		Viewer struct {
			ProjectsV2 struct {
				Nodes    []api.ProjectV2
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
			} `graphql:"projectsV2(first: 100, orderBy: {field: TITLE, direction: ASC}, after: $endCursor, query: $query)"`
		} `graphql:"viewer"`
	}

	variables := map[string]interface{}{
		"endCursor": (*githubv4.String)(nil),
		"query":     (*githubv4.String)(nil),
	}

	if query != nil {
		variables["query"] = githubv4.String(*query)
	}

	var projectsV2 []api.ProjectV2
	for {
		var query responseData
		err := client.Query(hostname, "UserProjectV2List", &query, variables)
		if err != nil {
			return nil, err
		}

		projectsV2 = append(projectsV2, query.Viewer.ProjectsV2.Nodes...)

		if !query.Viewer.ProjectsV2.PageInfo.HasNextPage {
			break
		}
		variables["endCursor"] = githubv4.String(query.Viewer.ProjectsV2.PageInfo.EndCursor)
	}

	return projectsV2, nil
}

func ListOrganizationProjects(client *http.Client, repo ghrepo.Interface, org string) ([]api.ProjectV2, error) {
	apiClient := api.NewClientFromHTTP(client)
	owner := ghrepo.New(org, repo.RepoName())
	return api.OrganizationProjectsV2(apiClient, owner)
}

func ListRepoProjects(client *http.Client, repo ghrepo.Interface) ([]api.ProjectV2, error) {
	apiClient := api.NewClientFromHTTP(client)
	return api.RepoProjectsV2(apiClient, repo)
}

func ListAllProjects(client *http.Client, repo ghrepo.Interface, org string) ([]api.ProjectV2, error) {
	var projects = []api.ProjectV2{}

	uProjects, error := CurrentUserProjects(client, repo)
	if error != nil {
		return projects, error
	}
	projects = append(projects, uProjects...)

	// TODO: get all orgs and iterate over them?
	if len(org) > 0 {
		oProjects, error := ListOrganizationProjects(client, repo, org)
		if error != nil {
			return projects, error
		}
		projects = append(projects, oProjects...)
	}

	rProjects, error := ListRepoProjects(client, repo)
	if error != nil {
		return projects, error
	}
	projects = append(projects, rProjects...)

	// TODO: filter out duplicates?

	return projects, nil

}

func GetProjectFromOrg(client *api.Client, repo ghrepo.Interface, org string, projectNumber int) (*ProjectV2View, error) {
	type responseData struct {
		Organization struct {
			ProjectV2 struct {
				ProjectV2View
			} `graphql:"projectV2(number: $number)"`
		} `graphql:"organization(login: $owner)"`
	}

	orgRepo := ghrepo.New(org, repo.RepoName())
	variables := map[string]interface{}{
		"owner":  githubv4.String(orgRepo.RepoOwner()),
		"number": githubv4.Int(projectNumber),
	}

	var query responseData
	err := client.Query(repo.RepoHost(), "OrganizationProjectV2View", &query, variables)
	if err != nil {
		return nil, err
	}

	projectV2 := query.Organization.ProjectV2.ProjectV2View
	return &projectV2, nil
}

func GetProjectFromRepo(client *api.Client, repo ghrepo.Interface, projectNumber int) (*ProjectV2View, error) {
	type responseData struct {
		Repository struct {
			ProjectV2 struct {
				ProjectV2View
			} `graphql:"projectV2(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(repo.RepoOwner()),
		"name":   githubv4.String(repo.RepoName()),
		"number": githubv4.Int(projectNumber),
	}

	var query responseData
	err := client.Query(repo.RepoHost(), "RepoProjectV2View", &query, variables)
	if err != nil {
		return nil, err
	}

	projectView := query.Repository.ProjectV2.ProjectV2View
	return &projectView, nil
}

func GetProjectFromCurrentUser(client *api.Client, repo ghrepo.Interface, projectNumber int) (*ProjectV2View, error) {
	type responseData struct {
		Viewer struct {
			ProjectV2 struct {
				ProjectV2View
			} `graphql:"projectV2(number: $number)"`
		} `graphql:"viewer"`
	}

	variables := map[string]interface{}{
		"number": githubv4.Int(projectNumber),
	}

	var query responseData
	err := client.Query(repo.RepoHost(), "UserProjectV2View", &query, variables)
	if err != nil {
		return nil, err
	}

	projectView := query.Viewer.ProjectV2.ProjectV2View
	return &projectView, nil
}

func GetProject(client *http.Client, repo ghrepo.Interface, org string, projectNumber int) (*ProjectV2View, error) {
	apiClient := api.NewClientFromHTTP(client)

	// try get from repo
	project, cuErr := GetProjectFromCurrentUser(apiClient, repo, projectNumber)
	if cuErr == nil && project != nil {
		return project, nil
	}

	// try get from repo
	project, repErr := GetProjectFromRepo(apiClient, repo, projectNumber)
	if repErr == nil && project != nil {
		return project, nil
	}

	// try get from org
	project, orgErr := GetProjectFromOrg(apiClient, repo, org, projectNumber)
	if orgErr == nil && project != nil {
		return project, nil
	}

	// prepare error
	err := fmt.Errorf("project not found")
	if cuErr != nil {
		err = fmt.Errorf("%w %w", err, cuErr)
	}
	if repErr != nil {
		err = fmt.Errorf("%w %w", err, repErr)
	}
	if orgErr != nil {
		err = fmt.Errorf("%w %w", err, orgErr)
	}

	return nil, err
}

type ProjectV2View struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Number           int    `json:"number"`
	ResourcePath     string `json:"resourcePath"`
	Closed           bool   `json:"closed"`
	ClosedAt         string `json:"closedAt"`
	Public           bool   `json:"public"`
	ShortDescription string `json:"shortDescription"`
	URL              string `json:"url"`
}
