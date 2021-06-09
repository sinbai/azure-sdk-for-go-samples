// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/go-autorest/autorest/to"
)

func TestDscpConfiguration(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	dscpConfigurationName := config.AppendRandomSuffix("dscpconfiguration")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	dscpConfigurationParameters := armnetwork.DscpConfiguration{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.DscpConfigurationPropertiesFormat{
			DestinationIPRanges: &[]*armnetwork.QosIPRange{
				{
					EndIP:   to.StringPtr("127.0.10.2"),
					StartIP: to.StringPtr("127.0.10.1"),
				},
				{
					EndIP:   to.StringPtr("127.0.11.2"),
					StartIP: to.StringPtr("127.0.11.1"),
				},
			},
			DestinationPortRanges: &[]*armnetwork.QosPortRange{
				{
					End:   to.Int32Ptr(15),
					Start: to.Int32Ptr(15),
				},
				{
					End:   to.Int32Ptr(26),
					Start: to.Int32Ptr(27),
				},
			},
			Markings: &[]*int32{to.Int32Ptr(46), to.Int32Ptr(10)},
			Protocol: armnetwork.ProtocolTypeTCP.ToPtr(),
			SourceIPRanges: &[]*armnetwork.QosIPRange{
				{
					EndIP:   to.StringPtr("127.0.0.2"),
					StartIP: to.StringPtr("127.0.0.1"),
				},
				{
					EndIP:   to.StringPtr("127.0.1.2"),
					StartIP: to.StringPtr("127.0.1.1"),
				},
			},
			SourcePortRanges: &[]*armnetwork.QosPortRange{
				{
					End:   to.Int32Ptr(11),
					Start: to.Int32Ptr(10),
				},
				{
					End:   to.Int32Ptr(21),
					Start: to.Int32Ptr(20),
				},
			},
		},
	}
	err = CreateDscpConfiguration(ctx, dscpConfigurationName, dscpConfigurationParameters)
	if err != nil {
		t.Fatalf("failed to create dscp configuration: % +v", err)
	}
	t.Logf("created dscp configuration")

	err = GetDscpConfiguration(ctx, dscpConfigurationName)
	if err != nil {
		t.Fatalf("failed to get dscp configuration: %+v", err)
	}
	t.Logf("got dscp configuration")

	err = ListAllDscpConfiguration(ctx)
	if err != nil {
		t.Fatalf("failed to list all dscp configuration: %+v", err)
	}
	t.Logf("listed all dscp configuration")

	err = DeleteDscpConfiguration(ctx, dscpConfigurationName)
	if err != nil {
		t.Fatalf("failed to delete dscp configuration: %+v", err)
	}
	t.Logf("deleted dscp configuration")
}
