// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
)

func getServiceTagsClient() armnetwork.ServiceTagsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewServiceTagsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

//Gets a list of service tag information in a subscription.
func ListServiceTags(ctx context.Context) error {
	client := getServiceTagsClient()
	_, err := client.List(ctx, config.Location(), nil)

	if err != nil {
		return err
	}
	return nil
}
