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

func getPacketCapturesClient() armnetwork.PacketCapturesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPacketCapturesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create PacketCaptures
func CreatePacketCapture(ctx context.Context, networkWatcherName string, packetCaptureName string, packetCaptureParameters armnetwork.PacketCapture) error {
	client := getPacketCapturesClient()
	poller, err := client.BeginCreate(
		ctx,
		config.GroupName(),
		networkWatcherName,
		packetCaptureName,
		packetCaptureParameters,
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

// Gets the specified packet capture in a specified resource group.
func GetPacketCapture(ctx context.Context, networkWatcherName string, packetCaptureName string) error {
	client := getPacketCapturesClient()
	_, err := client.Get(ctx, config.GroupName(), networkWatcherName, packetCaptureName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the packet capture in a subscription.
func ListPacketCapture(ctx context.Context, networkWatcherName string) error {
	client := getPacketCapturesClient()
	_, err := client.List(ctx, config.GroupName(), networkWatcherName, nil)

	if err != nil {
		return err
	}
	return nil
}

// Query the status of a running packet capture session.
func GetPacketCaptureStatus(ctx context.Context, networkWatcherName string, packetCaptureName string) error {
	client := getPacketCapturesClient()
	resp, err := client.BeginGetStatus(ctx, config.GroupName(), networkWatcherName, packetCaptureName, nil)

	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Stops a specified packet capture session.
func StopPacketCapture(ctx context.Context, networkWatcherName string, packetCaptureName string) error {
	client := getPacketCapturesClient()
	resp, err := client.BeginStop(ctx, config.GroupName(), networkWatcherName, packetCaptureName, nil)

	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified packet capture.
func DeletePacketCapture(ctx context.Context, networkWatcherName string, packetCaptureName string) error {
	client := getPacketCapturesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkWatcherName, packetCaptureName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
