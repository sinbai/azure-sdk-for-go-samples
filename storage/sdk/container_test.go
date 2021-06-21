// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/marstr/randname"
)

func TestStorageContainer(t *testing.T) {
	groupName := config.GenerateGroupName("storage")
	config.SetGroupName(groupName)

	storageAccountName := randname.Prefixed{Prefix: "storageaccount", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()
	containerName := randname.Prefixed{Prefix: "blobcontainer", Acceptable: randname.LowercaseAlphabet, Len: 5}.Generate()

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	storageAccountCreateParameters := armstorage.StorageAccountCreateParameters{
		Kind:     armstorage.KindStorage.ToPtr(),
		Location: to.StringPtr(config.Location()),
		SKU: &armstorage.SKU{
			Name: armstorage.SKUNameStandardLRS.ToPtr(),
		},
	}
	_, err = CreateStorageAccount(ctx, storageAccountName, storageAccountCreateParameters)
	if err != nil {
		t.Fatalf("failed to create storage account: % +v", err)
	}

	blobContainerParameters := armstorage.BlobContainer{}
	_, err = CreateBlobContainer(ctx, storageAccountName, containerName, blobContainerParameters)
	if err != nil {
		t.Fatalf("failed to create blob container: % +v", err)
	}
	t.Logf("created blob container")
}
