package main

import (
	"log"
	"os"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/resources"
	"github.com/ThorstenHans/akamai-functions-mcp/internal/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Setup debug logging
	logFile, err := os.OpenFile("mcp-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "[MCP] ", log.LstdFlags|log.Lshortfile)
	afTools := tools.NewAkamaiFunctionsTools(logger)
	afResources := resources.NewAkamaiFunctionsResources(logger)
	mcpServer := server.NewMCPServer(
		"Akamai Functions MCP Server",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithPromptCapabilities(false),
		server.WithPromptCapabilities(false),
		server.WithRecovery(),
		server.WithLogging(),
	)

	afTools.RegisterAllWith(mcpServer)
	afResources.RegisterAllWith(mcpServer)

	//	prompts.RegisterAll(mcpServer, logger)
	//resources.RegisterAll(mcpServer, logger)

	if err := server.ServeStdio(mcpServer); err != nil {
		logger.Printf("Stdio server failed: %v\n", err)
	}
}
