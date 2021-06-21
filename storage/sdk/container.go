// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package storage

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/armstorage"
)

func getBlobContainersClient() armstorage.BlobContainersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a new container under the specified account as described by request body. The container resource includes metadata and properties for
// that container. It does not include a list of the blobs contained by the container.
func CreateBlobContainer(ctx context.Context, accountName string, containerName string, blobContainerParameters armstorage.BlobContainer) (string, error) {
	client := getBlobContainersClient()
	resp, err := client.Create(
		ctx,
		config.GroupName(),
		accountName,
		containerName,
		blobContainerParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	if resp.BlobContainer.ID == nil {
		return resp.RawResponse.Request.URL.Path, nil
	}
	return *resp.BlobContainer.ID, nil
}
