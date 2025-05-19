package action

import (
	"encoding/json"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

// Interface that we are exposing as a plugin
type ActionRPC struct {
	client *rpc.Client
}

// Implement all action interface functions
func (a *ActionRPC) Execute(ctx *security_context.SecurityContext) error {
	b, err := json.Marshal(ctx)
	if err != nil {
		return fmt.Errorf("client: failed to marshal context: %w", err)
	}

	var resp error
	err = a.client.Call("Plugin.Execute", string(b), &resp)	
	if err != nil {
		return fmt.Errorf("failed to execute plugin: %v", err)
	}

	return resp
}

func (a *ActionRPC) GetType() string {
	var resp string
	err := a.client.Call("Plugin.GetType", new(interface{}), &resp)
	if err != nil {
		panic("Failed to get type of plugin")
	}

	return resp
}

func (a *ActionRPC) GetName() string {
	var resp string
	err := a.client.Call("Plugin.GetName", new(interface{}), &resp) 
	if err != nil {
		panic("Failed to get name of plugin")
	}

	return resp
}

func (a *ActionRPC) Validate() error {
	var resp error
	err := a.client.Call("Plugin.Validate", new(interface{}), &resp)
	if err != nil {
		return fmt.Errorf("failed to validate the plugin: %v", err)
	}

	return resp
}

// RPC server that the ActionRPC talks to
type ActionRPCServer struct {
	Impl Action
}

func (s *ActionRPCServer) Execute(arg string, resp *error) error {
	var ctx security_context.SecurityContext
	if err := json.Unmarshal([]byte(arg), &ctx); err != nil {
		return fmt.Errorf("plugin: failed to unmarshal context: %w", err)
	}
	*resp = s.Impl.Execute(&ctx)

	return nil
}

func (s *ActionRPCServer) GetType(args any, resp *string) error {
	*resp = s.Impl.GetType()
	return nil
}

func (s ActionRPCServer) GetName(args any, resp *string) error {
	*resp = s.Impl.GetName()
	return nil
}

func (s *ActionRPCServer)  Validate(args any, resp *error) error {
	*resp = s.Impl.Validate()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a ActionRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return ActionRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type ActionPlugin struct {
	Impl Action
}

func (p *ActionPlugin) Server(*plugin.MuxBroker) (any, error) {
	return &ActionRPCServer{Impl: p.Impl}, nil
}

func (ActionPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &ActionRPC{client: c}, nil
}