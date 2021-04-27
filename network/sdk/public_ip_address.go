// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getIPAddressClient() armnetwork.PublicIPAddressesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPAddressesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create public IP address
func CreatePublicIPAddress(ctx context.Context, addressName string) error {
	ipClient := getIPAddressClient()
	poller, err := ipClient.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		addressName,
		armnetwork.PublicIPAddress{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(addressName),
				Location: to.StringPtr(config.Location()),
			},

			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
				PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	if resp.RawResponse != nil {
		log.Printf("create public ip address, name: %v", addressName)
	}
	return nil
}

// Gets the specified public IP address in a specified resource group.
func GetPublicIPAddress(ctx context.Context, ipName string) error {
	ipClient := getIPAddressClient()
	resp, err := ipClient.Get(ctx, config.GroupName(), ipName, nil)
	if err != nil {
		return err
	}
	log.Printf("get public ip address, name: %v", *resp.PublicIPAddress.Name)
	return nil
}

// Gets all the public IP prefixes in a subscription.
func ListPublicIPAddress(ctx context.Context) error {
	ipClient := getIPAddressClient()
	pager := ipClient.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
		var resp = pager.PageResponse().PublicIPAddressListResult
		var b strings.Builder
		if resp != nil && resp.Value != nil {
			for _, v := range *resp.Value {
				b.WriteString(*v.Properties.IPAddress)
				b.WriteString(",")
			}
			log.Printf("list public ip address in a resource group, IPAddress: %v\n", strings.TrimRight(b.String(), ","))
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Gets all the public IP addresses in a subscription.
func ListAllPublicIPAddress(ctx context.Context) error {
	ipClient := getIPAddressClient()
	pager := ipClient.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
		var resp = pager.PageResponse().PublicIPAddressListResult
		var b strings.Builder
		if resp != nil && resp.Value != nil {
			for _, v := range *resp.Value {
				b.WriteString(*v.Name)
				b.WriteString(",")
			}
			log.Printf("list all public ip address in a subscription, name: %v\n", strings.TrimRight(b.String(), ","))
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Updates public IP address tags.
func UpdateAddressTags(ctx context.Context, prefixName string) error {
	ipClient := getIPAddressClient()
	resp, err := ipClient.UpdateTags(
		ctx,
		config.GroupName(),
		prefixName,
		armnetwork.TagsObject{
			Tags: &map[string]string{"tag1": "value1", "tag2": "value2"},
		},
		nil,
	)
	if err != nil {
		return err
	}
	if resp.RawResponse != nil {
		log.Printf("update prefix tags, name: %v\n", *resp.PublicIPAddress.Name)
	}
	return nil
}

// Deletes the specified public IP address.
func DeletePublicIPAddress(ctx context.Context, addressName string) error {
	ipClient := getIPAddressClient()
	resp, err := ipClient.BeginDelete(ctx, config.GroupName(), addressName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	log.Printf("delete public address, name: %v\n", addressName)
	return nil
}
