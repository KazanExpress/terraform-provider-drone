package drone

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateConfigBasic(
					"foo",
					"provider_test.yaml",
					"kind: pipeline",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_template.template",
						"namespace",
						"foo",
					),
					resource.TestCheckResourceAttr(
						"drone_template.template",
						"name",
						"provider_test.yaml",
					),
					resource.TestCheckResourceAttr(
						"drone_template.template",
						"data",
						"kind: pipeline",
					),
				),
			},
		},
	})
}

func testTemplateConfigBasic(namespace, name, data string) string {
	return fmt.Sprintf(`
	resource "drone_template" "template" {
		namespace = "%s"
		name      = "%s"
		data     = "%s"
	}
	`, namespace, name, data)
}

func testTemplateDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_template" {
			continue
		}

		namespace := resource.Primary.Attributes[templateNamespaceProperty]
		name := resource.Primary.Attributes[templateNameExample]

		err := client.TemplateDelete(namespace, name)
		if err == nil {
			return fmt.Errorf("template still exists: %s/%s", namespace, name)
		}
	}

	return nil
}
