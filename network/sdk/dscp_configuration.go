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

func getDscpConfigurationClient() armnetwork.DscpConfigurationClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewDscpConfigurationClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates a DSCP Configuration
func CreateDscpConfiguration(ctx context.Context, dscpConfigurationName string, dscpConfigurationParameters armnetwork.DscpConfiguration) error {
	client := getDscpConfigurationClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		dscpConfigurationName,
		dscpConfigurationParameters,
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

// Gets a DSCP Configuration.
func GetDscpConfiguration(ctx context.Context, dscpConfigurationName string) error {
	client := getDscpConfigurationClient()
	_, err := client.Get(ctx, config.GroupName(), dscpConfigurationName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all dscp configurations in a subscription.
func ListAllDscpConfiguration(ctx context.Context) error {
	client := getDscpConfigurationClient()
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

// Deletes a DSCP Configuration.
func DeleteDscpConfiguration(ctx context.Context, dscpConfigurationName string) error {
	client := getDscpConfigurationClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), dscpConfigurationName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
