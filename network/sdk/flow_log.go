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
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
)

func getFlowLogsClient() armnetwork.FlowLogsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewFlowLogsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create FlowLogs
func CreateFlowLog(ctx context.Context, networkWatcherName string, flowLogName string, flowLogParameters armnetwork.FlowLog) error {
	client := getFlowLogsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkWatcherName,
		flowLogName,
		flowLogParameters,
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

// Gets the specified flow log in a specified resource group.
func GetFlowLog(ctx context.Context, networkWatcherName string, flowLogName string) error {
	client := getFlowLogsClient()
	_, err := client.Get(ctx, config.GroupName(), networkWatcherName, flowLogName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified flow log.
func DeleteFlowLog(ctx context.Context, networkWatcherName string, flowLogName string) error {
	client := getFlowLogsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkWatcherName, flowLogName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
