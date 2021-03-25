
resource "vercel_project" "my_project" {
  name = "mercury"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }
}