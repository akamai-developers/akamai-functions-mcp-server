package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	mcp_golang "github.com/metoro-io/mcp-golang"
)

type SearchAppArguments struct {
	Query string `json:"query" jsonschema:"required,description=A query to search apps for"`
}

type SearchResult struct {
	AppName     string `json:"appName"`
	AppId       string `json:"appId"`
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
}

func SearchApp(args SearchAppArguments) (*mcp_golang.ToolResponse, error) {
	if len(args.Query) == 0 {
		return nil, fmt.Errorf("You must provide a query")
	}
	args.Query = strings.ToLower(args.Query)
	command := []string{"aka", "info", "--format", "json"}
	out, err := spin.RunCommand(command...)
	if err != nil {
		return nil, err
	}
	var accountInfo GetAccountInfoResponse

	err = json.Unmarshal(out, &accountInfo)
	if err != nil {
		return nil, err
	}

	apps := make([]SearchResult, 0)
	for _, account := range accountInfo.AuthInfo.Accounts {
		out, err := spin.RunCommand("aka", "apps", "list", "--format", "json", "--account-id", account.Id)
		if err != nil {
			return nil, err
		}
		var appsPerAccount []GetAppsResponse
		err = json.Unmarshal(out, &appsPerAccount)
		if err != nil {
			return nil, err
		}
		for _, app := range appsPerAccount {
			if strings.Contains(strings.ToLower(app.Name), args.Query) {
				apps = append(apps, SearchResult{
					AppId:       app.Id,
					AppName:     app.Name,
					AccountId:   account.Id,
					AccountName: account.Name,
				})
			}
		}
	}
	content, err := json.Marshal(apps)
	if err != nil {
		return nil, err
	}
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(content))), nil

}

type ListAppsArguments struct {
	Account `jsonschema:"description=Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
}

type GetAppsResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListApps(args ListAppsArguments) (*mcp_golang.ToolResponse, error) {
	command := []string{"aka", "apps", "list", "--format", "json"}
	if len(args.AccountId) > 0 {
		command = append(command, "--account-id", args.AccountId)
	}
	out, err := spin.RunCommand(command...)
	if err != nil {
		return nil, err
	}
	var apps []GetAppsResponse
	err = json.Unmarshal(out, &apps)
	if err != nil {
		return nil, err
	}
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(out))), nil
}

type AppStatusResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Urls        []string `json:"urls"`
	CreatedAt   string   `json:"created_at"`
	Invocations int      `json:"invocations"`
}

type GetAppStatusArguments struct {
	App     `jsonschema:"description=Specify the desired Akamai Functions application by either providing the app name or its identifier. If not specified by the user, you can check if there is a .spin-aka/config.toml in the application folder, it contains the Application ID"`
	Account `jsonschema:"description=Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
}

type App struct {
	AppId   string `json:"appId,omitempty" jsonschema:"description=Identifier of the application running on Akamai Functions"`
	AppName string `json:"appName,omitempty" jsonschema:"description=Name of the application running on Akamai Functions"`
}

type Account struct {
	AccountId   string `json:"accountId,omitempty" jsonschema:"description=Akamai Functions Account Id"`
	AccountName string `json:"accountName,omitempty" jsonschema:"description=Akamai Functions Account Name"`
}

func (a GetAppStatusArguments) isValid() bool {
	if len(a.AppId) == 0 && len(a.AppName) == 0 {
		return false
	} else if len(a.AppId) > 0 && len(a.AppName) > 0 {
		return false
	}
	if len(a.AccountId) > 0 && len(a.AccountName) > 0 {
		return false
	}
	return true
}

func (a GetAppStatusArguments) getCommandArgs() []string {
	result := []string{}
	if len(a.AccountId) > 0 {
		result = append(result, "--account-id", a.AccountId)
	}
	if len(a.AccountName) > 0 {
		result = append(result, "--account-name", a.AccountName)
	}
	if len(a.AppId) > 0 {
		result = append(result, "--app-id", a.AppId)
	}
	if len(a.AppName) > 0 {
		result = append(result, "--app-name", a.AppName)
	}
	return result
}

const deprecatedAkamaiFunctionsDomain = "aka.fermyon.tech"

