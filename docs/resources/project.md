---
page_title: "vercel_project Resource - terraform-provider-vercel"
subcategory: ""
description: |-
  Sample resource in the Terraform provider scaffolding.
---

# Resource `vercel_project`

Sample resource in the Terraform provider scaffolding.

## Example Usage

```terraform
resource "vercel_project" "my_project" {
  name = "mercury"
  git_repository {
    type = "github"
    repo = "chronark/mercury"
  }
}
```

## Schema

### Required

- **git_repository** (Block List, Min: 1, Max: 1) The Git Repository that will be connected to the project. Any pushes to the specified connected Git Repository will be automatically deployed. (see [below for nested schema](#nestedblock--git_repository))
- **name** (String) The name of the project.

### Optional

- **build_command** (String) The build command for this project. When null is used this value will be automatically detected.
- **dev_command** (String) The dev command for this project. When null is used this value will be automatically detected.
- **env** (Block List) A list of environment variables configured for the project. (see [below for nested schema](#nestedblock--env))
- **framework** (String) The framework that is being used for this project. When null is used no framework is selected.
- **id** (String) The ID of this resource.
- **install_command** (String) The install command for this project. When null is used this value will be automatically detected.
- **node_version** (String) The Node.js Version for this project.
- **output_directory** (String) The output directory of the project. When null is used this value will be automatically detected.
- **public_source** (Boolean) Specifies whether the source code and logs of the deployments for this project should be public or not.
- **root_directory** (String) The name of a directory or relative path to the source code of your project. When null is used it will default to the project root.
- **serverless_function_region** (String) The region to deploy Serverless Functions in this project.

### Read-only

- **account_id** (String) The unique ID of the user or team the project belongs to.
- **created_at** (Number) A number containing the date when the project was created in milliseconds.
- **updated_at** (Number) A number containing the date when the project was updated in milliseconds.

<a id="nestedblock--git_repository"></a>
### Nested Schema for `git_repository`

Required:

- **repo** (String) The name of the Git Repository.
- **type** (String) The Git Provider of the repository. Must be either `github`, `gitlab`, or `bitbucket`.


<a id="nestedblock--env"></a>
### Nested Schema for `env`

Optional:

- **key** (String) The name of this variable
- **type** (String) `secret` or `public`
- **value** (String) The value of this variable.

Read-only:

- **created_at** (Number) A number containing the date when the variable was created in milliseconds.
- **id** (String) Unique id for this variable.
- **updated_at** (Number) A number containing the date when the variable was updated in milliseconds.


