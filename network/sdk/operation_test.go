// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"testing"
	"time"
)

func TestOperation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	err := ListOperation(ctx)
	if err != nil {
		t.Fatalf("failed to list the available REST API operations of the Microsoft.Search provider: %+v", err)
	}
	t.Logf("listed the available REST API operations of the Microsoft.Search provider")
}
