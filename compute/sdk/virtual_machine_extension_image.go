// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getVirtualMachineExtensionImagesClient() armcompute.VirtualMachineExtensionImagesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineExtensionImagesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

//  Gets a virtual machine extension image.
func GetVirtualMachineExtensionImage(ctx context.Context, publisherName string, typeParameter string, version string) error {
	client := getVirtualMachineExtensionImagesClient()
	_, err := client.Get(ctx, config.Location(), publisherName, typeParameter, version, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets a list of virtual machine extension image types.
func ListVirtualMachineExtensionImageType(ctx context.Context, publisherName string) error {
	client := getVirtualMachineExtensionImagesClient()
	_, err := client.ListTypes(ctx, config.Location(), publisherName, nil)

	if err != nil {
		return err
	}
	return nil
}

// Gets a list of virtual machine extension image versions.
func ListVirtualMachineExtensionImageVersion(ctx context.Context, publisherName string, typeParameter string) error {
	client := getVirtualMachineExtensionImagesClient()
	_, err := client.ListVersions(ctx, config.Location(), publisherName, typeParameter, nil)
	if err != nil {
		return err
	}
	return nil
}
