// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

func TestVirtualApplianceSku(t *testing.T) {

	virtualApplianceSkuName := "ciscosdwan"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	err := GetVirtualApplianceSku(ctx, virtualApplianceSkuName)
	if err != nil {
		t.Fatalf("failed to get virtual appliance sku: %+v", err)
	}
	t.Logf("got virtual appliance sku")

	err = ListVirtualApplianceSku(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual appliance sku: %+v", err)
	}
	t.Logf("listed virtual appliance sku")

}
