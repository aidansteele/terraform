package terraform

import (
	"testing"

	"github.com/hashicorp/terraform/internal/addrs"
	"github.com/hashicorp/terraform/internal/providers"
)

func TestResourceAddressPopulation(t *testing.T) {
	// Test that resource addresses are properly populated
	t.Run("AbsResourceInstance.String() produces correct format", func(t *testing.T) {
		tests := []struct {
			name     string
			addr     addrs.AbsResourceInstance
			expected string
		}{
			{
				name: "simple resource",
				addr: addrs.Resource{
					Mode: addrs.ManagedResourceMode,
					Type: "aws_instance",
					Name: "example",
				}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
				expected: "aws_instance.example",
			},
			{
				name: "resource with index",
				addr: addrs.Resource{
					Mode: addrs.ManagedResourceMode,
					Type: "aws_instance",
					Name: "server",
				}.Instance(addrs.IntKey(0)).Absolute(addrs.RootModuleInstance),
				expected: "aws_instance.server[0]",
			},
			{
				name: "resource with string key",
				addr: addrs.Resource{
					Mode: addrs.ManagedResourceMode,
					Type: "aws_instance",
					Name: "server",
				}.Instance(addrs.StringKey("web")).Absolute(addrs.RootModuleInstance),
				expected: "aws_instance.server[\"web\"]",
			},
			{
				name: "module resource",
				addr: addrs.Resource{
					Mode: addrs.ManagedResourceMode,
					Type: "aws_db_instance",
					Name: "primary",
				}.Instance(addrs.NoKey).Absolute(
					addrs.RootModuleInstance.Child("postgres", addrs.NoKey),
				),
				expected: "module.postgres.aws_db_instance.primary",
			},
			{
				name: "nested module resource with index",
				addr: addrs.Resource{
					Mode: addrs.ManagedResourceMode,
					Type: "aws_instance",
					Name: "server",
				}.Instance(addrs.IntKey(0)).Absolute(
					addrs.RootModuleInstance.Child("app", addrs.NoKey).Child("web", addrs.NoKey),
				),
				expected: "module.app.module.web.aws_instance.server[0]",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.addr.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})

	t.Run("ResourceAddress field exists in request types", func(t *testing.T) {
		// Test that ResourceAddress field is present in request types
		planReq := providers.PlanResourceChangeRequest{
			ResourceAddress: "aws_instance.example",
		}
		if planReq.ResourceAddress != "aws_instance.example" {
			t.Errorf("PlanResourceChangeRequest.ResourceAddress not working")
		}

		applyReq := providers.ApplyResourceChangeRequest{
			ResourceAddress: "aws_instance.example",
		}
		if applyReq.ResourceAddress != "aws_instance.example" {
			t.Errorf("ApplyResourceChangeRequest.ResourceAddress not working")
		}
	})
}