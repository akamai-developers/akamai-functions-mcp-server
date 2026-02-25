package main

import (
	"github.com/ThorstenHans/akamai-functions-mcp/internal/prompts"
	"github.com/ThorstenHans/akamai-functions-mcp/internal/resources"
	"github.com/ThorstenHans/akamai-functions-mcp/internal/tools"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport(), mcp_golang.WithName("Akamai Functions MCP Server"))
	err := tools.RegisterAllTools(server)

	if err != nil {
		panic(err)
	}

	err = prompts.RegisterAllPrompts(server)

	if err != nil {
		panic(err)
	}

	err = resources.RegisterAllResources(server)

	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
