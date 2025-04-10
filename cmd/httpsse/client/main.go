package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// Create a new SSE MCP client connected to the server
	sseClient, err := client.NewSSEMCPClient("http://localhost:8080/sse")
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer sseClient.Close()
	ctx := context.Background()

	// Start the SSE client to listen for server events
	if err := sseClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start SSE client: %v", err)
	}
	fmt.Println("SSE client started.")

	// Initialize the client
	initRequest := mcp.InitializeRequest{}
	initRequest.Method = string(mcp.MethodInitialize)

	_, err = sseClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// List tools available on the server
	listToolReq := mcp.ListToolsRequest{}
	listToolReq.Method = string(mcp.MethodToolsList)
	toolsResult, err := sseClient.ListTools(ctx, listToolReq)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	// Print the available tools
	fmt.Println("Available tools:")
	for _, tool := range toolsResult.Tools {
		fmt.Printf(" - %s: %s\n", tool.Name, tool.Description)
	}

	// Prepare the request to run a container
	runContainerRequest := mcp.CallToolRequest{}
	runContainerRequest.Request.Method = string(mcp.MethodToolsCall)
	runContainerRequest.Params.Name = "run_container"
	runContainerRequest.Params.Arguments = map[string]interface{}{
		"image": "hello-world",
	}

	// Call the selected tool
	response, err := sseClient.CallTool(ctx, runContainerRequest)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	fmt.Println(response.Content)
}
