

resource "vercel_domain" "chronark_com" {
  name = "chronark.com"
}


resource "vercel_dns" "www" {
  domain = vercel_domain.chronark_com.name
  type   = "CNAME"
  value  = "www.${vercel_domain.chronark_com.name}"
  name   = "www"
}