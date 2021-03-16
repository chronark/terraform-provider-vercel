---
page_title: "vercel Provider"
subcategory: ""
description: |-
  
---

# vercel Provider



## Example Usage

```terraform
terraform {
  required_providers {
    vercel = {
      source  = "hashicorp.com/chronark/vercel"
      version = "9000.1"
    }
  }
}

provider "vercel" {
  token = "<YOUR_TOKEN>"
}
```

## Schema

### Required

- **token** (String, Sensitive)
