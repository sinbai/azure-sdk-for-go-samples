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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getP2sVpnGatewaysClient() armnetwork.P2SVPNGatewaysClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewP2SVPNGatewaysClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates a virtual wan p2s vpn gateway if it doesn't exist else updates the existing gateway.
func CreateP2sVpnGateway(ctx context.Context, gatewayName string, p2SVPNGatewayParameters armnetwork.P2SVPNGateway) error {
	client := getP2sVpnGatewaysClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		gatewayName,
		p2SVPNGatewayParameters,
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
