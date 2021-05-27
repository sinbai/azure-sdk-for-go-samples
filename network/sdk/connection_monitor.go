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

func getConnectionMonitorsClient() armnetwork.ConnectionMonitorsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewConnectionMonitorsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create ConnectionMonitors
func CreateConnectionMonitor(ctx context.Context, networkWatcherName string, connectionMonitorName string, connectionMonitorParameters armnetwork.ConnectionMonitor) error {
	client := getConnectionMonitorsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		networkWatcherName,
		connectionMonitorName,
		connectionMonitorParameters,
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

//  Gets a connection monitor by name.
func GetConnectionMonitor(ctx context.Context, networkWatcherName string, connectionMonitorName string) error {
	client := getConnectionMonitorsClient()
	_, err := client.Get(ctx, config.GroupName(), networkWatcherName, connectionMonitorName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all connection monitors for the specified Network Watcher.
func ListConnectionMonitor(ctx context.Context, networkWatcherName string) error {
	client := getConnectionMonitorsClient()
	_, err := client.List(ctx, config.GroupName(), networkWatcherName, nil)

	if err != nil {
		return err
	}
	return nil
}

// Update tags of the specified connection monitor.
func UpdateConnectionMonitorTags(ctx context.Context, networkWatcherName string, connectionMonitorName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getConnectionMonitorsClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkWatcherName,
		connectionMonitorName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified connection monitor.
func DeleteConnectionMonitor(ctx context.Context, networkWatcherName string, connectionMonitorName string) error {
	client := getConnectionMonitorsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkWatcherName, connectionMonitorName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
