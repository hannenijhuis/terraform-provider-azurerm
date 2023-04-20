package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorRuleSetsAssociationResource struct {
	// The 'azurerm_cdn_frontdoor_rule_set' configuration block(s) for the Front Door Rule Set Resource(s)
	RuleSetConfig string
}

func NewCdnFrontDoorRuleSetsAssociationResource(ruleSetConfig string) *CdnFrontDoorRuleSetsAssociationResource {
	return &CdnFrontDoorRuleSetsAssociationResource{
		RuleSetConfig: ruleSetConfig,
	}
}

func TestAccCdnFrontDoorRuleSetsAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_sets_association", "test")
	r := NewCdnFrontDoorRuleSetsAssociationResource(templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("1"),
			),
		},
	})
}

// NOTE: the 'requiresImport' test is not possible with this resource

func TestAccCdnFrontDoorRuleSetsAssociation_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_sets_association", "test")
	r := NewCdnFrontDoorRuleSetsAssociationResource(templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("1"),
			),
		},
		{
			Config: r.destroy(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func TestAccCdnFrontDoorRuleSetsAssociation_removeRouteRuleSetAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_sets_association", "test")
	r := NewCdnFrontDoorRuleSetsAssociationResource(templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("1"),
			),
		},
		{
			Config: r.removeRouteRuleSetAssociations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("0"),
			),
		},
	})
}

func TestAccCdnFrontDoorRuleSetsAssociation_removeMultipleRouteRuleSetAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_sets_association", "test")
	r := NewCdnFrontDoorRuleSetsAssociationResource(templateMultipleRuleSets(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRuleSets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("3"),
			),
		},
		{
			Config: r.removeRouteRuleSetAssociations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("0"),
			),
		},
	})
}

func TestAccCdnFrontDoorRuleSetsAssociation_multipleRuleSetAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_sets_association", "test")
	r := NewCdnFrontDoorRuleSetsAssociationResource(templateMultipleRuleSets(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRuleSets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("3"),
			),
		},
	})
}

func TestAccCdnFrontDoorRuleSetsAssociation_destroySingleRuleSetAssociation(t *testing.T) {
	// Regression test case for issue #20744
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_sets_association", "test")
	r := NewCdnFrontDoorRuleSetsAssociationResource(templateMultipleRuleSets(data))
	d := NewCdnFrontDoorRuleSetsAssociationResource(templateRuleSetsOneAndThreeOnly(data))
	n := NewCdnFrontDoorRuleSetsAssociationResource(templateRuleSetsOneAndThreeOnly(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRuleSets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("3"),
			),
		},
		{
			// NOTE: Need to delete one of rule sets and the rule set association,
			// this step is the actual regression part of this test case... This may
			// seem a bit confusing, but the route will still exist just the association
			// and rule set two will be destroyed with out an error being raised by the service...
			Config: d.destroy(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		{
			// NOTE: Now that the rule set and the association have been destroyed you need to add
			// the association back with the correct rule sets to be associated with the route
			// which still exist in Azure...
			Config: n.ruleSetsOneAndThreeOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(n),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("2"),
			),
		},
	})
}

func (r CdnFrontDoorRuleSetsAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorRuleSetAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorRoutesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnFrontDoorRuleSetsAssociationResource) multipleRuleSets(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_sets_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.two.id, azurerm_cdn_frontdoor_rule_set.three.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetsAssociationResource) ruleSetsOneAndThreeOnly(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_sets_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.three.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetsAssociationResource) removeRouteRuleSetAssociations(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_sets_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = []
}
`, template)
}

func (r CdnFrontDoorRuleSetsAssociationResource) destroy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
`, template)
}

func (r CdnFrontDoorRuleSetsAssociationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_sets_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetsAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-profile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "acctest-origin-group-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  session_affinity_enabled = true

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "GET"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-origin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-endpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  enabled                  = true
}

%[4]s

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "acctest-contoso-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  enabled                       = true

  https_redirect_enabled = true
  forwarding_protocol    = "HttpsOnly"
  patterns_to_match      = ["/%[3]s"]
  supported_protocols    = ["Http", "Https"]
  link_to_default_domain = true

  cache {
    compression_enabled           = true
    content_types_to_compress     = ["text/html", "text/javascript", "text/xml"]
    query_strings                 = ["account", "settings", "foo", "bar"]
    query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(10), r.RuleSetConfig)
}

func templateOneRuleSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_cdn_frontdoor_rule_set" "one" {
  name                     = "acctestrulesetone%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger)
}

func templateMultipleRuleSets(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_cdn_frontdoor_rule_set" "one" {
  name                     = "acctestrulesetone%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "two" {
  name                     = "acctestrulesettwo%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "three" {
  name                     = "acctestrulesetthree%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger)
}

func templateRuleSetsOneAndThreeOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_cdn_frontdoor_rule_set" "one" {
  name                     = "acctestrulesetone%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "three" {
  name                     = "acctestrulesetthree%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger)
}
