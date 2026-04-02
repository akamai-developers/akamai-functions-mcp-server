package resources

import (
	"context"
	_ "embed"
	"log"
	"net/url"
	"strings"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed aka.md
var spinAkaCommandReference string

type AkamaiFunctionsResources struct {
	logger *log.Logger
}

func NewAkamaiFunctionsResources(logger *log.Logger) *AkamaiFunctionsResources {
	return &AkamaiFunctionsResources{
		logger: logger,
	}
}

const (
	mimeTypeMarkdown        = "text/markdown"
	resourceIdAkaCommandRef = "akamai-functions://docs/reference/spin-aka"
)

type CmdHelp struct {
	Command string `json:"command" jsonschema:"The 'spin aka' sub-command you want help with. For example: 'apps list' or 'apps deploy'"`
}

func (a *AkamaiFunctionsResources) RegisterAllWith(s *server.MCPServer) {
	ref := mcp.NewResource(resourceIdAkaCommandRef, "The 'spin aka' command reference",
		mcp.WithResourceDescription("Reference documentation for the `spin aka` command"),
		mcp.WithMIMEType(mimeTypeMarkdown),
	)

	s.AddResource(ref, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{MIMEType: mimeTypeMarkdown, Text: spinAkaCommandReference},
		}, nil
	})

	cmdHelp := mcp.NewResourceTemplate("aka-help://{command}", "Help for 'spin aka' sub-commands",
		mcp.WithTemplateDescription("Help documentation for the `spin aka` sub-commands"),
		mcp.WithTemplateMIMEType("text/plain"),
	)

	s.AddResourceTemplate(cmdHelp, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		cmd := extractCommandFromUri(request.Params.URI)
		cmd, err := url.QueryUnescape(cmd)
		if err != nil {
			a.logger.Printf("Error unescaping command from URI %s: %v\n", request.Params.URI, err)
			return []mcp.ResourceContents{
				mcp.TextResourceContents{MIMEType: "text/plain", Text: "Error parsing command from URI"},
			}, err
		}
		cmdParts := strings.Split(cmd, " ")
		arguments := append([]string{"aka"}, cmdParts...)
		arguments = append(arguments, "--help")
		a.logger.Printf("Running command for help resource: %v\n", arguments)
		out, err := spin.RunCommand(arguments...)
		if err != nil {
			a.logger.Printf("Error running command %s: %v\nOutput was: %s\n", arguments, err, string(out))
			return []mcp.ResourceContents{
				mcp.TextResourceContents{MIMEType: "text/plain", Text: "Error retrieving help for command " + cmd},
			}, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{MIMEType: "text/plain", Text: string(out)},
		}, nil
	})
}

func extractCommandFromUri(uri string) string {
	// Extract ID from "users://123" format
	parts := strings.Split(uri, "://")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
