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
)

func TestCheckDnsNameAvailability(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	domainNameLabel := "testdns"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	err := GetCheckDnsNameAvailability(ctx, domainNameLabel)
	if err != nil {
		t.Fatalf("failed to check dns name availability: %+v", err)
	}
	t.Logf("checked dns name availability")
}
