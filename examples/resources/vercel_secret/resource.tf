



// Create a new secret
resource "vercel_secret" "my_secret" {
  name  = "my_secret_name" // `@` will be prefixed automatically
  value = "super secret"
}

// Use the secret
resource "vercel_env" "env" {
  type  = "secret"
  key   = "Hello"
  value = vercel_secret.my_secret.id

  //  irrelevant values omitted
}
