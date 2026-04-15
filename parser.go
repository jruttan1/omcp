package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Info struct {
	Title   string `yaml:"title"`
	Version string `yaml:"version"`
}

type Spec struct {
	Info  Info                      `yaml:"info"`
	Paths map[string]map[string]any `yaml:"paths"`
	// map of string: any -> '/route': 'get'
}

type Tool struct {
	Name    string
	Method  string
	Route   string
	Summary string
}

func (t Tool) Title() string       { return t.Method + " " + t.Route }
func (t Tool) Description() string { return t.Summary }
func (t Tool) FilterValue() string { return t.Method + " " + t.Route }

func parse(fileName string) ([]Tool, error) {

	contents, err := os.ReadFile(fileName) // contents = byte slice of input file text

	if err != nil {
		return nil, err
	}

	var apiSpec Spec // init an instance of the Spec object

	err = yaml.Unmarshal(contents, &apiSpec) // takes (input, output) and returns error or nil and turns yaml into a struct

	if err != nil {
		return nil, err
	}

	paths := apiSpec.Paths

	var tools []Tool

	for route, methodsMap := range paths {
		for method, details := range methodsMap {

			var summary string

			summaryMap, ok := details.(map[string]any) // check if data is type map[string]any. ok will be true or false
			if ok {
				summaryAny := summaryMap["summary"]
				if summaryAny != nil {
					summary = fmt.Sprint(summaryMap["summary"])
				}
			} else {
				summary = "No description provided"
			}

			cleanPath := strings.ReplaceAll(route, "/", "_")
			cleanMethod := strings.ToLower(method)
			name := cleanMethod + cleanPath

			tool := Tool{
				Name:    name,
				Route:   route,
				Method:  method,
				Summary: summary,
			}
			tools = append(tools, tool)
		}
	}

	return tools, nil

}
