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

func getIPPrefixClient() armnetwork.PublicIPPrefixesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPPrefixesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create public IP prefix
func CreatePublicIPPrefix(ctx context.Context, prefixName string) {
	ipClient := getIPPrefixClient()

	poller, err := ipClient.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		prefixName,
		armnetwork.PublicIPPrefix{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(prefixName),
				Location: to.StringPtr(config.Location()),
			},
			Properties: &armnetwork.PublicIPPrefixPropertiesFormat{
				PrefixLength:           to.Int32Ptr(30),
				PublicIPAddressVersion: armnetwork.IPVersionIPv4.ToPtr(),
			},
			SKU: &armnetwork.PublicIPPrefixSKU{
				Name: armnetwork.PublicIPPrefixSKUNameStandard.ToPtr(),
			},
		},
		nil,
	)

	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	log.Printf("create public ip prefix, name: %v", prefixName)
}

// Create public IP address
func CreatePublicIPAddress(ctx context.Context, addressName string) {
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
		log.Fatalf("failed to obtain a response: %v", err)
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	log.Printf("create public ip address, name: %v", addressName)
}

// Gets the specified public IP prefix in a specified resource group.
func GetPublicIPPrefix(ctx context.Context, ipName string) {
	ipClient := getIPPrefixClient()
	resp, err := ipClient.Get(ctx, config.GroupName(), ipName, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("get public ip prefix, name: %v", *resp.PublicIPPrefix.Name)
}

// Gets the specified public IP address in a specified resource group.
func GetPublicIPAddress(ctx context.Context, ipName string) {
	ipClient := getIPAddressClient()
	resp, err := ipClient.Get(ctx, config.GroupName(), ipName, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("get public ip address, name: %v", *resp.PublicIPAddress.Name)
}

// Gets all public IP addresses in a resource group.
func ListPublicIPPrefix(ctx context.Context) {
	ipClient := getIPPrefixClient()
	pager := ipClient.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			log.Fatalf("failed to obtain a response: %v", pager.Err())
		}
		var resp = pager.PageResponse().PublicIPPrefixListResult
		var b strings.Builder
		if resp != nil && resp.Value != nil {
			for _, v := range *resp.Value {
				b.WriteString(*v.Properties.IPPrefix)
				b.WriteString(",")
			}
			log.Printf("list public ip prefixs in a resource group, IPPrefix: %v\n", strings.TrimRight(b.String(), ","))
		}
	}

	if pager.Err() != nil {
		log.Fatalf("failed to obtain a response: %v", pager.Err())
	}
}

// Gets all the public IP prefixes in a subscription.
func ListPublicIPAddress(ctx context.Context) {
	ipClient := getIPAddressClient()
	pager := ipClient.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			log.Fatalf("failed to obtain a response: %v", pager.Err())
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
		log.Fatalf("failed to obtain a response: %v", pager.Err())
	}
}

// Gets all the public IP prefix in a subscription.
func ListAllPublicIPPrefix(ctx context.Context) {
	ipClient := getIPPrefixClient()
	pager := ipClient.ListAll(nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			log.Fatalf("failed to obtain a response: %v", pager.Err())
		}
		var resp = pager.PageResponse().PublicIPPrefixListResult
		var b strings.Builder
		if resp != nil && resp.Value != nil {
			for _, v := range *resp.Value {
				b.WriteString(*v.Name)
				b.WriteString(",")
			}
			log.Printf("list all public ip prefixs in asubscription, name: %v\n", strings.TrimRight(b.String(), ","))
		}
	}

	if pager.Err() != nil {
		log.Fatalf("failed to obtain a response: %v", pager.Err())
	}
}

// Gets all the public IP addresses in a subscription.
func ListAllPublicIPAddress(ctx context.Context) {
	ipClient := getIPAddressClient()
	pager := ipClient.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			log.Fatalf("failed to obtain a response: %v", pager.Err())
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
		log.Fatalf("failed to obtain a response: %v", pager.Err())
	}
}

// Updates public IP prefix tags.
func UpdatePrefixTags(ctx context.Context, prefixName string) {
	ipClient := getIPPrefixClient()
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("update prefix tags, name: %v\n", *resp.PublicIPPrefix.Name)

}

// Updates public IP address tags.
func UpdateAddressTags(ctx context.Context, prefixName string) {
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("update prefix tags, name: %v\n", *resp.PublicIPAddress.Name)
}

// Deletes the specified public IP prefix.
func DeletePublicIPPrefix(ctx context.Context, prefixName string) {
	ipClient := getIPPrefixClient()
	resp, err := ipClient.BeginDelete(ctx, config.GroupName(), prefixName, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
	log.Printf("delete public prefix, name: %v\n", prefixName)
}

// Deletes the specified public IP address.
func DeletePublicIPAddress(ctx context.Context, addressName string) {
	ipClient := getIPAddressClient()
	resp, err := ipClient.BeginDelete(ctx, config.GroupName(), addressName, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
	log.Printf("delete public address, name: %v\n", addressName)
}
