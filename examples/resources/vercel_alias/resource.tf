resource "vercel_project" "flare_tf" {
  name = "flare-tf"
  git_repository {
    type = "github"
    repo = "chronark/flare"
  }
}


resource "vercel_domain" "chronark_com" {
  name = "chronark.com"
}


resource "vercel_alias" "flare" {
  project_id = vercel_project.flare_tf.id
  domain     = "flare.${vercel_domain.chronark_com.name}"
}