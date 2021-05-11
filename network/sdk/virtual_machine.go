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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getVirtualMachinesClient() armcompute.VirtualMachinesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create VirtualMachines
func CreateVirtualMachine(ctx context.Context, virtualMachineName string, nicId *string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualMachineName,
		armcompute.VirtualMachine{
			Resource: armcompute.Resource{
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armcompute.VirtualMachineProperties{
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: armcompute.VirtualMachineSizeTypesStandardD2V2.ToPtr(),
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: &[]armcompute.NetworkInterfaceReference{
						{
							SubResource: armcompute.SubResource{
								ID: nicId,
							},
							Properties: &armcompute.NetworkInterfaceReferenceProperties{
								Primary: to.BoolPtr(true),
							},
						},
					},
				},
				OSProfile: &armcompute.OSProfile{
					AdminPassword: to.StringPtr("Aa1!zyx_"),
					AdminUsername: to.StringPtr("testuser"),
					ComputerName:  to.StringPtr("myVM"),
					WindowsConfiguration: &armcompute.WindowsConfiguration{
						EnableAutomaticUpdates: to.BoolPtr(true),
					},
				},
				StorageProfile: &armcompute.StorageProfile{
					DataDisks: &[]armcompute.DataDisk{
						{
							CreateOption: armcompute.DiskCreateOptionTypesEmpty.ToPtr(),
							DiskSizeGb:   to.Int32Ptr(1023),
							Lun:          to.Int32Ptr(0),
						},
						{
							CreateOption: armcompute.DiskCreateOptionTypesEmpty.ToPtr(),
							DiskSizeGb:   to.Int32Ptr(1023),
							Lun:          to.Int32Ptr(1),
						},
					},
					ImageReference: &armcompute.ImageReference{
						Offer:     to.StringPtr("WindowsServer"),
						Publisher: to.StringPtr("MicrosoftWindowsServer"),
						SKU:       to.StringPtr("2016-Datacenter"),
						Version:   to.StringPtr("latest"),
					},
					OSDisk: &armcompute.OSDisk{
						Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
						CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
						ManagedDisk: &armcompute.ManagedDiskParameters{
							StorageAccountType: armcompute.StorageAccountTypesStandardLrs.ToPtr(),
						},
						Name: to.StringPtr("myVMosdisk"),
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

// Deletes the specified virtual machine.
func DeleteVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
