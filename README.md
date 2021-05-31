
<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" height="50" align="right"></img>
</a>

# Terraform Provider for Vercel

Configure vercel resources such as projects, deployments and secrets as code with terraform.

<div align="center"><a href="https://codecov.io/gh/chronark/terraform-provider-vercel">
  <img src="https://codecov.io/gh/chronark/terraform-provider-vercel/branch/main/graph/badge.svg?token=pBJrBYgr9g"/>
</a></div>



## Features

This provider has not reached feauture parity with the vercel api yet. I am adding new features as I need them.
Please create an issue if you requrie a certain feature, I will work on them asap.

Available features can be found [here](https://registry.terraform.io/providers/chronark/vercel/latest/docs).

## Quickstart

1. Create a token [here](https://vercel.com/account/tokens)
2. Create a `vercel.tf` file with the following content. 
    - Replace `<YOUR_TOKEN>` with the token from step 1. Alternatively you can set the `VERCEL_TOKEN` environment variable
    - Change the `git_repository` to whatever you want to deploy.

```tf
terraform {
  required_providers {
    vercel = {
      source  = "hashicorp.com/chronark/vercel"
      version = "1.0.0"
    }
  }
}

provider "vercel" {
  token = "<YOUR_TOKEN>"
}

resource "vercel_project" "my_project" {
  name = "mercury-via-terraform"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }
}
```

3. Run
```sh
terraform init
terraform apply
```


4. Check vercel's [dashboard](https://vercel.com/dashboard) to see your project.
5. Push to the default branch of your repository to create your first deployment.

## Documentation

Documentation can be found [here](https://registry.terraform.io/providers/chronark/vercel/latest/docs)

## Development Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15
