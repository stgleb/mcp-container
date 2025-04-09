package main

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/testcontainers/testcontainers-go"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Container Runner",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add a container runner tool
	containerTool := mcp.NewTool("run_container",
		mcp.WithDescription("Run a Docker container and return its output"),
		mcp.WithString("image",
			mcp.Required(),
			mcp.Description("The Docker image to run"),
		),
	)

	// Add the container runner handler
	s.AddTool(containerTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		image := request.Params.Arguments["image"].(string)

		// Initialize a testcontainers context
		ctx = context.Background()

		// Define a container request
		req := testcontainers.ContainerRequest{
			Image: image,
		}

		// Start the container
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			return nil, errors.New("could not start container: " + err.Error())
		}
		defer container.Terminate(ctx)

		// Retrieve STDOUT of the container
		logReader, err := container.Logs(ctx)
		if err != nil {
			return nil, errors.New("could not get logs: " + err.Error())
		}
		logBytes, err := io.ReadAll(logReader)
		if err != nil {
			return nil, errors.New("could not read logs: " + err.Error())
		}
		return mcp.NewToolResultText(string(logBytes)), nil
	})

	// Start the server
	log.Printf("Starting stdio server")
	/*
		Paste to stdio:

		{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}

		To call tool

		{"jsonrpc":"2.0","method":"tools/call","params":{"name":"run_container","arguments":{"image":"hello-world"}},"id":2}

	*/
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
