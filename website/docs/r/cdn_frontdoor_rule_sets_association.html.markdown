---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_rule_sets_association"
description: |-
  Manages the association between one or more Front Door (standard/premium) Rule Sets and a Front Door (standard/premium) Route.
---

# azurerm_cdn_frontdoor_rule_sets_association

Manages the association(s) between one or more Front Door (standard/premium) Rule Sets and a Front Door (standard/premium) Route. Use this resource to avoid receiving the service side error `This resource is still associated with a route. Please delete the association with the route first before deleting this resource` when you attempt to `destroy`/`delete` a Front Door (standard/premium) Rule Set. 

!> **IMPORTANT:** When you delete this resource it will remove **all** of the Front Door Rule Set associations from the referenced Front Door Route resource.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_rule_sets_association" "example" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.contoso.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.static.id, azurerm_cdn_frontdoor_rule_set.images.id]
}
```

## Arguments Reference

The following arguments are supported:

* `cdn_frontdoor_route_id` - (Required) The ID of the Front Door Route that should be managed by the association resource. Changing this forces a new association resource to be created.

* `cdn_frontdoor_rule_set_ids` - (Required) One or more Front Door Rule Set IDs which are associated with the Front Door Route. Changing this forces a new association resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Rule Set Association.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Rule Set Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Rule Set Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Rule Set Association.

## Import

Front Door Rule Set Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_rule_sets_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/afdEndpoint1/associations/route1
```
