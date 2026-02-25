package tools

import (
	mcp_golang "github.com/metoro-io/mcp-golang"
)

func RegisterAllTools(server *mcp_golang.Server) error {
	err := server.RegisterTool("search_app", "Find an app using the specified query within all your Akamai Functions accounts", SearchApp)
	if err != nil {
		return err
	}
	err = server.RegisterTool("list_apps", "Get all Spin apps deployed to Akamai Functions", ListApps)
	if err != nil {
		return err
	}
	err = server.RegisterTool("list_accounts", "List all Akamai Functions accounts that the current user has access to", ListAccounts)
	if err != nil {
		return err
	}

	err = server.RegisterTool("get_app_status", "Retrieve the status for an application deployed to Akamai Functions", GetAppStatus)
	if err != nil {
		return err
	}

	err = server.RegisterTool("get_app_url", "Retrieve the public endpoint for an app deployed to Akamai Functions", GetAppUrl)
	if err != nil {
		return err
	}

	err = server.RegisterTool("get_app_logs", "Retrieve logs for a particular application deployed to an Akamai Functions account", GetAppLogs)
	if err != nil {
		return err
	}

	err = server.RegisterTool("get_app_history", "Retrieve the history of an app deployed to a particular Akamai Functions account", GetAppDeploymentHistory)
	if err != nil {
		return err
	}
	return nil
}
