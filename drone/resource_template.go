package drone

import (
	"fmt"

	"github.com/Lucretius/terraform-provider-drone/drone/utils"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const templateNameProperty string = "name"
const templateNamespaceProperty string = "namespace"
const templateDataProperty string = "data"
const templateResourceIdExample = "KazanExpress/go_service.yaml"

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			templateNameProperty: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			templateNamespaceProperty: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			templateDataProperty: {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceTemplateCreate,
		Read:   resourceTemplateRead,
		Update: resourceTemplateUpdate,
		Delete: resourceTemplateDelete,
		Exists: resourceTemplateExists,
	}
}

func resourceTemplateCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace := getNamespace(data)
	template := constructDroneTemplateModel(data)

	createdTemplate, err := client.TemplateCreate(namespace, template)

	if err != nil {
		return fmt.Errorf("unable to create template %s", template.Name)
	}

	readTemplate(data, createdTemplate)
	return nil
}

func resourceTemplateUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, templateName, err := utils.ParseOrgId(data.Id(), templateResourceIdExample)
	if err != nil {
		return err
	}

	template := constructDroneTemplateModel(data)

	updatedTemplate, err := client.TemplateUpdate(namespace, templateName, template)
	if err != nil {
		return err
	}

	readTemplate(data, updatedTemplate)
	return nil
}

func resourceTemplateRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, templateName, err := utils.ParseOrgId(data.Id(), templateResourceIdExample)
	if err != nil {
		return err
	}

	template, err := client.Template(namespace, templateName)
	if err != nil {
		return fmt.Errorf("failed to read Drone template from namespace: %s with name: %s", namespace,
			templateName)
	}

	readTemplate(data, template)
	return nil
}

func resourceTemplateDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, templateName, err := utils.ParseOrgId(data.Id(), templateResourceIdExample)
	if err != nil {
		return err
	}

	return client.TemplateDelete(namespace, templateName)
}

func resourceTemplateExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	namespace, templateName, err := utils.ParseOrgId(data.Id(), templateResourceIdExample)
	if err != nil {
		return false, err
	}

	template, err := client.Template(namespace, templateName)
	if err != nil {
		return false, fmt.Errorf("failed to read Drone template from namespace: %s with name: %s", namespace,
			templateName)
	}

	exists := template.Name == templateName

	return exists, nil
}

func getNamespace(data *schema.ResourceData) string {
	return data.Get(templateNamespaceProperty).(string)
}

func constructDroneTemplateModel(data *schema.ResourceData) (template *drone.Template) {
	template = &drone.Template{
		Name: data.Get(templateNameProperty).(string),
		Data: data.Get(templateDataProperty).(string),
	}

	return template
}

func readTemplate(data *schema.ResourceData, template *drone.Template) {
	namespace := data.Get(templateNamespaceProperty).(string)

	data.SetId(fmt.Sprintf("%s/%s", namespace, template.Name))

	data.Set(templateNamespaceProperty, namespace)
	data.Set(templateNameProperty, template.Name)
	data.Set(templateDataProperty, template.Data)
}
