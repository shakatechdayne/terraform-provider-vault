package generate

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/logical/framework"
)

// TODO this seems more for an auth resource, not generalizeable.
func TestMakeResource(t *testing.T) {
	pathName, pathItem := sample()
	info := &resourceInfo{
		packageName:"alicloud",
		pathToPackage:"resources/secrets/alicloud",
		fileName: "role",
	}

	bodyFields := getBodyFields(pathItem)

	parsedTemplate := template
	if pathItem.CreateSupported || pathItem.Post != nil {
		// TODO what about Exists? How does that factor in?
		parsedTemplate += "\n" + `		Create: resource{{packageNameTitle}}Create,`
		parsedTemplate += "\n" + `		Update: resource{{packageNameTitle}}Update,`
	}
	if pathItem.Get != nil {
		parsedTemplate += "\n" + `		Read:   resource{{packageNameTitle}}Read,`
	}
	if pathItem.Delete != nil {
		parsedTemplate += "\n" + `		Delete: resource{{packageNameTitle}}Delete,`
	}
	parsedTemplate += startSchema
	for fieldName, fieldType := range bodyFields {
		parsedTemplate += "\n" + attribute
		parsedTemplate = strings.Replace(parsedTemplate, "{{fieldName}}", fieldName, -1)
		parsedTemplate = strings.Replace(parsedTemplate, "{{fieldTerraformType}}", fieldType.TerraformType(), -1)
	}
	parsedTemplate += "\n" + endSchema

	if pathItem.CreateSupported || pathItem.Post != nil {
		parsedTemplate += "\n\n" + `func resource{{packageNameTitle}}Create(d *schema.ResourceData, meta interface{}) error {`
		parsedTemplate += "\n\n" + `	data := map[string]interface{}`
		for fieldName := range bodyFields {
			parsedTemplate += "\n" + `	data["{{fieldName}}"] = d.Get("{{fieldName}}")`
			parsedTemplate = strings.Replace(parsedTemplate, "{{fieldName}}", fieldName, -1)
		}
		parsedTemplate += strings.Replace("\n\n" + `	path := "{{pathName}}"`, "{{pathName}}", pathName, -1)
		for _, parameter := range pathItem.Parameters {
			parsedTemplate += "\n" + strings.Replace(`	{{parameterName}} := d.Get("{{parameterName}}")`, `{{parameterName}}`, parameter.Name, -1)
			parsedTemplate += "\n" + strings.Replace(`	path = strings.Replace(path, "{{{parameterName}}}", {{parameterName}}, -1)`, "{{parameterName}}", parameter.Name, -1)
		}
		parsedTemplate += "\n\n" + `	client := meta.(*api.Client)`
		parsedTemplate += "\n" + `	resp, err := client.Logical().Write(path, data)
	if err != nil {
		return err
	}`
		parsedTemplate += "\n" + `	d.SetId(resp.Auth.Accessor)
	d.Set("lease_started", time.Now().Format(time.RFC3339))
	d.Set("client_token", resp.Auth.ClientToken)`
		parsedTemplate += "\n\n" + `	return resource{{packageNameTitle}}Read(d, meta)
}`
	}

	// Create and Post both use the read method, so if either of those exist,
	// we need it.
	if pathItem.CreateSupported || pathItem.Post != nil || pathItem.Get != nil {
		parsedTemplate += "\n\n" + `func resource{{packageNameTitle}}Read(d *schema.ResourceData, meta interface{}) error {`
		parsedTemplate += "\n" + `	client := meta.(*api.Client)`

		parsedTemplate += "\n" + `	log.Printf("[DEBUG] Reading token %q", d.Id())
	resp, err := client.Auth().Token().LookupAccessor(d.Id())
	if err != nil {
		// If the token is not found (it has expired) we don't return an error
		if isExpiredTokenErr(err) {
			return nil
		}
		return fmt.Errorf("error reading token %q from Vault: %s", d.Id(), err)
	}
	if resp == nil {
		log.Printf("[DEBUG] Token %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] Read token %q", d.Id())
	if leaseExpiringSoon(d) {
		log.Printf("[DEBUG] Lease for %q expiring soon, renewing", d.Id())
		renewed, err := client.Auth().Token().Renew(d.Get("client_token").(string), d.Get("lease_duration").(int))
		if err != nil {
			log.Printf("[DEBUG] Error renewing token %q, bailing", d.Id())
		} else {
			resp = renewed
			d.Set("lease_started", time.Now().Format(time.RFC3339))
			d.Set("client_token", resp.Auth.ClientToken)
			d.SetId(resp.Auth.Accessor)
		}
	}`
		for fieldName := range bodyFields {
			parsedTemplate += "\n" + strings.Replace(`	d.Set("{{fieldName}}", resp.Data["{{fieldName}}"])`, "{{fieldName}}", fieldName, -1)
		}

		parsedTemplate += "\n" + `	return nil
}`
	}

	if pathItem.Delete != nil {
		parsedTemplate += "\n\n" + `func resource{{packageNameTitle}}Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	accessor := d.Id()

	log.Printf("[DEBUG] Revoking token %q", accessor)
	err := client.Auth().Token().RevokeAccessor(accessor)
	if err != nil {
		return fmt.Errorf("error revoking token %q", accessor)
	}
	log.Printf("[DEBUG] Revoked token %q", accessor)

	return nil
}`
	}

	parsedTemplate = strings.Replace(parsedTemplate, "{{packageName}}", info.packageName, -1)
	parsedTemplate = strings.Replace(parsedTemplate, "{{packageNameTitle}}", strings.Title(info.packageName), -1)

	fmt.Println(parsedTemplate)
}

