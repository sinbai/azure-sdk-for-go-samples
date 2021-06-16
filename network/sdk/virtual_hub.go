// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/helper/resource"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getVirtualHubsClient() armnetwork.VirtualHubsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewVirtualHubsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

func virtualHubCreateRefreshFunc(ctx context.Context, client *armnetwork.VirtualHubsClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if res.RawResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("virtual Hub %q (Resource Group %q) doesn't exist", resourceGroup, name)
			}

			return nil, "", fmt.Errorf("retrieving Virtual Hub %q (Resource Group %q): %+v", resourceGroup, name, err)
		}
		if res.VirtualHub.Properties == nil {
			return nil, "", fmt.Errorf("unexpected nil properties of Virtual Hub %q (Resource Group %q)", resourceGroup, name)
		}

		state := *res.VirtualHub.Properties.RoutingState
		if state == "Failed" {
			return nil, "", fmt.Errorf("failed to provision routing on Virtual Hub %q (Resource Group %q)", resourceGroup, name)
		}
		return res, string(state), nil
	}
}

func virtualHubUpdateRefreshFunc(ctx context.Context, client *armnetwork.VirtualHubsClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if res.RawResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("virtual Hub %q (Resource Group %q) doesn't exist", resourceGroup, name)
			}

			return nil, "", fmt.Errorf("retrieving Virtual Hub %q (Resource Group %q): %+v", resourceGroup, name, err)
		}
		if res.VirtualHub.Properties == nil {
			return nil, "", fmt.Errorf("unexpected nil properties of Virtual Hub %q (Resource Group %q)", resourceGroup, name)
		}

		state := *res.VirtualHub.Properties.ProvisioningState
		if state == "Failed" {
			return nil, "", fmt.Errorf("failed to provision routing on Virtual Hub %q (Resource Group %q)", resourceGroup, name)
		}
		return res, string(state), nil
	}
}

// Create VirtualHubs
func CreateVirtualHub(ctx context.Context, virtualHubName string, virtualHubParameters armnetwork.VirtualHub) (string, error) {
	client := getVirtualHubsClient()

	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualHubName,
		virtualHubParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	// Hub returns provisioned while the routing state is still "provisining". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, &client, config.GroupName(), virtualHubName),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	respRaw, err := stateConf.WaitForState()
	if err != nil {
		return "", fmt.Errorf("waiting for Virtual Hub %q (Host Group Name %q) provisioning route: %+v", virtualHubName, config.GroupName(), err)
	}
	response := respRaw.(armnetwork.VirtualHubResponse)
	if response.VirtualHub.ID == nil {
		return "", fmt.Errorf("cannot read Virtual Hub %q (Resource Group %q) ID", virtualHubName, config.GroupName())
	}

	return *response.VirtualHub.ID, nil
}

// Retrieves the details of a VirtualHub
func GetVirtualHub(ctx context.Context, virtualHubName string) error {
	client := getVirtualHubsClient()
	_, err := client.Get(ctx, config.GroupName(), virtualHubName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all the VirtualHubs in a subscription.
func ListVirtualHub(ctx context.Context) error {
	client := getVirtualHubsClient()
	pager := client.List(nil)

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

//  Updates VirtualHub tags.
func UpdateVirtualHubTags(ctx context.Context, virtualHubName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getVirtualHubsClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		virtualHubName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}

	// Hub returns state is "updating". This might cause deletion to fail.
	// As a workaround, we will poll the hub state and ensure it is "Succeeded".
	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Updating"},
		Target:                    []string{"Succeeded", "Failed", "None"},
		Refresh:                   virtualHubUpdateRefreshFunc(ctx, &client, config.GroupName(), virtualHubName),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	respRaw, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for Virtual Hub %q (Host Group Name %q) update: %+v", virtualHubName, config.GroupName(), err)
	}
	response := respRaw.(armnetwork.VirtualHubResponse)
	if response.VirtualHub.ID == nil {
		return fmt.Errorf("cannot read Virtual Hub %q (Resource Group %q) ID", virtualHubName, config.GroupName())
	}
	return nil
}

// Deletes a VirtualHub.
func DeleteVirtualHub(ctx context.Context, virtualHubName string) error {
	client := getVirtualHubsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualHubName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Lists all the VirtualHubs in a resource group.
func ListVirtualHubByResourceGroup(ctx context.Context) error {
	client := getVirtualHubsClient()
	pager := client.ListByResourceGroup(config.GroupName(), nil)
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
