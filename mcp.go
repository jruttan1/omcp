package main

import (
	"path/filepath"
	"strings"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

func createMcp(filename string, serverInfo Info) *server.MCPServer {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	description := serverInfo.Description
	version := serverInfo.Version

	s := server.NewMCPServer(
		name+" omcp server",
		version,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
		server.WithInstructions(description),
	)

	return s

}

// add selected tools to the server instance
func createTools(selectedTools []Tool, server *server.MCPServer, baseURL string) {
	for _, tool := range selectedTools {
		opts := []mcp.ToolOption{
			mcp.WithDescription(tool.Summary),
		}

		// loop through a tools parameters
		for _, p := range tool.Parameters {
			switch p.Type {
			case "int", "integer", "number":
				if p.Required == true {
					opts = append(opts, mcp.WithNumber(p.Name, mcp.Required()))
				} else {
					opts = append(opts, mcp.WithNumber(p.Name))
				}
			case "bool", "boolean":
				if p.Required == true {
					opts = append(opts, mcp.WithBoolean(p.Name, mcp.Required()))
				} else {
					opts = append(opts, mcp.WithBoolean(p.Name))
				}
			default:
				if p.Required == true {
					opts = append(opts, mcp.WithString(p.Name, mcp.Required()))
				} else {
					opts = append(opts, mcp.WithString(p.Name))
				}
			}
		}
		server.AddTool(mcp.NewTool(tool.Name, opts...), makeHandler(tool, baseURL)) // once parameters are accumulated, pass to add tool
	}
}

func startServer(s *server.MCPServer) error {
	return server.ServeStdio(s)
}
