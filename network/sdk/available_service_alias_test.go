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
)

func TestAvailableServiceAlias(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	err := ListAvailableServiceAlias(ctx)
	if err != nil {
		t.Fatalf("failed to list available service alias: %+v", err)
	}
	t.Logf("listed available service alias")

	err = ListAvailableServiceAliasByResourceGroup(ctx)
	if err != nil {
		t.Fatalf("failed to listavailable service alias by resource group: %+v", err)
	}
	t.Logf("listedavailable service alias by resource group")

}
