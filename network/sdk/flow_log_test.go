// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2021-01-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/marstr/randname"
)

func TestFlowLog(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	flowLogName := config.AppendRandomSuffix("flowlog")
	networkWatcherName := config.AppendRandomSuffix("networkwatcher")
	networkSecurityGroupName := config.AppendRandomSuffix("networksecuritygroup")
	storageAccountName := randname.Prefixed{Prefix: "storageaccount", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	err = CreateNetworkWatcher(ctx, networkWatcherName)
	if err != nil {
		t.Fatalf("failed to create network watcher: % +v", err)
	}

	networkSecurityGroupId, err := CreateNetworkSecurityGroup(ctx, networkSecurityGroupName)
	if err != nil {
		t.Fatalf("failed to create network security group: % +v", err)
	}

	storageAccountCreateParameters := armstorage.StorageAccountCreateParameters{
		Kind:     armstorage.KindStorage.ToPtr(),
		Location: to.StringPtr(config.Location()),
		SKU: &armstorage.SKU{
			Name: armstorage.SKUNameStandardLRS.ToPtr(),
		},
	}
	stroageAountId, err := CreateStorageAccount(ctx, storageAccountName, storageAccountCreateParameters)
	if err != nil {
		t.Fatalf("failed to create storage account: % +v", err)
	}

	flowLogParameters := armnetwork.FlowLog{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.FlowLogPropertiesFormat{
			Enabled: to.BoolPtr(true),
			Format: &armnetwork.FlowLogFormatParameters{
				Type:    armnetwork.FlowLogFormatTypeJSON.ToPtr(),
				Version: to.Int32Ptr(1),
			},
			StorageID:        &stroageAountId,
			TargetResourceID: &networkSecurityGroupId,
		},
	}
	err = CreateFlowLog(ctx, networkWatcherName, flowLogName, flowLogParameters)
	if err != nil {
		t.Fatalf("failed to create flow log: % +v", err)
	}
	t.Logf("created flow log")

	err = GetFlowLog(ctx, networkWatcherName, flowLogName)
	if err != nil {
		t.Fatalf("failed to get flow log: %+v", err)
	}
	t.Logf("got flow log")

	err = DeleteFlowLog(ctx, networkWatcherName, flowLogName)
	if err != nil {
		t.Fatalf("failed to delete flow log: %+v", err)
	}
	t.Logf("deleted flow log")
}
