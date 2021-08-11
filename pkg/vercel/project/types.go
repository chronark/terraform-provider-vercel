package project

import "github.com/chronark/terraform-provider-vercel/pkg/vercel/env"

type Domain struct {
	Name string `json:"name"`
	// Redirect           string `json:"redirect"`
	// RedirectStatusCode int    `json:"redirectStatusCode"`
	// GitBranch          string `json:"gitBranch"`
}

// Project houses all the information vercel offers about a project via their api
type Project struct {
	AccountID string `json:"accountId"`
	Alias     []struct {
		ConfiguredBy        string `json:"configuredBy"`
		ConfiguredChangedAt int64  `json:"configuredChangedAt"`
		CreatedAt           int64  `json:"createdAt"`
		Deployment          struct {
			Alias         []string      `json:"alias"`
			AliasAssigned int64         `json:"aliasAssigned"`
			Builds        []interface{} `json:"builds"`
			CreatedAt     int64         `json:"createdAt"`
			CreatedIn     string        `json:"createdIn"`
			Creator       struct {
				UID         string `json:"uid"`
				Email       string `json:"email"`
				Username    string `json:"username"`
				GithubLogin string `json:"githubLogin"`
			} `json:"creator"`
			DeploymentHostname string `json:"deploymentHostname"`
			Forced             bool   `json:"forced"`
			ID                 string `json:"id"`
			Meta               struct {
				GithubCommitRef         string `json:"githubCommitRef"`
				GithubRepo              string `json:"githubRepo"`
				GithubOrg               string `json:"githubOrg"`
				GithubCommitSha         string `json:"githubCommitSha"`
				GithubRepoID            string `json:"githubRepoId"`
				GithubCommitMessage     string `json:"githubCommitMessage"`
				GithubCommitAuthorLogin string `json:"githubCommitAuthorLogin"`
				GithubDeployment        string `json:"githubDeployment"`
				GithubCommitOrg         string `json:"githubCommitOrg"`
				GithubCommitAuthorName  string `json:"githubCommitAuthorName"`
				GithubCommitRepo        string `json:"githubCommitRepo"`
				GithubCommitRepoID      string `json:"githubCommitRepoId"`
			} `json:"meta"`
			Name       string `json:"name"`
			Plan       string `json:"plan"`
			Private    bool   `json:"private"`
			ReadyState string `json:"readyState"`
			Target     string `json:"target"`
			TeamID     string `json:"teamId"`
			Type       string `json:"type"`
			URL        string `json:"url"`
			UserID     string `json:"userId"`
			WithCache  bool   `json:"withCache"`
		} `json:"deployment"`
		Domain      string `json:"domain"`
		Environment string `json:"environment"`
		Target      string `json:"target"`
	} `json:"alias"`
	Analytics struct {
		ID         string `json:"id"`
		EnabledAt  int64  `json:"enabledAt"`
		DisabledAt int64  `json:"disabledAt"`
		CanceledAt int64  `json:"canceledAt"`
	} `json:"analytics"`
	AutoExposeSystemEnvs            bool      `json:"autoExposeSystemEnvs"`
	BuildCommand                    string    `json:"buildCommand"`
	CreatedAt                       int64     `json:"createdAt"`
	DevCommand                      string    `json:"devCommand"`
	DirectoryListing                bool      `json:"directoryListing"`
	Env                             []env.Env `json:"env"`
	Framework                       string    `json:"framework"`
	ID                              string    `json:"id"`
	InstallCommand                  string    `json:"installCommand"`
	Name                            string    `json:"name"`
	NodeVersion                     string    `json:"nodeVersion"`
	OutputDirectory                 string    `json:"outputDirectory"`
	PublicSource                    bool      `json:"publicSource"`
	RootDirectory                   string    `json:"rootDirectory"`
	ServerlessFunctionRegion        string    `json:"serverlessFunctionRegion"`
	SourceFilesOutsideRootDirectory bool      `json:"sourceFilesOutsideRootDirectory"`
	UpdatedAt                       int64     `json:"updatedAt"`
	Link                            struct {
		Type             string        `json:"type"`
		Repo             string        `json:"repo"`
		RepoID           int           `json:"repoId"`
		Org              string        `json:"org"`
		GitCredentialID  string        `json:"gitCredentialId"`
		CreatedAt        int64         `json:"createdAt"`
		UpdatedAt        int64         `json:"updatedAt"`
		Sourceless       bool          `json:"sourceless"`
		ProductionBranch string        `json:"productionBranch"`
		DeployHooks      []interface{} `json:"deployHooks"`
	} `json:"link"`
	LatestDeployments []struct {
		Alias         []string      `json:"alias"`
		AliasAssigned int64         `json:"aliasAssigned"`
		Builds        []interface{} `json:"builds"`
		CreatedAt     int64         `json:"createdAt"`
		CreatedIn     string        `json:"createdIn"`
		Creator       struct {
			UID         string `json:"uid"`
			Email       string `json:"email"`
			Username    string `json:"username"`
			GithubLogin string `json:"githubLogin"`
		} `json:"creator"`
		DeploymentHostname string `json:"deploymentHostname"`
		Forced             bool   `json:"forced"`
		ID                 string `json:"id"`
		Meta               struct {
			GithubCommitRef         string `json:"githubCommitRef"`
			GithubRepo              string `json:"githubRepo"`
			GithubOrg               string `json:"githubOrg"`
			GithubCommitSha         string `json:"githubCommitSha"`
			GithubCommitAuthorLogin string `json:"githubCommitAuthorLogin"`
			GithubCommitMessage     string `json:"githubCommitMessage"`
			GithubRepoID            string `json:"githubRepoId"`
			GithubDeployment        string `json:"githubDeployment"`
			GithubCommitOrg         string `json:"githubCommitOrg"`
			GithubCommitAuthorName  string `json:"githubCommitAuthorName"`
			GithubCommitRepo        string `json:"githubCommitRepo"`
			GithubCommitRepoID      string `json:"githubCommitRepoId"`
		} `json:"meta"`
		Name       string      `json:"name"`
		Plan       string      `json:"plan"`
		Private    bool        `json:"private"`
		ReadyState string      `json:"readyState"`
		Target     interface{} `json:"target"`
		TeamID     string      `json:"teamId"`
		Type       string      `json:"type"`
		URL        string      `json:"url"`
		UserID     string      `json:"userId"`
		WithCache  bool        `json:"withCache"`
	} `json:"latestDeployments"`
	Targets struct {
		Production struct {
			Alias         []string      `json:"alias"`
			AliasAssigned int64         `json:"aliasAssigned"`
			Builds        []interface{} `json:"builds"`
			CreatedAt     int64         `json:"createdAt"`
			CreatedIn     string        `json:"createdIn"`
			Creator       struct {
				UID         string `json:"uid"`
				Email       string `json:"email"`
				Username    string `json:"username"`
				GithubLogin string `json:"githubLogin"`
			} `json:"creator"`
			DeploymentHostname string `json:"deploymentHostname"`
			Forced             bool   `json:"forced"`
			ID                 string `json:"id"`
			Meta               struct {
				GithubCommitRef         string `json:"githubCommitRef"`
				GithubRepo              string `json:"githubRepo"`
				GithubOrg               string `json:"githubOrg"`
				GithubCommitSha         string `json:"githubCommitSha"`
				GithubRepoID            string `json:"githubRepoId"`
				GithubCommitMessage     string `json:"githubCommitMessage"`
				GithubCommitAuthorLogin string `json:"githubCommitAuthorLogin"`
				GithubDeployment        string `json:"githubDeployment"`
				GithubCommitOrg         string `json:"githubCommitOrg"`
				GithubCommitAuthorName  string `json:"githubCommitAuthorName"`
				GithubCommitRepo        string `json:"githubCommitRepo"`
				GithubCommitRepoID      string `json:"githubCommitRepoId"`
			} `json:"meta"`
			Name       string `json:"name"`
			Plan       string `json:"plan"`
			Private    bool   `json:"private"`
			ReadyState string `json:"readyState"`
			Target     string `json:"target"`
			TeamID     string `json:"teamId"`
			Type       string `json:"type"`
			URL        string `json:"url"`
			UserID     string `json:"userId"`
			WithCache  bool   `json:"withCache"`
		} `json:"production"`
	} `json:"targets"`
}

