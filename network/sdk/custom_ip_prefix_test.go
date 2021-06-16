// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"testing"
)

func TestCustomIpPrefix(t *testing.T) {
	/*Seems it’s the permission issue. After tried to register this feature “az feature register --namespace Microsoft.Network --name AllowBringYourOwnIpAddressForThirdParties”,
	seems it requires additional action to onboard the service from service team. So disable it with the reason for now.

	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	customIpPrefixName := config.AppendRandomSuffix("customipprefix")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	customIPPrefixParameters := armnetwork.CustomIPPrefix{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.CustomIPPrefixPropertiesFormat{
			Cidr: to.StringPtr("0.0.0.0/24"),
		},
	}
	err = CreateCustomIpPrefix(ctx, customIpPrefixName, customIPPrefixParameters)
	if err != nil {
		t.Fatalf("failed to create custom ip prefix: % +v", err)
	}
	t.Logf("created custom ip prefix")

	err = GetCustomIpPrefix(ctx, customIpPrefixName)
	if err != nil {
		t.Fatalf("failed to get custom ip prefix: %+v", err)
	}
	t.Logf("got custom ip prefix")

	err = ListCustomIpPrefix(ctx)
	if err != nil {
		t.Fatalf("failed to list custom ip prefix: %+v", err)
	}
	t.Logf("listed custom ip prefix")

	err = ListAllCustomIpPrefix(ctx)
	if err != nil {
		t.Fatalf("failed to list all custom ip prefix: %+v", err)
	}
	t.Logf("listed all custom ip prefix")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateCustomIpPrefixTags(ctx, customIpPrefixName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for custom ip prefix: %+v", err)
	}
	t.Logf("updated custom ip prefix tags")

	err = DeleteCustomIpPrefix(ctx, customIpPrefixName)
	if err != nil {
		t.Fatalf("failed to delete custom ip prefix: %+v", err)
	}
	t.Logf("deleted custom ip prefix")
	*/
}
