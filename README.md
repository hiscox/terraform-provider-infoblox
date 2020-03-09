# terraform-provider-infoblox

A terraform provider for creating records in Infoblox. Tested with Terraform 0.12.7 and Infoblox WAPI 2.10.1

Note that if someone manually updates/creates a record in Infoblox terraform will attempt to change it to what is defined in code.

## Build the Provider

Go version 1.14

```shell
go build -o terraform-provider-infoblox
```

## Configure the Provider

```terraform
provider "infoblox" {
  host         = "server.example.com"
  username     = "user"
  password     = "changme"
  wapi_version = "2.10.1"
  tls_verify   = false
}
```

## Create an A-Record

```terraform
resource "infoblox_a_record" "test" {
  ipv4addr = "192.168.13.9"
  name     = "dev.service.domain.com"
  comment  = "Test of automation"
  view     = "Internal"
}
```

## Create a Txt Record

```terraform
resource "infoblox_txt_record" "test" {
  name  = "example.service.domain.com"
  text  = "Test of automation"
  view  = "Internal"
}
```

## Create a Cname Record

```terraform
resource "infoblox_cname_record" "test" {
  name      = "alias.service.domain.com"
  canonical = "service.domain.com"
  comment   = "Test of automation"
  view      = "Internal"
}
```

## To-do

* Add validations to byte arrays in POST and PUT requests
  * Enhance logging to clearly indicate errors when constructing bodies
* Learn how to mock for `go test`
* Configureable TTLs
* Add comment field to record:txt
* Extend functionality to support other Infoblox objects

## License

Released under MIT [LICENSE](LICENSE).
