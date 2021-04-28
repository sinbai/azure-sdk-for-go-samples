// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
)

func TestMain(m *testing.M) {
	err := setupEnvironment()
	if err != nil {
		log.Fatalf("could not set up environment: %v\n", err)
	}

	os.Exit(m.Run())
}

func setupEnvironment() error {
	err1 := config.ParseEnvironment()
	err2 := config.AddFlags()

	for _, err := range []error{err1, err2} {
		if err != nil {
			return err
		}
	}

	flag.Parse()
	return nil
}