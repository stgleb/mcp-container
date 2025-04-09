package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	ctx := context.Background()

	// Initialize the Stdio MCP Client with the server command
	client, err := client.NewStdioMCPClient("./bin/server", []string{}) // Adjust as necessary
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer client.Close()

	// Initialize the client
	initRequest := mcp.InitializeRequest{}
	initRequest.Method = string(mcp.MethodInitialize)

	_, err = client.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	listToolReq := mcp.ListToolsRequest{}
	listToolReq.Method = string(mcp.MethodToolsList)
	listToolResult, err := client.ListTools(ctx, listToolReq)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Println("Tools available")
	for _, tool := range listToolResult.Tools {
		fmt.Printf("%s \n%s\n", tool.Name, tool.Description)
	}

	// Prepare the request to run a container
	runContainerRequest := mcp.CallToolRequest{}
	runContainerRequest.Request.Method = string(mcp.MethodToolsCall)
	runContainerRequest.Params.Name = "run_container"
	runContainerRequest.Params.Arguments = map[string]interface{}{
		"image": "hello-world",
	}

	// Call the tool
	callToolResult, err := client.CallTool(ctx, runContainerRequest)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	fmt.Printf("Container output: \n\n%s", callToolResult.Content)
}
