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
func createTools(selectedTools []Tool, server *server.MCPServer) {
	for _, tool := range selectedTools {
		summary := tool.Summary

	}
}
