
// This is optional and only useful to reference a project id.
resource "vercel_project" "my_project" {
  name = "mercury-via-terraform"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }
}

resource "vercel_env" "my_env" {
  project_id = vercel_project.my_project.id // or use a hardcoded value of an existing project
  type       = "plain"
  key        = "hello"
  value      = "world"
  target     = ["production", "preview", "development"]
}