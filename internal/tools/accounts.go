package tools

import (
	"encoding/json"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	mcp_golang "github.com/metoro-io/mcp-golang"
)

type GetAccountInfoResponse struct {
	AuthInfo AuthInfo `json:"auth_info"`
}

type AuthInfo struct {
	Accounts []AccountsInfo `json:"accounts"`
}
type AccountsInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListAccountsArgs struct{}

func ListAccounts(args ListAccountsArgs) (*mcp_golang.ToolResponse, error) {
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

	res, err := json.Marshal(accountInfo.AuthInfo.Accounts)
	if err != nil {
		return nil, err
	}

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(res))), nil
}
