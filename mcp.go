package main

import (
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/server"
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
