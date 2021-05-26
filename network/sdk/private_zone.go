// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/iam"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
)

func getPrivateZonesClient() privatedns.PrivateZonesClient {
	privateZonesClient := privatedns.NewPrivateZonesClient(config.SubscriptionID())
	a, _ := iam.GetResourceManagementAuthorizer()
	privateZonesClient.Authorizer = a
	privateZonesClient.AddToUserAgent(config.UserAgent())
	return privateZonesClient
}

// creates or updates a Private DNS zone. Does not modify Links to virtual networks or DNS records
func CreatePrivateZone(ctx context.Context, privateZoneName string, privateZoneParameters privatedns.PrivateZone) (privateZone privatedns.PrivateZone, err error) {
	client := getPrivateZonesClient()
	future, err := client.CreateOrUpdate(
		ctx,
		config.GroupName(),
		privateZoneName,
		privateZoneParameters,
		"",
		"",
	)

	if err != nil {
		return privateZone, err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return privateZone, err
	}

	return future.Result(client)
}