// CreateProject has all the fields the user can set when creating a new project
type CreateProject struct {
	Name          string `json:"name"`
	GitRepository struct {
		Type string `json:"type"`
		Repo string `json:"repo"`
	} `json:"gitRepository,omitempty"`
	UpdateProject
}

// UpdateProject has all the values a user can update without recreating a project
// https://vercel.com/docs/api#endpoints/projects/update-a-single-project
type UpdateProject struct {
	// The framework that is being used for this project. When null is used no framework is selected.
	Framework string `json:"framework,omitempty"`

	// Specifies whether the source code and logs of the deployments for this project should be public or not.
	PublicSource bool `json:"publicSource,omitempty"`

	// The install command for this project. When null is used this value will be automatically detected.
	InstallCommand string `json:"installCommand,omitempty"`

	// The build command for this project. When null is used this value will be automatically detected.
	BuildCommand string `json:"buildCommand,omitempty"`

	// The dev command for this project. When null is used this value will be automatically detected.
	DevCommand string `json:"devCommand,omitempty"`

	// The output directory of the project. When null is used this value will be automatically detected.
	OutputDirectory string `json:"outputDirectory,omitempty"`

	// The region to deploy Serverless Functions in this project.
	ServerlessFunctionRegion string `json:"serverlessFunctionRegion,omitempty"`

	// The name of a directory or relative path to the source code of your project. When null is used it will default to the project root.
	RootDirectory string `json:"rootDirectory,omitempty"`

	// A new name for this project.
	Name string `json:"name,omitempty"`

	// The Node.js Version for this project.
	NodeVersion string `json:"nodeVersion,omitempty"`
}
