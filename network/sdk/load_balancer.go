// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getLoadBalancersClient() armnetwork.LoadBalancersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancersClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create LoadBalancers
func CreateLoadBalancer(ctx context.Context, loadBalancerName string, loadBalancerParameters armnetwork.LoadBalancer) (string, error) {
	client := getLoadBalancersClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		loadBalancerName,
		loadBalancerParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.LoadBalancer.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.LoadBalancer.ID, nil
}

// Gets the specified load balancer in a specified resource group.
func GetLoadBalancer(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancersClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the load balancers in a resource group
func ListLoadBalancer(ctx context.Context) error {
	client := getLoadBalancersClient()
	pager := client.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Gets all the load balancer in a subscription.
func ListAllLoadBalancer(ctx context.Context) error {
	client := getLoadBalancersClient()
	pager := client.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Updates a load balancer tags.
func UpdateLoadBalancerTags(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancersClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		loadBalancerName,
		armnetwork.TagsObject{
			Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
		},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified load balancer.
func DeleteLoadBalancer(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancersClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), loadBalancerName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func getLoadBalancerFrontendIPConfigurationsClient() armnetwork.LoadBalancerFrontendIPConfigurationsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancerFrontendIPConfigurationsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets load balancer frontend IP configuration.
func GetLoadBalancerFrontendIPConfiguration(ctx context.Context, loadBalancerName string, loadBalancerFrontendIPConfigurationName string) error {
	client := getLoadBalancerFrontendIPConfigurationsClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerName, loadBalancerFrontendIPConfigurationName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the load balancer frontend IP configurations.
func ListLoadBalancerFrontendIPConfiguration(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancerFrontendIPConfigurationsClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getLoadBalancerBackendAddressPoolsClient() armnetwork.LoadBalancerBackendAddressPoolsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancerBackendAddressPoolsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets load balancer backend address pool.
func GetLoadBalancerBackendAddressPool(ctx context.Context, loadBalancerBackendAddressPoolName string, backendAddressPoolName string) error {
	client := getLoadBalancerBackendAddressPoolsClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerBackendAddressPoolName, backendAddressPoolName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the load balancer backed address pools.
func ListLoadBalancerBackendAddressPool(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancerBackendAddressPoolsClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getLoadBalancerLoadBalancingRulesClient() armnetwork.LoadBalancerLoadBalancingRulesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancerLoadBalancingRulesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets the specified load balancer load balancing rule.
func GetLoadBalancerLoadBalancingRule(ctx context.Context, loadBalancerName string, loadBalancingRuleName string) error {
	client := getLoadBalancerLoadBalancingRulesClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerName, loadBalancingRuleName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the load balancing rules in a load balancer.
func ListLoadBalancerLoadBalancingRule(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancerLoadBalancingRulesClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getLoadBalancerOutboundRulesClient() armnetwork.LoadBalancerOutboundRulesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancerOutboundRulesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets the specified load balancer outbound rule.
func GetLoadBalancerOutboundRule(ctx context.Context, loadBalancerName string, outBoundRuleName string) error {
	client := getLoadBalancerOutboundRulesClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerName, outBoundRuleName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the outbound rules in a load balancer.
func ListLoadBalancerOutboundRule(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancerOutboundRulesClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getLoadBalancerProbesClient() armnetwork.LoadBalancerProbesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancerProbesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets load balancer probe.
func GetLoadBalancerProbe(ctx context.Context, loadBalancerName string, probeName string) error {
	client := getLoadBalancerProbesClient()
	_, err := client.Get(ctx, config.GroupName(), loadBalancerName, probeName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the load balancer probes.
func ListLoadBalancerProbe(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancerProbesClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

func getLoadBalancerNetworkInterfacesClient() armnetwork.LoadBalancerNetworkInterfacesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancerNetworkInterfacesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets associated load balancer network interfaces.
func ListLoadBalancerNetworkInterface(ctx context.Context, loadBalancerName string) error {
	client := getLoadBalancerNetworkInterfacesClient()
	pager := client.List(config.GroupName(), loadBalancerName, nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}
