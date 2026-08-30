package mocklet

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIPreview struct {
	Endpoints []OpenAPIPreviewEndpoint `json:"endpoints"`
	Warnings  []string                 `json:"warnings,omitempty"`
}

type OpenAPIPreviewEndpoint struct {
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	StatusCode  int        `json:"status_code"`
	Body        string     `json:"body"`
	ContentType string     `json:"content_type"`
	Scenarios   []Scenario `json:"scenarios,omitempty"`
}

func previewOpenAPI(data []byte) (OpenAPIPreview, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = false
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return OpenAPIPreview{}, fmt.Errorf("invalid OpenAPI document: %w", err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		return OpenAPIPreview{}, fmt.Errorf("OpenAPI validation failed: %w", err)
	}
	preview := OpenAPIPreview{}
	paths := doc.Paths.Map()
	pathNames := make([]string, 0, len(paths))
	for path := range paths {
		pathNames = append(pathNames, path)
	}
	sort.Strings(pathNames)
	for _, path := range pathNames {
		if !pathPattern.MatchString(path) {
			return OpenAPIPreview{}, fmt.Errorf("unsupported path %q", path)
		}
		item := paths[path]
		methodNames := make([]string, 0, len(item.Operations()))
		for method := range item.Operations() {
			methodNames = append(methodNames, method)
		}
		sort.Strings(methodNames)
		for _, method := range methodNames {
			if !methods[strings.ToUpper(method)] {
				continue
			}
			operation := item.GetOperation(method)
			if operation == nil || operation.Responses == nil {
				continue
			}
			responses := operation.Responses.Map()
			codes := make([]int, 0, len(responses))
			for code := range responses {
				status, parseErr := strconv.Atoi(code)
				if parseErr != nil || status < 100 || status > 599 {
					continue
				}
				codes = append(codes, status)
			}
			sort.Ints(codes)
			if len(codes) == 0 {
				continue
			}
			base := OpenAPIPreviewEndpoint{Method: strings.ToUpper(method), Path: path, StatusCode: codes[0], ContentType: "application/json", Body: `{}`}
			for _, status := range codes {
				response := responses[strconv.Itoa(status)].Value
				if response == nil {
					continue
				}
				contentType, body := responseBody(response)
				scenario := Scenario{Name: fmt.Sprintf("status-%d", status), IsDefault: status == codes[0], StatusCode: status, ContentType: contentType, Body: body}
				if status == codes[0] {
					base.StatusCode, base.ContentType, base.Body = status, contentType, body
				}
				base.Scenarios = append(base.Scenarios, scenario)
			}
			preview.Endpoints = append(preview.Endpoints, base)
			if len(preview.Endpoints) > maxEndpoints {
				return OpenAPIPreview{}, fmt.Errorf("OpenAPI import is limited to %d endpoints", maxEndpoints)
			}
		}
	}
	if len(preview.Endpoints) == 0 {
		return OpenAPIPreview{}, fmt.Errorf("OpenAPI document contains no supported response routes")
	}
	return preview, nil
}

func responseBody(response *openapi3.Response) (string, string) {
	if response.Content == nil || len(response.Content) == 0 {
		return "application/json", `{}`
	}
	contentTypes := make([]string, 0, len(response.Content))
	for contentType := range response.Content {
		contentTypes = append(contentTypes, contentType)
	}
	sort.Strings(contentTypes)
	contentType := contentTypes[0]
	media := response.Content[contentType]
	if media.Example != nil {
		return contentType, marshalExample(media.Example)
	}
	exampleNames := make([]string, 0, len(media.Examples))
	for name := range media.Examples {
		exampleNames = append(exampleNames, name)
	}
	sort.Strings(exampleNames)
	for _, name := range exampleNames {
		if media.Examples[name] != nil && media.Examples[name].Value != nil && media.Examples[name].Value.Value != nil {
			return contentType, marshalExample(media.Examples[name].Value.Value)
		}
	}
	if media.Schema != nil && media.Schema.Value != nil {
		schema := media.Schema.Value
		if schema.Example != nil {
			return contentType, marshalExample(schema.Example)
		}
		if schema.Default != nil {
			return contentType, marshalExample(schema.Default)
		}
		return contentType, marshalExample(generateSchemaExample(schema))
	}
	return contentType, `{}`
}

func marshalExample(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return `{}`
	}
	return string(data)
}

func generateSchemaExample(schema *openapi3.Schema) any {
	if len(schema.Enum) > 0 {
		return schema.Enum[0]
	}
	if schema.Type == nil {
		return map[string]any{}
	}
	switch {
	case schema.Type.Is("object"):
		result := map[string]any{}
		propertyNames := make([]string, 0, len(schema.Properties))
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)
		for _, name := range propertyNames {
			if schema.Properties[name] != nil && schema.Properties[name].Value != nil {
				result[name] = generateSchemaExample(schema.Properties[name].Value)
			}
		}
		return result
	case schema.Type.Is("array"):
		if schema.Items != nil && schema.Items.Value != nil {
			return []any{generateSchemaExample(schema.Items.Value)}
		}
		return []any{}
	case schema.Type.Is("integer"):
		return 0
	case schema.Type.Is("number"):
		return 0
	case schema.Type.Is("boolean"):
		return false
	default:
		return ""
	}
}
