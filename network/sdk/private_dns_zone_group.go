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
)

func getPrivateDnsZoneGroupsClient() armnetwork.PrivateDNSZoneGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPrivateDNSZoneGroupsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a private dns zone group in the specified private endpoint.
func CreatePrivateDnsZoneGroup(ctx context.Context, privateEndpointName string, privateDnsZoneGroupName string, privateDNSZoneGroupParameters armnetwork.PrivateDNSZoneGroup) error {
	client := getPrivateDnsZoneGroupsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		privateEndpointName,
		privateDnsZoneGroupName,
		privateDNSZoneGroupParameters,
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

// Gets the private dns zone group resource by specified private dns zone group name.
func GetPrivateDnsZoneGroup(ctx context.Context, privateEndpointName string, privateDnsZoneGroupName string) error {
	client := getPrivateDnsZoneGroupsClient()
	_, err := client.Get(ctx, config.GroupName(), privateEndpointName, privateDnsZoneGroupName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all private dns zone groups in a private endpoint.
func ListPrivateDnsZoneGroup(ctx context.Context, privateEndpointName string) error {
	client := getPrivateDnsZoneGroupsClient()
	pager := client.List(privateEndpointName, config.GroupName(), nil)

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

// Deletes the specified private dns zone group.
func DeletePrivateDnsZoneGroup(ctx context.Context, privateEndpointName string, privateDnsZoneGroupName string) error {
	client := getPrivateDnsZoneGroupsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), privateEndpointName, privateDnsZoneGroupName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