func getBodyFields(pathItem *framework.OASPathItem) map[string]FieldType {
	if pathItem.Post == nil {
		// TODO will this cause a panic? It's below too.
		return nil
	}
	jsonContent, ok := pathItem.Post.RequestBody.Content["application/json"]
	if !ok {
		return nil
	}
	results := map[string]FieldType{}
	for fieldName, fieldInfo := range jsonContent.Schema.Properties {
		switch fieldInfo.Type {
		case "boolean":
			results[fieldName] = FieldTypeBool
		case "number":
			results[fieldName] = FieldTypeNumber
		case "string":
			results[fieldName] = FieldTypeString
		case "array":
			results[fieldName] = FieldTypeArray
		case "object":
			results[fieldName] = FieldTypeObject
		}
	}
	return results
}

// TODO still needed?
func snakeCaseToCamelCase(snake string) string {
	originalParts := strings.Split(snake, "_")
	resultingParts := make([]string, len(originalParts))
	for i, originalPart := range originalParts {
		if i == 0 {
			resultingParts[i] = originalPart
		} else {
			resultingParts[i] = strings.Title(originalPart)
		}
	}
	return strings.Join(resultingParts, "")
}

const (
	template = `package {{packageName}}

import (
    "github.com/hashicorp/terraform/helper/schema"
)

func resource{{packageNameTitle}}() *schema.Resource {
    return &schema.Resource{`

	startSchema = `
		Schema: map[string]*schema.Schema{`

	attribute = `			"{{fieldName}}": {
				Type:     {{fieldTerraformType}},
			},`

	endSchema = `		},
    }
}`

next = `func resourceExampleWidgetRead(d *schema.ResourceData, meta interface{}) error {
    // ... other logic ...

    d.Set("existing_attribute", /* ... */)

    // ... other logic ...
    return nil
}

func resourceExampleWidgetUpdate(d *schema.ResourceData, meta interface{}) error {
    // ... other logic ...

    existingAttribute := d.Get("existing_attribute").(string)
    // add attribute to provider update API call

    // ... other logic ...
    return resourceExampleWidgetRead(d, meta)
}
`
)

func sample() (string, *framework.OASPathItem) {
	sampleJson := `{
      "description": "Read, write and reference policies and roles that API keys or STS credentials can be made for.",
      "parameters": [{
        "name": "name",
        "description": "The name of the role.",
        "in": "path",
        "schema": {
          "type": "string"
        },
        "required": true
      }],
      "x-vault-create-supported": true,
      "get": {
        "summary": "Read, write and reference policies and roles that API keys or STS credentials can be made for.",
        "tags": ["secrets"],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      },
      "post": {
        "summary": "Read, write and reference policies and roles that API keys or STS credentials can be made for.",
        "tags": ["secrets"],
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "inline_policies": {
                    "type": "string",
                    "description": "JSON of policies to be dynamically applied to users of this role."
                  },
                  "max_ttl": {
                    "type": "number",
                    "description": "The maximum allowed lifetime of tokens issued using this role.",
                    "format": "seconds"
                  },
                  "remote_policies": {
                    "type": "array",
                    "description": "The name and type of each remote policy to be applied. Example: \"name:AliyunRDSReadOnlyAccess,type:System\".",
                    "items": {
                      "type": "string"
                    }
                  },
                  "role_arn": {
                    "type": "string",
                    "description": "ARN of the role to be assumed. If provided, inline_policies and remote_policies should be blank. At creation time, this role must have configured trusted actors, and the access key and secret that will be used to assume the role (in /config) must qualify as a trusted actor."
                  },
                  "ttl": {
                    "type": "number",
                    "description": "Duration in seconds after which the issued token should expire. Defaults to 0, in which case the value will fallback to the system/mount defaults.",
                    "format": "seconds"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      },
      "delete": {
        "summary": "Read, write and reference policies and roles that API keys or STS credentials can be made for.",
        "tags": ["secrets"],
        "responses": {
          "204": {
            "description": "empty body"
          }
        }
      }
    }`

	item := &framework.OASPathItem{}
	json.NewDecoder(strings.NewReader(sampleJson)).Decode(item)
	return "/alicloud/role/{name}", item
}

type FieldType int

const (
	FieldTypeBool FieldType = iota
	FieldTypeNumber
	FieldTypeString
	FieldTypeArray
	FieldTypeObject
)

func (t FieldType) String() string {
	switch t {
	case FieldTypeBool:
		return "boolean"
	case FieldTypeNumber:
		return "number"
	case FieldTypeString:
		return "string"
	case FieldTypeArray:
		return "array"
	case FieldTypeObject:
		return "object"
	}
	return ""
}

func (t FieldType) TerraformType() string {
	switch t {
	case FieldTypeBool:
		return "schema.TypeBool"
	case FieldTypeNumber:
		return "schema.TypeInt"
	case FieldTypeString:
		return "schema.TypeString"
	case FieldTypeArray:
		return "schema.TypeList"
	case FieldTypeObject:
		return "schema.TypeMap"
	}
	return ""
}