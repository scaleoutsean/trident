// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage

import (
	"encoding/json"
	"fmt"
	"testing"

	dvp "github.com/netapp/netappdvp/storage_drivers"

	"github.com/netapp/trident/config"
)

func TestGetCommonInternalVolumeName(t *testing.T) {
	const name = "volume"
	for _, test := range []struct {
		prefix   string
		expected string
	}{
		{
			prefix:   "\"specific\"",
			expected: fmt.Sprintf("specific-%s", name),
		},
		{
			prefix:   "",
			expected: fmt.Sprintf("%s-%s", config.OrchestratorName, name),
		},
		{
			prefix:   "\"\"",
			expected: fmt.Sprintf("%s-%s", config.OrchestratorName, name),
		},
		{
			prefix:   "{}",
			expected: fmt.Sprintf("%s-%s", config.OrchestratorName, name),
		},
	} {
		c := dvp.CommonStorageDriverConfig{
			Version:           1,
			StorageDriverName: "fake",
			StoragePrefixRaw:  json.RawMessage(test.prefix),
		}
		got := GetCommonInternalVolumeName(&c, name)
		if test.expected != got {
			t.Errorf("Mismatch between volume names.  Expected %s, got %s",
				test.expected, got)
		}
	}
}