func getAppStatus(args GetAppStatusArguments) (*AppStatusResponse, error) {
	if !args.isValid() {
		return nil, fmt.Errorf("Invalid arguments provided")
	}
	command := []string{"aka", "app", "status", "--format", "json"}
	extraArgs := args.getCommandArgs()
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	out, err := spin.RunCommand(command...)
	if err != nil {
		return nil, err
	}
	var status AppStatusResponse
	err = json.Unmarshal(out, &status)
	if err != nil {
		return nil, err
	}

	n := 0
	for _, url := range status.Urls {
		if !strings.HasSuffix(url, deprecatedAkamaiFunctionsDomain) {
			status.Urls[n] = url
			n++
		}
	}

	status.Urls = status.Urls[:n]
	return &status, nil
}

func GetAppStatus(args GetAppStatusArguments) (*mcp_golang.ToolResponse, error) {
	status, err := getAppStatus(args)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(out))), nil
}

func GetAppUrl(args GetAppStatusArguments) (*mcp_golang.ToolResponse, error) {
	status, err := getAppStatus(args)
	if err != nil {
		return nil, err
	}
	for _, url := range status.Urls {
		if !strings.HasSuffix(url, deprecatedAkamaiFunctionsDomain) {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(url)), nil
		}
	}
	return nil, fmt.Errorf("Could not determine URL for app %s", status.Name)
}

type GetAppLogArguments struct {
	App      `jsonschema:"description=Specify the desired Akamai Functions application by either providing the app name or its identifier. If not specified by the user, you can check if there is a .spin-aka/config.toml in the application folder, it contains the Application ID"`
	Account  `jsonschema:"description=Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
	MaxLines int `json:"maxLines,omitempty" jsonschema:"description=Maximum number of log lines to retrieve. (Defaults to 10)"`
}

func (a GetAppLogArguments) getCommandArgs() []string {
	result := []string{}
	if len(a.AccountId) > 0 {
		result = append(result, "--account-id", a.AccountId)
	}
	if len(a.AccountName) > 0 {
		result = append(result, "--account-name", a.AccountName)
	}
	if len(a.AppId) > 0 {
		result = append(result, "--app-id", a.AppId)
	}
	if len(a.AppName) > 0 {
		result = append(result, "--app-name", a.AppName)
	}
	if a.MaxLines != 10 {
		result = append(result, "--max-lines", strconv.Itoa(a.MaxLines))
	}
	return result
}

func GetAppLogs(args GetAppLogArguments) (*mcp_golang.ToolResponse, error) {
	if args.MaxLines == 0 {
		args.MaxLines = 10
	}
	command := []string{"aka", "logs"}
	extraArgs := args.getCommandArgs()
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	out, err := spin.RunCommand(command...)
	if err != nil {
		return nil, err
	}
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(out))), nil
}

type GetAppDeploymentArguments struct {
	App     `jsonschema:"description=Specify the desired Akamai Functions application by either providing the app name or its identifier. If not specified by the user, you can check if there is a .spin-aka/config.toml in the application folder, it contains the Application ID"`
	Account `jsonschema:"description=Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
}

func (a GetAppDeploymentArguments) getCommandArgs() []string {
	result := []string{}
	if len(a.AccountId) > 0 {
		result = append(result, "--account-id", a.AccountId)
	}
	if len(a.AccountName) > 0 {
		result = append(result, "--account-name", a.AccountName)
	}
	if len(a.AppId) > 0 {
		result = append(result, "--app-id", a.AppId)
	}
	if len(a.AppName) > 0 {
		result = append(result, "--app-name", a.AppName)
	}
	return result
}

type AppHistory struct {
	EventType string `json:"event_type"`
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
}

func GetAppDeploymentHistory(args GetAppDeploymentArguments) (*mcp_golang.ToolResponse, error) {
	command := []string{"aka", "app", "history", "--format", "json"}
	extraArgs := args.getCommandArgs()
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	out, err := spin.RunCommand(command...)
	if err != nil {
		return nil, err
	}
	var history []AppHistory
	err = json.Unmarshal(out, &history)
	if err != nil {
		return nil, err
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(out))), nil
}
