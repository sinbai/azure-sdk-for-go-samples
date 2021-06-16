// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package storage

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/storage/2021-01-01/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getStorageAccountsClient() armstorage.StorageAccountsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create StorageAccounts
func CreateStorageAccount(ctx context.Context, storageAccountName string, storageAccountCreateParametersParameters armstorage.StorageAccountCreateParameters) (string, error) {
	client := getStorageAccountsClient()
	poller, err := client.BeginCreate(
		ctx,
		config.GroupName(),
		storageAccountName,
		storageAccountCreateParametersParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	if resp.StorageAccount.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.StorageAccount.ID, nil
}
