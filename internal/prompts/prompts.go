package prompts

import (
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

type FindAppByNameArguments struct {
	Query string `json:"query" jsonschema:"description=The name of the app you're looking for"`
}

func RegisterAllPrompts(server *mcp_golang.Server) error {
	err := server.RegisterPrompt(
		"find_app_by_name",
		"Provides instructions to retrieve an application by looking at all accounts you've access to", // Description
		func(args FindAppByNameArguments) (*mcp_golang.PromptResponse, error) {

			instruction := fmt.Sprintf("Use the `search_app` tool using the following query: %s", args.Query)
			instruction += "\nAfter fetching, provide a summarized list including the app name, app id, account name and account id"

			// Construct the response
			return mcp_golang.NewPromptResponse(
				"Find an app deployed to any of the Akamai Functions accounts you've access to", // Title of the response
				mcp_golang.NewPromptMessage(
					mcp_golang.NewTextContent(instruction),
					mcp_golang.RoleUser, // The role the instruction appears as
				),
			), nil
		},
	)
	return err
}
