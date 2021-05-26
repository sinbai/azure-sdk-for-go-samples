// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

func TestVirtualMachineRunCommand(t *testing.T) {
	groupName := config.GenerateGroupName("compute")
	config.SetGroupName(groupName)

	runCommandName := "RunPowerShellScript"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	err := GetVirtualMachineRunCommand(ctx, runCommandName)
	if err != nil {
		t.Fatalf("failed to get virtual machine run command: %+v", err)
	}
	t.Logf("got virtual machine run command")

	err = ListVirtualMachineRunCommand(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual machine run command: %+v", err)
	}
	t.Logf("listed virtual machine run command")

}
