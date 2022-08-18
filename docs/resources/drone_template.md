# drone_template Resource

Manage an org drone templates.

## Example Usage

```hcl
resource "drone_template" "go_service" {
  namespace  = "KazanExpress"
  name       = "go_service.yaml"
  data      = file(“go_service.yaml”)
}
```

## Argument Reference

* `namespace` - (Required) Organization name (e.g. `KazanExpress`).
* `name` - (Required) Template name.
* `data` - (Required) Template file content (Must be a valid json string).

~> In order to use the `drone_template` resource you must have admin privileges within your Drone environment.
