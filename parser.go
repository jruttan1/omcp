package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// represents the spec as a whole
type Spec struct {
	Info  Info                      `yaml:"info"`
	Paths map[string]map[string]any `yaml:"paths"`
	// map of string: any -> '/route': 'get'
}

// metadata info for each spec
type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// all the fields of a tool
type Tool struct {
	Name       string
	Method     string
	Route      string
	Summary    string
	Parameters []Param
}

// sub-struct of Tool
type Param struct {
	Name     string
	In       string // where you put the param: "path", "query", "body"
	Required bool
	Type     string // "string", "integer", "boolean"
}

func (t Tool) Title() string       { return t.Method + " " + t.Route }
func (t Tool) Description() string { return t.Summary }
func (t Tool) FilterValue() string { return t.Method + " " + t.Route }

func parse(fileName string) ([]Tool, Info, error) {

	contents, err := os.ReadFile(fileName) // contents = byte slice of input file text

	if err != nil {
		return nil, Info{}, err
	}

	var apiSpec Spec // init an instance of the Spec object

	err = yaml.Unmarshal(contents, &apiSpec) // takes (input, output) and returns error or nil and turns yaml into a struct

	if err != nil {
		return nil, Info{}, err
	}

	paths := apiSpec.Paths

	var tools []Tool

	for route, methodsMap := range paths {
		for method, details := range methodsMap {

			var summary string

			detailsMap, ok := details.(map[string]any) // check if data is type map[string]any. ok will be true or false
			if ok {
				summaryAny := detailsMap["summary"]
				if summaryAny != nil {
					summary = fmt.Sprint(detailsMap["summary"])
				}
			} else {
				summary = "No description provided"
			}

			var params []Param

			if paramsAny, ok := detailsMap["parameters"]; ok {
				for _, p := range paramsAny.([]any) { // type assertation, the whole list of parameters is type []any
					pMap := p.(map[string]any)
					param := Param{
						Name:     fmt.Sprint(pMap["name"]),
						In:       fmt.Sprint(pMap["in"]),
						Required: pMap["required"] == true,
					}
					if schema, ok := pMap["schema"].(map[string]any); ok {
						param.Type = fmt.Sprint(schema["type"])
					}
					params = append(params, param)
				}
			}

			cleanPath := strings.ReplaceAll(route, "/", "_")
			cleanMethod := strings.ToLower(method)
			name := cleanMethod + cleanPath

			tool := Tool{
				Name:       name,
				Route:      route,
				Method:     method,
				Summary:    summary,
				Parameters: params,
			}
			tools = append(tools, tool)
		}
	}

	return tools, apiSpec.Info, nil
}
