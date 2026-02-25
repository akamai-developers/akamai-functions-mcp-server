package resources

import (
	_ "embed"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

//go:embed aka.md
var spinAkaCommandReference string

const (
	mimeTypeMarkdown        = "text/markdown"
	resourceIdAkaCommandRef = "akamai-functions://docs/reference/spin-aka"
)

func RegisterAllResources(server *mcp_golang.Server) error {
	err := server.RegisterResource(
		resourceIdAkaCommandRef,
		"The 'spin aka' command reference",
		"Reference documentation for the `spin aka` command",
		mimeTypeMarkdown,
		func() (*mcp_golang.ResourceResponse, error) {
			return mcp_golang.NewResourceResponse(
				mcp_golang.NewTextEmbeddedResource(
					resourceIdAkaCommandRef,
					spinAkaCommandReference,
					mimeTypeMarkdown,
				),
			), nil
		},
	)
	return err
}
