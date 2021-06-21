// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
)

func getVirtualMachineImagesClient() armcompute.VirtualMachineImagesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineImagesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets a virtual machine image.
func GetVirtualMachineImage(ctx context.Context, publisherName string, offer string, skus string, version string) error {
	client := getVirtualMachineImagesClient()
	_, err := client.Get(ctx, config.Location(), publisherName, offer, skus, version, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets a list of all virtual machine image versions for the specified location, publisher, offer, and SKU.
func ListVirtualMachineImage(ctx context.Context, publisherName string, offer string, skus string) error {
	client := getVirtualMachineImagesClient()
	_, err := client.List(ctx, config.Location(), publisherName, offer, skus, nil)

	if err != nil {
		return err
	}
	return nil
}

// Gets a list of virtual machine image offers for the specified location and publisher.
func ListVirtualMachineImageOffer(ctx context.Context, publisherName string) error {
	client := getVirtualMachineImagesClient()
	_, err := client.ListOffers(ctx, config.Location(), publisherName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets a list of virtual machine image publishers for the specified Azure location.
func LisVirtualMachineImagePublisher(ctx context.Context) error {
	client := getVirtualMachineImagesClient()
	_, err := client.ListPublishers(ctx, config.Location(), nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets a list of virtual machine image SKUs for the specified location, publisher, and offer.
func ListVirtualMachineImageSKU(ctx context.Context, publisherName string, offer string) error {
	client := getVirtualMachineImagesClient()
	_, err := client.ListSKUs(ctx, config.Location(), publisherName, offer, nil)
	if err != nil {
		return err
	}
	return nil
}
