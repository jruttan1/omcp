package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// args is param.name: ai response
func buildURL(tool Tool, args map[string]any, baseURL string) string {
	requestURL := baseURL
	requestURL += tool.Route
	queryParams := url.Values{}
	for _, p := range tool.Parameters {
		arg := args[p.Name]
		switch p.In {
		case "path":
			requestURL = strings.ReplaceAll(requestURL, "{"+p.Name+"}", fmt.Sprint(arg))
		case "query":
			queryParams.Set(p.Name, fmt.Sprint(arg))
		}
	}
	if len(queryParams) > 0 {
		requestURL += "?" + queryParams.Encode()
	}

	return requestURL
}

func makeHandler(tool Tool, baseURL string) server.ToolHandlerFunc {
	// returns a custom mcp handler for each tool
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		requestURL := buildURL(tool, args, baseURL)
		method := tool.Method

		// create http request object with method type and request URL
		httpReq, err := http.NewRequestWithContext(ctx, method, requestURL, nil)
		if err != nil {
			return nil, err
		}

		// send the request
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close() // defer schedules this to run when function returns

		// reads response bytes
		r, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		response := string(r) // byte slice to string

		// returns
		return mcp.NewToolResultText(response), nil
	}
}
