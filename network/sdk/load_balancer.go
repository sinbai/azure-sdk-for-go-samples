// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"net/url"
	"strings"
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
func CreateLoadBalancer(ctx context.Context, loadBalancerName string, publicIpAddressName string,
	frontendIpConfigurationName string, backendAddressPoolName string, probeName string,
	loadBalancingRuleName string, outBoundRuleName string) error {
	urlPathAddress := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/{publicIpAddressName}"
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{publicIpAddressName}", url.PathEscape(publicIpAddressName))
	urlPathAddress = strings.ReplaceAll(urlPathAddress, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))

	urlPathFrontendIPConfiguration := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/frontendIPConfigurations/{frontendIpConfigurationName}"
	urlPathFrontendIPConfiguration = strings.ReplaceAll(urlPathFrontendIPConfiguration, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathFrontendIPConfiguration = strings.ReplaceAll(urlPathFrontendIPConfiguration, "{loadBalancerName}", url.PathEscape(loadBalancerName))
	urlPathFrontendIPConfiguration = strings.ReplaceAll(urlPathFrontendIPConfiguration, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))
	urlPathFrontendIPConfiguration = strings.ReplaceAll(urlPathFrontendIPConfiguration, "{frontendIpConfigurationName}", url.PathEscape(frontendIpConfigurationName))

	urlPathBackendAddressPool := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/backendAddressPools/{backendAddressPoolName}"
	urlPathBackendAddressPool = strings.ReplaceAll(urlPathBackendAddressPool, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathBackendAddressPool = strings.ReplaceAll(urlPathBackendAddressPool, "{loadBalancerName}", url.PathEscape(loadBalancerName))
	urlPathBackendAddressPool = strings.ReplaceAll(urlPathBackendAddressPool, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))
	urlPathBackendAddressPool = strings.ReplaceAll(urlPathBackendAddressPool, "{backendAddressPoolName}", url.PathEscape(backendAddressPoolName))

	urlPathProb := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/probes/{probeName}"
	urlPathProb = strings.ReplaceAll(urlPathProb, "{resourceGroupName}", url.PathEscape(config.GroupName()))
	urlPathProb = strings.ReplaceAll(urlPathProb, "{loadBalancerName}", url.PathEscape(loadBalancerName))
	urlPathProb = strings.ReplaceAll(urlPathProb, "{subscriptionId}", url.PathEscape(config.SubscriptionID()))
	urlPathProb = strings.ReplaceAll(urlPathProb, "{probeName}", url.PathEscape(probeName))

	client := getLoadBalancersClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		loadBalancerName,
		armnetwork.LoadBalancer{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				BackendAddressPools: &[]*armnetwork.BackendAddressPool{
					{
						Name: &backendAddressPoolName,
					},
				},
				FrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{
					{
						Name: &frontendIpConfigurationName,
						Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
							PublicIPAddress: &armnetwork.PublicIPAddress{
								Resource: armnetwork.Resource{
									ID: &urlPathAddress,
								},
							},
						},
					},
				},
				LoadBalancingRules: &[]*armnetwork.LoadBalancingRule{
					{
						Name: &loadBalancingRuleName,
						Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
							BackendAddressPool: &armnetwork.SubResource{
								ID: &urlPathBackendAddressPool,
							},
							BackendPort:         to.Int32Ptr(80),
							DisableOutboundSnat: to.BoolPtr(true),
							EnableFloatingIP:    to.BoolPtr(true),
							EnableTCPReset:      new(bool),
							FrontendIPConfiguration: &armnetwork.SubResource{
								ID: &urlPathFrontendIPConfiguration,
							},
							FrontendPort:         to.Int32Ptr(80),
							IdleTimeoutInMinutes: to.Int32Ptr(15),
							LoadDistribution:     armnetwork.LoadDistributionDefault.ToPtr(),
							Probe: &armnetwork.SubResource{
								ID: &urlPathProb,
							},
							Protocol: armnetwork.TransportProtocolTCP.ToPtr(),
						},
					},
				},
				OutboundRules: &[]*armnetwork.OutboundRule{
					{
						Name: &outBoundRuleName,
						Properties: &armnetwork.OutboundRulePropertiesFormat{
							BackendAddressPool: &armnetwork.SubResource{
								ID: &urlPathBackendAddressPool,
							},
							FrontendIPConfigurations: &[]*armnetwork.SubResource{
								{
									ID: &urlPathFrontendIPConfiguration,
								},
							},
							Protocol: armnetwork.LoadBalancerOutboundRuleProtocolAll.ToPtr(),
						},
					},
				},
				Probes: &[]*armnetwork.Probe{
					{
						Name: &probeName,
						Properties: &armnetwork.ProbePropertiesFormat{
							IntervalInSeconds: to.Int32Ptr(15),
							NumberOfProbes:    to.Int32Ptr(2),
							Port:              to.Int32Ptr(80),
							Protocol:          armnetwork.ProbeProtocolHTTP.ToPtr(),
							RequestPath:       to.StringPtr("healthcheck.aspx"),
						},
					},
				},
			},
			SKU: &armnetwork.LoadBalancerSKU{
				Name: armnetwork.LoadBalancerSKUNameStandard.ToPtr(),
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
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
