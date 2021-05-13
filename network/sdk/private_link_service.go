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
	"github.com/Azure/go-autorest/autorest/to"
)

func getPrivateLinkServicesClient() armnetwork.PrivateLinkServicesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPrivateLinkServicesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create PrivateLinkServices
func CreatePrivateLinkService(ctx context.Context, privateLinkServiceName string, virtualNetworkName string, loadBalancerName string, ipConfigurationName string, subNetName string) error {
	client := getPrivateLinkServicesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		privateLinkServiceName,
		armnetwork.PrivateLinkService{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.PrivateLinkServiceProperties{
				AutoApproval: &armnetwork.PrivateLinkServicePropertiesAutoApproval{
					ResourceSet: armnetwork.ResourceSet{
						Subscriptions: &[]*string{to.StringPtr(config.SubscriptionID())},
					},
				},
				Fqdns: &[]*string{to.StringPtr("fqdn1"),
					to.StringPtr("fqdn2"),
					to.StringPtr("fqdn3")},
				IPConfigurations: &[]*armnetwork.PrivateLinkServiceIPConfiguration{{
					Name: &ipConfigurationName,
					Properties: &armnetwork.PrivateLinkServiceIPConfigurationProperties{
						PrivateIPAddress:          to.StringPtr("10.0.1.5"),
						PrivateIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
						PrivateIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
						Subnet: &armnetwork.Subnet{
							SubResource: armnetwork.SubResource{
								ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/virtualNetworks/" + virtualNetworkName + "/subnets/" + subNetName),
							},
						},
					},
				}},
				LoadBalancerFrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{{
					SubResource: armnetwork.SubResource{
						ID: to.StringPtr("/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName + "/frontendIPConfigurations/" + ipConfigurationName),
					},
				}},
				Visibility: &armnetwork.PrivateLinkServicePropertiesVisibility{
					ResourceSet: armnetwork.ResourceSet{
						Subscriptions: &[]*string{to.StringPtr(config.SubscriptionID())},
					},
				},
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

// Updates private endpoint connection
func UpdatePrivateEndpointConnection(ctx context.Context, privateLinkServiceName string, peConnectionName string) error {
	client := getPrivateLinkServicesClient()
	_, err := client.UpdatePrivateEndpointConnection(
		ctx,
		config.GroupName(),
		privateLinkServiceName,
		peConnectionName,
		armnetwork.PrivateEndpointConnection{
			Name: &peConnectionName,
			Properties: &armnetwork.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: &armnetwork.PrivateLinkServiceConnectionState{
					Description: to.StringPtr("approved it for some reason."),
					Status:      to.StringPtr("Approved"),
				},
			},
		},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Get the specific private end point connection by specific private link service in the resource group.
func GetPrivateEndpointConnection(ctx context.Context, privateLinkServiceName string, peConnectionName string) error {
	client := getPrivateLinkServicesClient()
	_, err := client.GetPrivateEndpointConnection(ctx, config.GroupName(), privateLinkServiceName, peConnectionName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all private end point connections for a specific private link service.
func ListPrivateEndpointConnection(ctx context.Context, privateLinkServiceName string) error {
	client := getPrivateLinkServicesClient()
	pager := client.ListPrivateEndpointConnections(config.GroupName(), privateLinkServiceName, nil)

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

// Returns all of the private link service ids that can be linked to a Private Endpoint with auto approved
func ListAutoApprovedPrivateLinkServicesByResourceGroup(ctx context.Context) error {
	client := getPrivateLinkServicesClient()
	pager := client.ListAutoApprovedPrivateLinkServicesByResourceGroup(config.Location(), config.GroupName(), nil)

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

// Gets the specified private link service in a specified resource group.
func GetPrivateLinkService(ctx context.Context, privateLinkServiceName string) (string, error) {
	client := getPrivateLinkServicesClient()
	result, err := client.Get(ctx, config.GroupName(), privateLinkServiceName, nil)
	if err != nil {
		return "", err
	}
	if len((*(result.PrivateLinkService.Properties.PrivateEndpointConnections))) > 0 {
		return *((*(result.PrivateLinkService.Properties.PrivateEndpointConnections))[0].Name), nil
	}

	return "", nil
}

// Gets all private end point connections for a specific private link service.
func ListPrivateEndpointConnections(ctx context.Context, privateLinkServiceName string) error {
	client := getPrivateLinkServicesClient()
	pager := client.ListPrivateEndpointConnections(config.GroupName(), privateLinkServiceName, nil)

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

// Gets all private link service in a subscription.
func ListBySubscription(ctx context.Context) error {
	client := getPrivateLinkServicesClient()
	pager := client.ListBySubscription(nil)

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

// Deletes the specified private link service.
func BeginDeletePrivateEndpointConnection(ctx context.Context, privateLinkServiceName string, peConnectionName string) error {
	client := getPrivateLinkServicesClient()
	resp, err := client.BeginDeletePrivateEndpointConnection(ctx, config.GroupName(), privateLinkServiceName, peConnectionName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified private link service.
func DeletePrivateLinkService(ctx context.Context, privateLinkServiceName string) error {
	client := getPrivateLinkServicesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), privateLinkServiceName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
