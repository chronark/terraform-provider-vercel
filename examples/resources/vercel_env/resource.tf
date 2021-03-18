
resource "vercel_project" "my_project" {
  // ...
}

resource "vercel_env" "my_env" {
  project_id = vercel_project.my_project.id // or use a hardcoded value of an existing project
  type       = "plain"
  key        = "hello"
  value      = "world"
  target     = ["production", "preview", "development"]
}