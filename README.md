# Akamai Functions MCP Server

An MCP (Model Context Protocol) server that connects your AI Assistant or AI Agent to Akamai Functions, allowing you to retrieve essential information about the accounts you've access to and the applications you've deployed to Akamai Functions.

This MCP server supports `stdio`

## Prerequisites

The *Akamai Functions MCP Server* requires the following:

- The `spin` CLI being installed on your system (See [https://spinframework.dev](https://spinframework.dev))
- The `aka` plugin for `spin` must be installed and authorized for interacting with Akamai Functions on your behalf

## Available Tools 🛠️

The *Akamai Functions MCP Server* provides the following tools ⚒️:

### Account-specific Tools

- `list_accounts` - List all Akamai Functions accounts you can access
- `list_apps` - List all apps for a particular Akamai Functions account (not providing an `account id`, will list tools from your default account)
- `search_app` - Find an application by a search term (`query`). This command iterates over all Akamai Functions accounts you've access to

### App-specific Tools

- `get_app_status` - Retrieve the status of an application deployed to Akamai Functions
- `get_app_url` - Retrieve the public URL of an application deployed to Akamai Functions
- `get_app_history` - Retrieve deployment history of an application deployed to Akamai Functions
- `get_app_logs` - Retrieve logs of an application deployed to Akamai Functions

## Available Resources 📄

- The `spin aka` command Reference - Returns the full command reference as markdown

## Available Resource Templates

- Help for `spin aka` sub-commands - Provides command documentation (provide a sub-command (supports nested sub-commands))
