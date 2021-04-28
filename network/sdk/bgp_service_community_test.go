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

func TestBGPServiceCommunities(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	err := ListServiceTags(ctx)
	if err != nil {
		t.Fatalf("failed to list bgp service communities: %+v", err)
	}
	t.Logf("listed bgp service communities")
}
