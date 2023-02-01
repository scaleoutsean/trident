// Copyright 2022 NetApp, Inc. All Rights Reserved.

package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/netapp/trident/config"
	. "github.com/netapp/trident/logger"
	mockpersistentstore "github.com/netapp/trident/mocks/mock_persistent_store"
	mockstorage "github.com/netapp/trident/mocks/mock_storage"
	persistentstore "github.com/netapp/trident/persistent_store"
	"github.com/netapp/trident/storage"
	"github.com/netapp/trident/storage/fake"
	sa "github.com/netapp/trident/storage_attribute"
	storageclass "github.com/netapp/trident/storage_class"
	fakedriver "github.com/netapp/trident/storage_drivers/fake"
	tu "github.com/netapp/trident/storage_drivers/fake/test_utils"
	"github.com/netapp/trident/utils"
)

var (
	debug = flag.Bool("debug", false, "Enable debugging output")

	inMemoryClient *persistentstore.InMemoryClient
	ctx            = context.Background
)

func init() {
	testing.Init()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	inMemoryClient = persistentstore.NewInMemoryClient()
}

type deleteTest struct {
	name            string
	expectedSuccess bool
}

type recoveryTest struct {
	name           string
	volumeConfig   *storage.VolumeConfig
	snapshotConfig *storage.SnapshotConfig
	expectDestroy  bool
}

func cleanup(t *testing.T, o *TridentOrchestrator) {
	err := o.storeClient.DeleteBackends(ctx())
	if err != nil && !persistentstore.MatchKeyNotFoundErr(err) {
		t.Fatal("Unable to clean up backends: ", err)
	}
	storageClasses, err := o.storeClient.GetStorageClasses(ctx())
	if err != nil && !persistentstore.MatchKeyNotFoundErr(err) {
		t.Fatal("Unable to retrieve storage classes: ", err)
	} else if err == nil {
		for _, psc := range storageClasses {
			sc := storageclass.NewFromPersistent(psc)
			err := o.storeClient.DeleteStorageClass(ctx(), sc)
			if err != nil {
				t.Fatalf("Unable to clean up storage class %s: %v", sc.GetName(), err)
			}
		}
	}
	err = o.storeClient.DeleteVolumes(ctx())
	if err != nil && !persistentstore.MatchKeyNotFoundErr(err) {
		t.Fatal("Unable to clean up volumes: ", err)
	}
	err = o.storeClient.DeleteSnapshots(ctx())
	if err != nil && !persistentstore.MatchKeyNotFoundErr(err) {
		t.Fatal("Unable to clean up snapshots: ", err)
	}

	// Clear the InMemoryClient state so that it looks like we're
	// bootstrapping afresh next time.
	if err = inMemoryClient.Stop(); err != nil {
		t.Fatalf("Unable to stop in memory client for orchestrator: %v", o)
	}
}

func diffConfig(expected, got interface{}, fieldToSkip string) []string {
	diffs := make([]string, 0)
	expectedStruct := reflect.Indirect(reflect.ValueOf(expected))
	gotStruct := reflect.Indirect(reflect.ValueOf(got))

	for i := 0; i < expectedStruct.NumField(); i++ {

		// Optionally skip a field
		typeName := expectedStruct.Type().Field(i).Name
		if typeName == fieldToSkip {
			continue
		}

		// Compare each field in the structs
		expectedField := expectedStruct.FieldByName(typeName).Interface()
		gotField := gotStruct.FieldByName(typeName).Interface()

		if !reflect.DeepEqual(expectedField, gotField) {
			diffs = append(diffs, fmt.Sprintf("%s: expected %v, got %v", typeName, expectedField, gotField))
		}
	}

	return diffs
}

// To be called after reflect.DeepEqual has failed.
func diffExternalBackends(t *testing.T, expected, got *storage.BackendExternal) {
	diffs := make([]string, 0)

	if expected.Name != got.Name {
		diffs = append(diffs, fmt.Sprintf("Name: expected %s, got %s", expected.Name, got.Name))
	}
	if expected.State != got.State {
		diffs = append(diffs, fmt.Sprintf("Online: expected %s, got %s", expected.State, got.State))
	}

	// Diff configs
	expectedConfig := expected.Config
	gotConfig := got.Config

	expectedConfigTypeName := reflect.TypeOf(expectedConfig).Name()
	gotConfigTypeName := reflect.TypeOf(gotConfig).Name()
	if expectedConfigTypeName != gotConfigTypeName {
		t.Errorf("Config type mismatch: %v != %v", expectedConfigTypeName, gotConfigTypeName)
	}

	expectedConfigValue := reflect.ValueOf(expectedConfig)
	gotConfigValue := reflect.ValueOf(gotConfig)

	expectedCSDCIntf := expectedConfigValue.FieldByName("CommonStorageDriverConfig").Interface()
	gotCSDCIntf := gotConfigValue.FieldByName("CommonStorageDriverConfig").Interface()

	var configDiffs []string

	// Compare the common storage driver config
	configDiffs = diffConfig(expectedCSDCIntf, gotCSDCIntf, "")
	diffs = append(diffs, configDiffs...)

	// Compare the base config, without the common storage driver config
	configDiffs = diffConfig(expectedConfig, gotConfig, "CommonStorageDriverConfig")
	diffs = append(diffs, configDiffs...)

	t.Logf("expectedConfig %v", expectedConfig)
	t.Logf("gotConfig %v", gotConfig)
	t.Logf("diffs %v", diffs)
	// Diff storage
	for name, expectedVC := range expected.Storage {
		if gotVC, ok := got.Storage[name]; !ok {
			diffs = append(diffs, fmt.Sprintf("Storage: did not get expected VC %s", name))
		} else if !reflect.DeepEqual(expectedVC, gotVC) {
			expectedJSON, err := json.Marshal(expectedVC)
			if err != nil {
				t.Fatal("Unable to marshal expected JSON for VC ", name)
			}
			gotJSON, err := json.Marshal(gotVC)
			if err != nil {
				t.Fatal("Unable to marshal got JSON for VC ", name)
			}
			diffs = append(
				diffs, fmt.Sprintf(
					"Storage: pool %s differs:\n\t\t"+
						"Expected: %s\n\t\tGot: %s", name, string(expectedJSON), string(gotJSON),
				),
			)
		}
	}
	for name := range got.Storage {
		if _, ok := expected.Storage[name]; !ok {
			diffs = append(diffs, fmt.Sprintf("Storage: got unexpected VC %s", name))
		}
	}

	// Diff volumes
	expectedVolMap := make(map[string]bool, len(expected.Volumes))
	gotVolMap := make(map[string]bool, len(got.Volumes))
	for _, v := range expected.Volumes {
		expectedVolMap[v] = true
	}
	for _, v := range got.Volumes {
		gotVolMap[v] = true
	}
	for name := range expectedVolMap {
		if _, ok := gotVolMap[name]; !ok {
			diffs = append(diffs, fmt.Sprintf("Volumes: did not get expected volume %s", name))
		}
	}
	for name := range gotVolMap {
		if _, ok := expectedVolMap[name]; !ok {
			diffs = append(diffs, fmt.Sprintf("Volumes: got unexpected volume %s", name))
		}
	}
	if len(diffs) > 0 {
		t.Errorf("External backends differ:\n\t%s", strings.Join(diffs, "\n\t"))
	}
}

func runDeleteTest(
	t *testing.T, d *deleteTest, orchestrator *TridentOrchestrator,
) {
	var (
		backendUUID string
		backend     storage.Backend
		found       bool
	)
	if d.expectedSuccess {
		orchestrator.mutex.Lock()

		backendUUID = orchestrator.volumes[d.name].BackendUUID
		backend, found = orchestrator.backends[backendUUID]
		if !found {
			t.Errorf("Backend %v isn't managed by the orchestrator!", backendUUID)
		}
		if _, found = backend.Volumes()[d.name]; !found {
			t.Errorf("Volume %s doesn't exist on backend %s!", d.name, backendUUID)
		}
		orchestrator.mutex.Unlock()
	}
	err := orchestrator.DeleteVolume(ctx(), d.name)
	if err == nil && !d.expectedSuccess {
		t.Errorf("%s: volume delete succeeded when it should not have.", d.name)
	} else if err != nil && d.expectedSuccess {
		t.Errorf("%s: delete failed: %v", d.name, err)
	} else if d.expectedSuccess {
		volume, err := orchestrator.GetVolume(ctx(), d.name)
		if volume != nil || err == nil {
			t.Errorf("%s: got volume where none expected.", d.name)
		}
		orchestrator.mutex.Lock()
		if _, found = backend.Volumes()[d.name]; found {
			t.Errorf("Volume %s shouldn't exist on backend %s!", d.name, backendUUID)
		}
		externalVol, err := orchestrator.storeClient.GetVolume(ctx(), d.name)
		if err != nil {
			if !persistentstore.MatchKeyNotFoundErr(err) {
				t.Errorf(
					"%s: unable to communicate with backing store: "+
						"%v", d.name, err,
				)
			}
			// We're successful if we get to here; we expect an
			// ErrorCodeKeyNotFound.
		} else if externalVol != nil {
			t.Errorf("%s: volume not properly deleted from backing store", d.name)
		}
		orchestrator.mutex.Unlock()
	}
}

type storageClassTest struct {
	config   *storageclass.Config
	expected []*tu.PoolMatch
}

func getOrchestrator(t *testing.T, monitorTransactions bool) *TridentOrchestrator {
	var (
		storeClient persistentstore.Client
		err         error
	)
	// This will have been created as not nil in init
	// We can't create a new one here because tests that exercise
	// bootstrapping need to have their data persist.
	storeClient = inMemoryClient

	o := NewTridentOrchestrator(storeClient)
	if err = o.Bootstrap(monitorTransactions); err != nil {
		t.Fatal("Failure occurred during bootstrapping: ", err)
	}
	return o
}

func validateStorageClass(
	t *testing.T,
	o *TridentOrchestrator,
	name string,
	expected []*tu.PoolMatch,
) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	sc, ok := o.storageClasses[name]
	if !ok {
		t.Errorf("%s: Storage class not found in backend.", name)
	}
	remaining := make([]*tu.PoolMatch, len(expected))
	copy(remaining, expected)
	for _, protocol := range []config.Protocol{config.File, config.Block, config.BlockOnFile} {
		for _, pool := range sc.GetStoragePoolsForProtocol(ctx(), protocol, config.ReadWriteOnce) {
			nameFound := false
			for _, scName := range pool.StorageClasses() {
				if scName == name {
					nameFound = true
					break
				}
			}
			if !nameFound {
				t.Errorf("%s: Storage class name not found in storage pool %s", name, pool.Name())
			}
			matchIndex := -1
			for i, r := range remaining {
				if r.Matches(pool) {
					matchIndex = i
					break
				}
			}
			if matchIndex >= 0 {
				// If we match, remove the match from the potential matches.
				remaining[matchIndex] = remaining[len(remaining)-1]
				remaining[len(remaining)-1] = nil
				remaining = remaining[:len(remaining)-1]
			} else {
				t.Errorf("%s: Found unexpected match for storage class: %s:%s", name, pool.Backend().Name(),
					pool.Name())
			}
		}
	}
	if len(remaining) > 0 {
		remainingNames := make([]string, len(remaining))
		for i, r := range remaining {
			remainingNames[i] = r.String()
		}
		t.Errorf("%s: Storage class failed to match storage pools %s", name, strings.Join(remainingNames, ", "))
	}
	persistentSC, err := o.storeClient.GetStorageClass(ctx(), name)
	if err != nil {
		t.Fatalf("Unable to get storage class %s from backend: %v", name, err)
	}
	if !reflect.DeepEqual(
		persistentSC,
		sc.ConstructPersistent(),
	) {
		gotSCJSON, err := json.Marshal(persistentSC)
		if err != nil {
			t.Fatalf("Unable to marshal persisted storage class %s: %v", name, err)
		}
		expectedSCJSON, err := json.Marshal(sc.ConstructPersistent())
		if err != nil {
			t.Fatalf("Unable to marshal expected persistent storage class %s: %v", name, err)
		}
		t.Errorf("%s: Storage class persisted incorrectly.\n\tExpected %s\n\tGot %s", name, expectedSCJSON, gotSCJSON)
	}
}

// This test is fairly heavyweight, but, due to the need to accumulate state
// to run the later tests, it's easier to do this all in one go at the moment.
// Consider breaking this up if it gets unwieldy, though.
func TestAddStorageClassVolumes(t *testing.T) {
	mockPools := tu.GetFakePools()
	orchestrator := getOrchestrator(t, false)

	errored := false
	for _, c := range []struct {
		name      string
		protocol  config.Protocol
		poolNames []string
	}{
		{
			name:      "fast-a",
			protocol:  config.File,
			poolNames: []string{tu.FastSmall, tu.FastThinOnly},
		},
		{
			name:      "fast-b",
			protocol:  config.File,
			poolNames: []string{tu.FastThinOnly, tu.FastUniqueAttr},
		},
		{
			name:      "slow-file",
			protocol:  config.File,
			poolNames: []string{tu.SlowNoSnapshots, tu.SlowSnapshots},
		},
		{
			name:      "slow-block",
			protocol:  config.Block,
			poolNames: []string{tu.SlowNoSnapshots, tu.SlowSnapshots, tu.MediumOverlap},
		},
	} {
		pools := make(map[string]*fake.StoragePool, len(c.poolNames))
		for _, poolName := range c.poolNames {
			pools[poolName] = mockPools[poolName]
		}
		volumes := make([]fake.Volume, 0)
		fakeConfig, err := fakedriver.NewFakeStorageDriverConfigJSON(c.name, c.protocol, pools, volumes)
		if err != nil {
			t.Fatalf("Unable to generate config JSON for %s: %v", c.name, err)
		}
		_, err = orchestrator.AddBackend(ctx(), fakeConfig, "")
		if err != nil {
			t.Errorf("Unable to add backend %s: %v", c.name, err)
			errored = true
		}
		orchestrator.mutex.Lock()
		backend, err := orchestrator.getBackendByBackendName(c.name)
		if err != nil {
			t.Fatalf("Backend %s not stored in orchestrator, err %s", c.name, err)
		}
		persistentBackend, err := orchestrator.storeClient.GetBackend(ctx(), c.name)
		if err != nil {
			t.Fatalf("Unable to get backend %s from persistent store: %v", c.name, err)
		} else if !reflect.DeepEqual(backend.ConstructPersistent(ctx()), persistentBackend) {
			t.Error("Wrong data stored for backend ", c.name)
		}
		orchestrator.mutex.Unlock()
	}
	if errored {
		t.Fatal("Failed to add all backends; aborting remaining tests.")
	}

	// Add storage classes
	scTests := []storageClassTest{
		{
			config: &storageclass.Config{
				Name: "slow",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(40),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "slow-file", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "fast",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
		{
			config: &storageclass.Config{
				Name: "fast-unique",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
					sa.UniqueOptions:    sa.NewStringRequest("baz"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
		{
			config: &storageclass.Config{
				Name: "pools",
				Pools: map[string][]string{
					"fast-a":     {tu.FastSmall},
					"slow-block": {tu.SlowNoSnapshots, tu.MediumOverlap},
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
		},
		{
			config: &storageclass.Config{
				Name: "additionalPools",
				AdditionalPools: map[string][]string{
					"fast-a":     {tu.FastThinOnly},
					"slow-block": {tu.SlowNoSnapshots, tu.MediumOverlap},
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
		},
		{
			config: &storageclass.Config{
				Name: "poolsWithAttributes",
				Attributes: map[string]sa.Request{
					sa.IOPS:      sa.NewIntRequest(2000),
					sa.Snapshots: sa.NewBoolRequest(true),
				},
				Pools: map[string][]string{
					"fast-a":     {tu.FastThinOnly},
					"slow-block": {tu.SlowNoSnapshots, tu.MediumOverlap},
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastThinOnly},
			},
		},
		{
			config: &storageclass.Config{
				Name: "additionalPoolsWithAttributes",
				Attributes: map[string]sa.Request{
					sa.IOPS:      sa.NewIntRequest(2000),
					sa.Snapshots: sa.NewBoolRequest(true),
				},
				AdditionalPools: map[string][]string{
					"fast-a":     {tu.FastThinOnly},
					"slow-block": {tu.SlowNoSnapshots},
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "additionalPoolsWithAttributesAndPools",
				Attributes: map[string]sa.Request{
					sa.IOPS:      sa.NewIntRequest(2000),
					sa.Snapshots: sa.NewBoolRequest(true),
				},
				Pools: map[string][]string{
					"fast-a":     {tu.FastThinOnly},
					"slow-block": {tu.SlowNoSnapshots, tu.MediumOverlap},
				},
				AdditionalPools: map[string][]string{
					"fast-b":     {tu.FastThinOnly},
					"slow-block": {tu.SlowNoSnapshots},
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "additionalPoolsNoMatch",
				AdditionalPools: map[string][]string{
					"unknown": {tu.FastThinOnly},
				},
			},
			expected: []*tu.PoolMatch{},
		},
		{
			config: &storageclass.Config{
				Name: "mixed",
				AdditionalPools: map[string][]string{
					"slow-file": {tu.SlowNoSnapshots},
					"fast-b":    {tu.FastThinOnly, tu.FastUniqueAttr},
				},
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
				{Backend: "slow-file", Pool: tu.SlowNoSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "emptyStorageClass",
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
				{Backend: "slow-file", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-file", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
		},
	}
	for _, s := range scTests {
		_, err := orchestrator.AddStorageClass(ctx(), s.config)
		if err != nil {
			t.Errorf("Unable to add storage class %s: %v", s.config.Name, err)
			continue
		}
		validateStorageClass(t, orchestrator, s.config.Name, s.expected)
	}
	preSCDeleteTests := make([]*deleteTest, 0)
	postSCDeleteTests := make([]*deleteTest, 0)
	for _, s := range []struct {
		name            string
		config          *storage.VolumeConfig
		expectedSuccess bool
		expectedMatches []*tu.PoolMatch
		expectedCount   int
		deleteAfterSC   bool
	}{
		{
			name:            "basic",
			config:          tu.GenerateVolumeConfig("basic", 1, "fast", config.File),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
			expectedCount: 1,
			deleteAfterSC: false,
		},
		{
			name:            "large",
			config:          tu.GenerateVolumeConfig("large", 100, "fast", config.File),
			expectedSuccess: false,
			expectedMatches: []*tu.PoolMatch{},
			expectedCount:   0,
			deleteAfterSC:   false,
		},
		{
			name:            "block",
			config:          tu.GenerateVolumeConfig("block", 1, "pools", config.Block),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
			expectedCount: 1,
			deleteAfterSC: false,
		},
		{
			name:            "block2",
			config:          tu.GenerateVolumeConfig("block2", 1, "additionalPools", config.Block),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
			expectedCount: 1,
			deleteAfterSC: false,
		},
		{
			name:            "invalid-storage-class",
			config:          tu.GenerateVolumeConfig("invalid", 1, "nonexistent", config.File),
			expectedSuccess: false,
			expectedMatches: []*tu.PoolMatch{},
			expectedCount:   0,
			deleteAfterSC:   false,
		},
		{
			name:            "repeated",
			config:          tu.GenerateVolumeConfig("basic", 20, "fast", config.File),
			expectedSuccess: false,
			expectedMatches: []*tu.PoolMatch{},
			expectedCount:   1,
			deleteAfterSC:   false,
		},
		{
			name:            "postSCDelete",
			config:          tu.GenerateVolumeConfig("postSCDelete", 20, "fast", config.File),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
			expectedCount: 1,
			deleteAfterSC: false,
		},
	} {
		vol, err := orchestrator.AddVolume(ctx(), s.config)
		if err != nil && s.expectedSuccess {
			t.Errorf("%s: got unexpected error %v", s.name, err)
			continue
		} else if err == nil && !s.expectedSuccess {
			t.Errorf("%s: volume create succeeded unexpectedly.", s.name)
			continue
		}
		orchestrator.mutex.Lock()
		volume, found := orchestrator.volumes[s.config.Name]
		if s.expectedCount == 1 && !found {
			t.Errorf("%s: did not get volume where expected.", s.name)
		} else if s.expectedCount == 0 && found {
			t.Errorf("%s: got a volume where none expected.", s.name)
		}
		if !s.expectedSuccess {
			deleteTest := &deleteTest{
				name:            s.config.Name,
				expectedSuccess: false,
			}
			if s.deleteAfterSC {
				postSCDeleteTests = append(postSCDeleteTests, deleteTest)
			} else {
				preSCDeleteTests = append(preSCDeleteTests, deleteTest)
			}
			orchestrator.mutex.Unlock()
			continue
		}
		matched := false
		for _, potentialMatch := range s.expectedMatches {
			volumeBackend, err := orchestrator.getBackendByBackendUUID(volume.BackendUUID)
			if volumeBackend == nil || err != nil {
				continue
			}
			if potentialMatch.Backend == volumeBackend.Name() &&
				potentialMatch.Pool == volume.Pool {
				matched = true
				deleteTest := &deleteTest{
					name:            s.config.Name,
					expectedSuccess: true,
				}
				if s.deleteAfterSC {
					postSCDeleteTests = append(postSCDeleteTests, deleteTest)
				} else {
					preSCDeleteTests = append(preSCDeleteTests, deleteTest)
				}
				break
			}
		}
		if !matched {
			t.Errorf(
				"%s: Volume placed on unexpected backend and storage pool: %s, %s",
				s.name,
				volume.BackendUUID,
				volume.Pool,
			)
		}

		externalVolume, err := orchestrator.storeClient.GetVolume(ctx(), s.config.Name)
		if err != nil {
			t.Errorf("%s: unable to communicate with backing store: %v", s.name, err)
		}
		if !reflect.DeepEqual(externalVolume, vol) {
			t.Errorf("%s: external volume %s stored in backend does not match created volume.", s.name,
				externalVolume.Config.Name)
			externalVolJSON, err := json.Marshal(externalVolume)
			if err != nil {
				t.Fatal("Unable to remarshal JSON: ", err)
			}
			origVolJSON, err := json.Marshal(vol)
			if err != nil {
				t.Fatal("Unable to remarshal JSON: ", err)
			}
			t.Logf("\tExpected: %s\n\tGot: %s\n", string(externalVolJSON), string(origVolJSON))
		}
		orchestrator.mutex.Unlock()
	}
	for _, d := range preSCDeleteTests {
		runDeleteTest(t, d, orchestrator)
	}

	// Delete storage classes.  Note: there are currently no error cases.
	for _, s := range scTests {
		err := orchestrator.DeleteStorageClass(ctx(), s.config.Name)
		if err != nil {
			t.Errorf("%s delete: Unable to remove storage class: %v", s.config.Name, err)
		}
		orchestrator.mutex.Lock()
		if _, ok := orchestrator.storageClasses[s.config.Name]; ok {
			t.Errorf("%s delete: Storage class still found in map.", s.config.Name)
		}
		// Ensure that the storage class was cleared from its backends.
		for _, poolMatch := range s.expected {
			b, err := orchestrator.getBackendByBackendName(poolMatch.Backend)
			if b == nil || err != nil {
				t.Errorf("%s delete: backend %s not found in orchestrator.", s.config.Name, poolMatch.Backend)
				continue
			}
			p, ok := b.Storage()[poolMatch.Pool]
			if !ok {
				t.Errorf("%s delete: storage pool %s not found for backend %s", s.config.Name, poolMatch.Pool,
					poolMatch.Backend)
				continue
			}
			for _, sc := range p.StorageClasses() {
				if sc == s.config.Name {
					t.Errorf("%s delete: storage class name not removed from backend %s, storage pool %s",
						s.config.Name, poolMatch.Backend, poolMatch.Pool)
				}
			}
		}
		externalSC, err := orchestrator.storeClient.GetStorageClass(ctx(), s.config.Name)
		if err != nil {
			if !persistentstore.MatchKeyNotFoundErr(err) {
				t.Errorf("%s: unable to communicate with backing store: %v", s.config.Name, err)
			}
			// We're successful if we get to here; we expect an
			// ErrorCodeKeyNotFound.
		} else if externalSC != nil {
			t.Errorf("%s: storageClass not properly deleted from backing store", s.config.Name)
		}
		orchestrator.mutex.Unlock()
	}
	for _, d := range postSCDeleteTests {
		runDeleteTest(t, d, orchestrator)
	}
	cleanup(t, orchestrator)
}

func TestUpdateVolume_LUKSPassphraseNames(t *testing.T) {
	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: luksPassphraseNames field updated
	orchestrator := getOrchestrator(t, false)
	vol := &storage.Volume{
		Config:      &storage.VolumeConfig{Name: "test-vol", LUKSPassphraseNames: []string{}},
		BackendUUID: "12345",
	}
	orchestrator.volumes[vol.Config.Name] = vol
	err := orchestrator.storeClient.AddVolume(context.TODO(), vol)
	assert.NoError(t, err)
	assert.Empty(t, orchestrator.volumes[vol.Config.Name].Config.LUKSPassphraseNames)

	err = orchestrator.UpdateVolume(context.TODO(), "test-vol", &[]string{"A"})
	desiredPassphraseNames := []string{"A"}
	assert.NoError(t, err)
	assert.Equal(t, desiredPassphraseNames, orchestrator.volumes[vol.Config.Name].Config.LUKSPassphraseNames)

	storedVol, err := orchestrator.storeClient.GetVolume(context.TODO(), "test-vol")
	assert.NoError(t, err)
	assert.Equal(t, desiredPassphraseNames, storedVol.Config.LUKSPassphraseNames)

	err = orchestrator.storeClient.DeleteVolume(context.TODO(), vol)
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: luksPassphraseNames field nil
	orchestrator = getOrchestrator(t, false)
	vol = &storage.Volume{
		Config:      &storage.VolumeConfig{Name: "test-vol", LUKSPassphraseNames: []string{}},
		BackendUUID: "12345",
	}
	orchestrator.volumes[vol.Config.Name] = vol
	err = orchestrator.storeClient.AddVolume(context.TODO(), vol)
	assert.NoError(t, err)
	assert.Empty(t, orchestrator.volumes[vol.Config.Name].Config.LUKSPassphraseNames)

	err = orchestrator.UpdateVolume(context.TODO(), "test-vol", nil)
	desiredPassphraseNames = []string{}
	assert.NoError(t, err)
	assert.Equal(t, desiredPassphraseNames, orchestrator.volumes[vol.Config.Name].Config.LUKSPassphraseNames)

	storedVol, err = orchestrator.storeClient.GetVolume(context.TODO(), "test-vol")
	assert.NoError(t, err)
	assert.Equal(t, desiredPassphraseNames, storedVol.Config.LUKSPassphraseNames)

	err = orchestrator.storeClient.DeleteVolume(context.TODO(), vol)
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: failed to update persistence, volume not found
	orchestrator = getOrchestrator(t, false)
	vol = &storage.Volume{
		Config:      &storage.VolumeConfig{Name: "test-vol", LUKSPassphraseNames: []string{}},
		BackendUUID: "12345",
	}
	orchestrator.volumes[vol.Config.Name] = vol

	err = orchestrator.UpdateVolume(context.TODO(), "test-vol", &[]string{"A"})
	desiredPassphraseNames = []string{}
	assert.Error(t, err)
	assert.Equal(t, desiredPassphraseNames, orchestrator.volumes[vol.Config.Name].Config.LUKSPassphraseNames)

	_, err = orchestrator.storeClient.GetVolume(context.TODO(), "test-vol")
	// Not found
	assert.Error(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: bootstrap error
	orchestrator = getOrchestrator(t, false)
	bootstrapError := fmt.Errorf("my bootstrap error")
	orchestrator.bootstrapError = bootstrapError

	err = orchestrator.UpdateVolume(context.TODO(), "test-vol", &[]string{"A"})
	assert.Error(t, err)
	assert.ErrorIs(t, err, bootstrapError)
}

func TestCloneVolume_SnapshotDataSource_LUKS(t *testing.T) {
	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: luksPassphraseNames field updated
	// // Setup
	mockPools := tu.GetFakePools()
	orchestrator := getOrchestrator(t, false)

	// Make a backend
	poolNames := []string{tu.SlowSnapshots}
	pools := make(map[string]*fake.StoragePool, len(poolNames))
	for _, poolName := range poolNames {
		pools[poolName] = mockPools[poolName]
	}
	volumes := make([]fake.Volume, 0)
	cfg, err := fakedriver.NewFakeStorageDriverConfigJSON("slow-block", "block", pools, volumes)
	assert.NoError(t, err)
	_, err = orchestrator.AddBackend(ctx(), cfg, "")
	assert.NoError(t, err)
	defer orchestrator.DeleteBackend(ctx(), "slow-block")

	// Make a StorageClass
	storageClass := &storageclass.Config{Name: "specific"}
	_, err = orchestrator.AddStorageClass(ctx(), storageClass)
	defer orchestrator.DeleteStorageClass(ctx(), storageClass.Name)
	assert.NoError(t, err)

	// Create the original volume
	volConfig := tu.GenerateVolumeConfig("block", 1, "specific", config.Block)
	volConfig.LUKSEncryption = "true"
	volConfig.LUKSPassphraseNames = []string{"A", "B"}
	_, err = orchestrator.AddVolume(ctx(), volConfig)
	assert.NoError(t, err)
	defer orchestrator.DeleteVolume(ctx(), volConfig.Name)

	// Create a snapshot
	snapshotConfig := generateSnapshotConfig("test-snapshot", volConfig.Name, volConfig.InternalName)
	snapshot, err := orchestrator.CreateSnapshot(ctx(), snapshotConfig)
	assert.NoError(t, err)
	assert.Equal(t, snapshot.Config.LUKSPassphraseNames, []string{"A", "B"})
	defer orchestrator.DeleteSnapshot(ctx(), volConfig.Name, snapshotConfig.Name)

	// "rotate" the luksPassphraseNames of the volume
	err = orchestrator.UpdateVolume(ctx(), volConfig.Name, &[]string{"A"})
	assert.NoError(t, err)
	vol, err := orchestrator.GetVolume(ctx(), volConfig.Name)
	assert.NoError(t, err)
	volConfig = vol.Config
	assert.Equal(t, vol.Config.LUKSPassphraseNames, []string{"A"})

	// Now clone the snapshot and ensure everything looks fine
	cloneName := volConfig.Name + "_clone"
	cloneConfig := &storage.VolumeConfig{
		Name:                cloneName,
		StorageClass:        volConfig.StorageClass,
		CloneSourceVolume:   volConfig.Name,
		CloneSourceSnapshot: snapshotConfig.Name,
		VolumeMode:          volConfig.VolumeMode,
	}
	cloneResult, err := orchestrator.CloneVolume(ctx(), cloneConfig)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, cloneResult.Config.LUKSPassphraseNames)
	defer orchestrator.DeleteVolume(ctx(), cloneResult.Config.Name)
}

func TestCloneVolume_VolumeDataSource_LUKS(t *testing.T) {
	// // Setup
	mockPools := tu.GetFakePools()
	orchestrator := getOrchestrator(t, false)

	// Create backend
	poolNames := []string{tu.SlowSnapshots}
	pools := make(map[string]*fake.StoragePool, len(poolNames))
	for _, poolName := range poolNames {
		pools[poolName] = mockPools[poolName]
	}
	volumes := make([]fake.Volume, 0)
	cfg, err := fakedriver.NewFakeStorageDriverConfigJSON("slow-block", "block", pools, volumes)
	assert.NoError(t, err)
	_, err = orchestrator.AddBackend(ctx(), cfg, "")
	assert.NoError(t, err)
	defer orchestrator.DeleteBackend(ctx(), "slow-block")

	// Create a StorageClass
	storageClass := &storageclass.Config{Name: "specific"}
	_, err = orchestrator.AddStorageClass(ctx(), storageClass)
	defer orchestrator.DeleteStorageClass(ctx(), storageClass.Name)
	assert.NoError(t, err)

	// Create the Volume
	volConfig := tu.GenerateVolumeConfig("block", 1, "specific", config.Block)
	volConfig.LUKSEncryption = "true"
	volConfig.LUKSPassphraseNames = []string{"A", "B"}
	// Create the source volume
	_, err = orchestrator.AddVolume(ctx(), volConfig)
	assert.NoError(t, err)
	defer orchestrator.DeleteVolume(ctx(), volConfig.Name)

	// Create a snapshot
	snapshotConfig := generateSnapshotConfig("test-snapshot", volConfig.Name, volConfig.InternalName)
	snapshot, err := orchestrator.CreateSnapshot(ctx(), snapshotConfig)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, snapshot.Config.LUKSPassphraseNames)
	defer orchestrator.DeleteSnapshot(ctx(), volConfig.Name, snapshotConfig.Name)

	// Now clone the volume and ensure everything looks fine
	cloneName := volConfig.Name + "_clone"
	cloneConfig := &storage.VolumeConfig{
		Name:              cloneName,
		StorageClass:      volConfig.StorageClass,
		CloneSourceVolume: volConfig.Name,
		VolumeMode:        volConfig.VolumeMode,
	}
	cloneResult, err := orchestrator.CloneVolume(ctx(), cloneConfig)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, cloneResult.Config.LUKSPassphraseNames)
	defer orchestrator.DeleteVolume(ctx(), cloneResult.Config.Name)
}

// This test is modeled after TestAddStorageClassVolumes, but we don't need all the
// tests around storage class deletion, etc.
func TestCloneVolumes(t *testing.T) {
	mockPools := tu.GetFakePools()
	orchestrator := getOrchestrator(t, false)

	errored := false
	for _, c := range []struct {
		name      string
		protocol  config.Protocol
		poolNames []string
	}{
		{
			name:      "fast-a",
			protocol:  config.File,
			poolNames: []string{tu.FastSmall, tu.FastThinOnly},
		},
		{
			name:      "fast-b",
			protocol:  config.File,
			poolNames: []string{tu.FastThinOnly, tu.FastUniqueAttr},
		},
		{
			name:      "slow-file",
			protocol:  config.File,
			poolNames: []string{tu.SlowNoSnapshots, tu.SlowSnapshots},
		},
		{
			name:      "slow-block",
			protocol:  config.Block,
			poolNames: []string{tu.SlowNoSnapshots, tu.SlowSnapshots, tu.MediumOverlap},
		},
	} {
		pools := make(map[string]*fake.StoragePool, len(c.poolNames))
		for _, poolName := range c.poolNames {
			pools[poolName] = mockPools[poolName]
		}

		volumes := make([]fake.Volume, 0)
		cfg, err := fakedriver.NewFakeStorageDriverConfigJSON(
			c.name, c.protocol,
			pools, volumes,
		)
		if err != nil {
			t.Fatalf("Unable to generate cfg JSON for %s: %v", c.name, err)
		}
		_, err = orchestrator.AddBackend(ctx(), cfg, "")
		if err != nil {
			t.Errorf("Unable to add backend %s: %v", c.name, err)
			errored = true
		}
		orchestrator.mutex.Lock()
		backend, err := orchestrator.getBackendByBackendName(c.name)
		if backend == nil || err != nil {
			t.Fatalf("Backend %s not stored in orchestrator", c.name)
		}
		persistentBackend, err := orchestrator.storeClient.GetBackend(ctx(), c.name)
		if err != nil {
			t.Fatalf("Unable to get backend %s from persistent store: %v", c.name, err)
		} else if !reflect.DeepEqual(backend.ConstructPersistent(ctx()), persistentBackend) {
			t.Error("Wrong data stored for backend ", c.name)
		}
		orchestrator.mutex.Unlock()
	}
	if errored {
		t.Fatal("Failed to add all backends; aborting remaining tests.")
	}

	// Add storage classes
	storageClasses := []storageClassTest{
		{
			config: &storageclass.Config{
				Name: "slow",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(40),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "slow-file", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "fast",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
		{
			config: &storageclass.Config{
				Name: "fast-unique",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
					sa.UniqueOptions:    sa.NewStringRequest("baz"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
		{
			config: &storageclass.Config{
				Name: "specific",
				AdditionalPools: map[string][]string{
					"fast-a":     {tu.FastThinOnly},
					"slow-block": {tu.SlowNoSnapshots, tu.MediumOverlap},
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
		},
		{
			config: &storageclass.Config{
				Name: "specificNoMatch",
				AdditionalPools: map[string][]string{
					"unknown": {tu.FastThinOnly},
				},
			},
			expected: []*tu.PoolMatch{},
		},
		{
			config: &storageclass.Config{
				Name: "mixed",
				AdditionalPools: map[string][]string{
					"slow-file": {tu.SlowNoSnapshots},
					"fast-b":    {tu.FastThinOnly, tu.FastUniqueAttr},
				},
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
				{Backend: "slow-file", Pool: tu.SlowNoSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "emptyStorageClass",
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
				{Backend: "slow-file", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-file", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
		},
	}
	for _, s := range storageClasses {
		_, err := orchestrator.AddStorageClass(ctx(), s.config)
		if err != nil {
			t.Errorf("Unable to add storage class %s: %v", s.config.Name, err)
			continue
		}
		validateStorageClass(t, orchestrator, s.config.Name, s.expected)
	}

	for _, s := range []struct {
		name            string
		config          *storage.VolumeConfig
		expectedSuccess bool
		expectedMatches []*tu.PoolMatch
	}{
		{
			name:            "file",
			config:          tu.GenerateVolumeConfig("file", 1, "fast", config.File),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
		{
			name:            "block",
			config:          tu.GenerateVolumeConfig("block", 1, "specific", config.Block),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "slow-block", Pool: tu.SlowNoSnapshots},
				{Backend: "slow-block", Pool: tu.MediumOverlap},
			},
		},
	} {
		// Create the source volume
		_, err := orchestrator.AddVolume(ctx(), s.config)
		if err != nil {
			t.Errorf("%s: got unexpected error %v", s.name, err)
			continue
		}

		// Now clone the volume and ensure everything looks fine
		cloneName := s.config.Name + "_clone"
		cloneConfig := &storage.VolumeConfig{
			Name:              cloneName,
			StorageClass:      s.config.StorageClass,
			CloneSourceVolume: s.config.Name,
			VolumeMode:        s.config.VolumeMode,
		}
		cloneResult, err := orchestrator.CloneVolume(ctx(), cloneConfig)
		if err != nil {
			t.Errorf("%s: got unexpected error %v", s.name, err)
			continue
		}

		orchestrator.mutex.Lock()

		volume, found := orchestrator.volumes[s.config.Name]
		if !found {
			t.Errorf("%s: did not get volume where expected.", s.name)
		}
		clone, found := orchestrator.volumes[cloneName]
		if !found {
			t.Errorf("%s: did not get volume clone where expected.", cloneName)
		}

		// Clone must reside in the same place as the source
		if clone.BackendUUID != volume.BackendUUID {
			t.Errorf("%s: Clone placed on unexpected backend: %s", cloneName, clone.BackendUUID)
		}

		// Clone should be registered in the store just like any other volume
		externalClone, err := orchestrator.storeClient.GetVolume(ctx(), cloneName)
		if err != nil {
			t.Errorf("%s: unable to communicate with backing store: %v", cloneName, err)
		}
		if !reflect.DeepEqual(externalClone, cloneResult) {
			t.Errorf("%s: external volume %s stored in backend does not match created volume.", cloneName,
				externalClone.Config.Name)
			externalCloneJSON, err := json.Marshal(externalClone)
			if err != nil {
				t.Fatal("Unable to remarshal JSON: ", err)
			}
			origCloneJSON, err := json.Marshal(cloneResult)
			if err != nil {
				t.Fatal("Unable to remarshal JSON: ", err)
			}
			t.Logf("\tExpected: %s\n\tGot: %s\n", string(externalCloneJSON), string(origCloneJSON))
		}

		orchestrator.mutex.Unlock()
	}

	cleanup(t, orchestrator)
}

func addBackend(
	t *testing.T, orchestrator *TridentOrchestrator, backendName string, backendProtocol config.Protocol,
) {
	volumes := []fake.Volume{
		{Name: "origVolume01", RequestedPool: "primary", PhysicalPool: "primary", SizeBytes: 1000000000},
		{Name: "origVolume02", RequestedPool: "primary", PhysicalPool: "primary", SizeBytes: 1000000000},
	}
	configJSON, err := fakedriver.NewFakeStorageDriverConfigJSON(
		backendName,
		backendProtocol,
		map[string]*fake.StoragePool{
			"primary": {
				Attrs: map[string]sa.Offer{
					sa.Media:            sa.NewStringOffer("hdd"),
					sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
					// testingAttribute is here to ensure that only one
					// storage class will match this backend.
					sa.TestingAttribute: sa.NewBoolOffer(true),
				},
				Bytes: 100 * 1024 * 1024 * 1024,
			},
		},
		volumes,
	)
	if err != nil {
		t.Fatal("Unable to create mock driver config JSON: ", err)
	}
	_, err = orchestrator.AddBackend(ctx(), configJSON, "")
	if err != nil {
		t.Fatalf("Unable to add initial backend: %v", err)
	}
}

// addBackendStorageClass creates a backend and storage class for tests
// that don't care deeply about this functionality.
func addBackendStorageClass(
	t *testing.T,
	orchestrator *TridentOrchestrator,
	backendName string,
	scName string,
	backendProtocol config.Protocol,
) {
	addBackend(t, orchestrator, backendName, backendProtocol)
	_, err := orchestrator.AddStorageClass(
		ctx(), &storageclass.Config{
			Name: scName,
			Attributes: map[string]sa.Request{
				sa.Media:            sa.NewStringRequest("hdd"),
				sa.ProvisioningType: sa.NewStringRequest("thick"),
				sa.TestingAttribute: sa.NewBoolRequest(true),
			},
		},
	)
	if err != nil {
		t.Fatal("Unable to add storage class: ", err)
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)
	f()
	return buf.String()
}

func TestBackendUpdateAndDelete(t *testing.T) {
	const (
		backendName       = "updateBackend"
		scName            = "updateBackendTest"
		newSCName         = "updateBackendTest2"
		volumeName        = "updateVolume"
		offlineVolumeName = "offlineVolume"
		backendProtocol   = config.File
	)
	// Test setup
	orchestrator := getOrchestrator(t, false)
	addBackendStorageClass(t, orchestrator, backendName, scName, backendProtocol)

	orchestrator.mutex.Lock()
	sc, ok := orchestrator.storageClasses[scName]
	if !ok {
		t.Fatal("Storage class not found in orchestrator map")
	}
	orchestrator.mutex.Unlock()

	_, err := orchestrator.AddVolume(ctx(), tu.GenerateVolumeConfig(volumeName, 50, scName, config.File))
	if err != nil {
		t.Fatal("Unable to create volume: ", err)
	}

	orchestrator.mutex.Lock()
	volume, ok := orchestrator.volumes[volumeName]
	if !ok {
		t.Fatalf("Volume %s not tracked by the orchestrator!", volumeName)
	}
	log.WithFields(
		log.Fields{
			"volume.BackendUUID": volume.BackendUUID,
			"volume.Config.Name": volume.Config.Name,
			"volume.Config":      volume.Config,
		},
	).Debug("Found volume.")
	startingBackend, err := orchestrator.getBackendByBackendName(backendName)
	if startingBackend == nil || err != nil {
		t.Fatalf("Backend %s not stored in orchestrator", backendName)
	}
	if _, ok = startingBackend.Volumes()[volumeName]; !ok {
		t.Fatalf("Volume %s not tracked by the backend %s!", volumeName, backendName)
	}
	orchestrator.mutex.Unlock()

	// Test updates that should succeed
	previousBackends := make([]storage.Backend, 1)
	previousBackends[0] = startingBackend
	for _, c := range []struct {
		name  string
		pools map[string]*fake.StoragePool
	}{
		{
			name: "New pool",
			pools: map[string]*fake.StoragePool{
				"primary": {
					Attrs: map[string]sa.Offer{
						sa.Media:            sa.NewStringOffer("hdd"),
						sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
						sa.TestingAttribute: sa.NewBoolOffer(true),
					},
					Bytes: 100 * 1024 * 1024 * 1024,
				},
				"secondary": {
					Attrs: map[string]sa.Offer{
						sa.Media:            sa.NewStringOffer("ssd"),
						sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
						sa.TestingAttribute: sa.NewBoolOffer(true),
					},
					Bytes: 100 * 1024 * 1024 * 1024,
				},
			},
		},
		{
			name: "Removed pool",
			pools: map[string]*fake.StoragePool{
				"primary": {
					Attrs: map[string]sa.Offer{
						sa.Media:            sa.NewStringOffer("hdd"),
						sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
						sa.TestingAttribute: sa.NewBoolOffer(true),
					},
					Bytes: 100 * 1024 * 1024 * 1024,
				},
			},
		},
		{
			name: "Expanded offer",
			pools: map[string]*fake.StoragePool{
				"primary": {
					Attrs: map[string]sa.Offer{
						sa.Media:            sa.NewStringOffer("ssd", "hdd"),
						sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
						sa.TestingAttribute: sa.NewBoolOffer(true),
					},
					Bytes: 100 * 1024 * 1024 * 1024,
				},
			},
		},
	} {
		// Make sure we're starting with an active backend
		previousBackend, err := orchestrator.getBackendByBackendName(backendName)
		if previousBackend == nil || err != nil {
			t.Fatalf("Backend %s not stored in orchestrator", backendName)
		}
		if !previousBackend.Driver().Initialized() {
			t.Errorf("Backend %s is not initialized", backendName)
		}

		var volumes []fake.Volume
		newConfigJSON, err := fakedriver.NewFakeStorageDriverConfigJSON(
			backendName,
			config.File, c.pools, volumes,
		)
		if err != nil {
			t.Errorf("%s: unable to generate new backend config: %v", c.name, err)
			continue
		}

		_, err = orchestrator.UpdateBackend(ctx(), backendName, newConfigJSON, "")

		if err != nil {
			t.Errorf("%s: unable to update backend with a nonconflicting change: %v", c.name, err)
			continue
		}

		orchestrator.mutex.Lock()
		newBackend, err := orchestrator.getBackendByBackendName(backendName)
		if newBackend == nil || err != nil {
			t.Fatalf("Backend %s not stored in orchestrator", backendName)
		}
		if previousBackend.Driver().Initialized() {
			t.Errorf("Previous backend %s still initialized", backendName)
		}
		if !newBackend.Driver().Initialized() {
			t.Errorf("Updated backend %s is not initialized.", backendName)
		}
		pools := sc.GetStoragePoolsForProtocol(ctx(), config.File, config.ReadWriteMany)
		foundNewBackend := false
		for _, pool := range pools {
			for i, b := range previousBackends {
				if pool.Backend() == b {
					t.Errorf(
						"%s: backend %d not cleared from storage class",
						c.name, i+1,
					)
				}
				if pool.Backend() == newBackend {
					foundNewBackend = true
				}
			}
		}
		if !foundNewBackend {
			t.Errorf("%s: Storage class does not point to new backend.", c.name)
		}
		matchingPool, ok := newBackend.Storage()["primary"]
		if !ok {
			t.Errorf("%s: storage pool for volume not found", c.name)
			continue
		}
		if len(matchingPool.StorageClasses()) != 1 {
			t.Errorf("%s: unexpected number of storage classes for main storage pool: %d", c.name,
				len(matchingPool.StorageClasses()))
		}
		volumeBackend, err := orchestrator.getBackendByBackendUUID(volume.BackendUUID)
		if volumeBackend == nil || err != nil {
			for backendUUID, backend := range orchestrator.backends {
				log.WithFields(
					log.Fields{
						"volume.BackendUUID":  volume.BackendUUID,
						"backend":             backend,
						"backend.BackendUUID": backend.BackendUUID(),
						"uuid":                backendUUID,
					},
				).Debug("Found backend.")
			}
			t.Fatalf("Backend %s not stored in orchestrator, err: %v", volume.BackendUUID, err)
		}
		if volumeBackend != newBackend {
			t.Errorf("%s: volume backend does not point to the new backend", c.name)
		}
		if volume.Pool != matchingPool.Name() {
			t.Errorf("%s: volume does not point to the right storage pool.", c.name)
		}
		persistentBackend, err := orchestrator.storeClient.GetBackend(ctx(), backendName)
		if err != nil {
			t.Error("Unable to retrieve backend from store client: ", err)
		} else if !cmp.Equal(newBackend.ConstructPersistent(ctx()), persistentBackend) {
			if diff := cmp.Diff(newBackend.ConstructPersistent(ctx()), persistentBackend); diff != "" {
				t.Errorf("Failure for %v; mismatch (-want +got):\n%s", c.name, diff)
			}
			t.Errorf("Backend not correctly updated in persistent store.")
		}
		previousBackends = append(previousBackends, newBackend)
		orchestrator.mutex.Unlock()
	}

	backend := previousBackends[len(previousBackends)-1]
	pool := volume.Pool

	// Test backend offlining.
	err = orchestrator.DeleteBackend(ctx(), backendName)
	if err != nil {
		t.Fatalf("Unable to delete backend: %v", err)
	}
	if !backend.Driver().Initialized() {
		t.Errorf("Deleted backend with volumes %s is not initialized.", backendName)
	}
	_, err = orchestrator.AddVolume(ctx(), tu.GenerateVolumeConfig(offlineVolumeName, 50, scName, config.File))
	if err == nil {
		t.Error("Created volume volume on offline backend.")
	}
	orchestrator.mutex.Lock()
	pools := sc.GetStoragePoolsForProtocol(ctx(), config.File, config.ReadWriteOnce)
	if len(pools) == 1 {
		t.Error("Offline backend not removed from storage pool in storage class.")
	}
	foundBackend, err := orchestrator.getBackendByBackendUUID(volume.BackendUUID)
	if err != nil {
		t.Errorf("Couldn't find backend: %v", err)
	}
	if foundBackend != backend {
		t.Error("Backend changed for volume after offlining.")
	}
	if volume.Pool != pool {
		t.Error("Storage pool changed for volume after backend offlined.")
	}
	persistentBackend, err := orchestrator.storeClient.GetBackend(ctx(), backendName)
	if err != nil {
		t.Errorf("Unable to retrieve backend from store client after offlining: %v", err)
	} else if persistentBackend.State.IsOnline() {
		t.Error("Online not set to true in the backend.")
	}
	orchestrator.mutex.Unlock()

	// Ensure that new storage classes do not get the offline backend assigned
	// to them.
	newSCExternal, err := orchestrator.AddStorageClass(
		ctx(), &storageclass.Config{
			Name: newSCName,
			Attributes: map[string]sa.Request{
				sa.Media:            sa.NewStringRequest("hdd"),
				sa.TestingAttribute: sa.NewBoolRequest(true),
			},
		},
	)
	if err != nil {
		t.Fatal("Unable to add new storage class after offlining: ", err)
	}
	if _, ok = newSCExternal.StoragePools[backendName]; ok {
		t.Error("Offline backend added to new storage class.")
	}

	// Test that online gets set properly after bootstrapping.
	newOrchestrator := getOrchestrator(t, false)
	// We need to lock the orchestrator mutex here because we call
	// ConstructExternal on the original backend in the else if clause.
	orchestrator.mutex.Lock()
	if bootstrappedBackend, _ := newOrchestrator.GetBackend(ctx(), backendName); bootstrappedBackend == nil {
		t.Error("Unable to find backend after bootstrapping.")
	} else if !reflect.DeepEqual(bootstrappedBackend, backend.ConstructExternal(ctx())) {
		diffExternalBackends(t, backend.ConstructExternal(ctx()), bootstrappedBackend)
	}

	orchestrator.mutex.Unlock()

	newOrchestrator.mutex.Lock()
	for _, name := range []string{scName, newSCName} {
		newSC, ok := newOrchestrator.storageClasses[name]
		if !ok {
			t.Fatalf(
				"Unable to find storage class %s after bootstrapping.",
				name,
			)
		}
		pools = newSC.GetStoragePoolsForProtocol(ctx(), config.File, config.ReadOnlyMany)
		if len(pools) == 1 {
			t.Errorf(
				"Offline backend readded to storage class %s after "+
					"bootstrapping.", name,
			)
		}
	}
	newOrchestrator.mutex.Unlock()

	// Test that deleting the volume causes the backend to be deleted.
	err = orchestrator.DeleteVolume(ctx(), volumeName)
	if err != nil {
		t.Fatal("Unable to delete volume for offline backend: ", err)
	}
	if backend.Driver().Initialized() {
		t.Errorf("Deleted backend %s is still initialized.", backendName)
	}
	persistentBackend, err = orchestrator.storeClient.GetBackend(ctx(), backendName)
	if err == nil {
		t.Error(
			"Backend remained on store client after deleting the last " +
				"volume present.",
		)
	}
	orchestrator.mutex.Lock()

	missingBackend, _ := orchestrator.getBackendByBackendName(backendName)
	if missingBackend != nil {
		t.Error("Empty offlined backend not removed from memory.")
	}
	orchestrator.mutex.Unlock()
	cleanup(t, orchestrator)
}

func backendPasswordsInLogsHelper(t *testing.T, debugTraceFlags map[string]bool) {
	backendName := "passwordBackend"
	backendProtocol := config.File

	orchestrator := getOrchestrator(t, false)

	fakeConfig, err := fakedriver.NewFakeStorageDriverConfigJSONWithDebugTraceFlags(
		backendName, backendProtocol,
		debugTraceFlags, "prefix1_",
	)
	if err != nil {
		t.Fatalf("Unable to generate config JSON for %s: %v", backendName, err)
	}

	_, err = orchestrator.AddBackend(ctx(), fakeConfig, "")
	if err != nil {
		t.Errorf("Unable to add backend %s: %v", backendName, err)
	}

	newConfigJSON, err := fakedriver.NewFakeStorageDriverConfigJSONWithDebugTraceFlags(
		backendName, backendProtocol,
		debugTraceFlags, "prefix2_",
	)
	if err != nil {
		t.Errorf("%s: unable to generate new backend config: %v", backendName, err)
	}

	output := captureOutput(
		func() {
			_, err = orchestrator.UpdateBackend(ctx(), backendName, newConfigJSON, "")
		},
	)

	if err != nil {
		t.Errorf("%s: unable to update backend with a nonconflicting change: %v", backendName, err)
	}

	assert.Contains(t, output, "configJSON")
	outputArr := strings.Split(output, "configJSON")
	outputArr = strings.Split(outputArr[1], "=\"")
	outputArr = strings.Split(outputArr[1], "\"")

	assert.Equal(t, outputArr[0], "<suppressed>")
	cleanup(t, orchestrator)
}

func TestBackendPasswordsInLogs(t *testing.T) {
	backendPasswordsInLogsHelper(t, nil)
	backendPasswordsInLogsHelper(t, map[string]bool{"method": true})
}

func TestEmptyBackendDeletion(t *testing.T) {
	const (
		backendName     = "emptyBackend"
		backendProtocol = config.File
	)

	orchestrator := getOrchestrator(t, false)
	// Note that we don't care about the storage class here, but it's easier
	// to reuse functionality.
	addBackendStorageClass(t, orchestrator, backendName, "none", backendProtocol)
	backend, errLookup := orchestrator.getBackendByBackendName(backendName)
	if backend == nil || errLookup != nil {
		t.Fatalf("Backend %s not stored in orchestrator", backendName)
	}

	err := orchestrator.DeleteBackend(ctx(), backendName)
	if err != nil {
		t.Fatalf("Unable to delete backend: %v", err)
	}
	if backend.Driver().Initialized() {
		t.Errorf("Deleted backend %s is still initialized.", backendName)
	}
	_, err = orchestrator.storeClient.GetBackend(ctx(), backendName)
	if err == nil {
		t.Error("Empty backend remained on store client after offlining")
	}
	orchestrator.mutex.Lock()
	missingBackend, _ := orchestrator.getBackendByBackendName(backendName)
	if missingBackend != nil {
		t.Error("Empty offlined backend not removed from memory.")
	}
	orchestrator.mutex.Unlock()
	cleanup(t, orchestrator)
}

func TestBootstrapSnapshotMissingVolume(t *testing.T) {
	const (
		offlineBackendName = "snapNoVolBackend"
		scName             = "snapNoVolSC"
		volumeName         = "snapNoVolVolume"
		snapName           = "snapNoVolSnapshot"
		backendProtocol    = config.File
	)

	orchestrator := getOrchestrator(t, false)
	defer cleanup(t, orchestrator)
	addBackendStorageClass(t, orchestrator, offlineBackendName, scName, backendProtocol)
	_, err := orchestrator.AddVolume(
		ctx(), tu.GenerateVolumeConfig(
			volumeName, 50,
			scName, config.File,
		),
	)
	if err != nil {
		t.Fatal("Unable to create volume: ", err)
	}

	// For the full test, we create everything and recreate the AddSnapshot transaction.
	snapshotConfig := generateSnapshotConfig(snapName, volumeName, volumeName)
	if _, err := orchestrator.CreateSnapshot(ctx(), snapshotConfig); err != nil {
		t.Fatal("Unable to add snapshot: ", err)
	}

	// Simulate deleting the existing volume without going through Trident then bootstrapping
	vol, ok := orchestrator.volumes[volumeName]
	if !ok {
		t.Fatalf("Unable to find volume %s in backend.", volumeName)
	}
	orchestrator.mutex.Lock()
	err = orchestrator.storeClient.DeleteVolume(ctx(), vol)
	if err != nil {
		t.Fatalf("Unable to delete volume from store: %v", err)
	}
	orchestrator.mutex.Unlock()

	newOrchestrator := getOrchestrator(t, false)
	bootstrappedSnapshot, err := newOrchestrator.GetSnapshot(ctx(), snapshotConfig.VolumeName, snapshotConfig.Name)
	if err != nil {
		t.Fatalf("error getting snapshot: %v", err)
	}
	if bootstrappedSnapshot == nil {
		t.Error("Volume not found during bootstrap.")
	}
	if !bootstrappedSnapshot.State.IsMissingVolume() {
		t.Error("Unexpected snapshot state.")
	}
	// Delete volume in missing_volume state
	err = newOrchestrator.DeleteSnapshot(ctx(), volumeName, snapName)
	if err != nil {
		t.Error("could not delete snapshot with missing volume")
	}
}

func TestBootstrapSnapshotMissingBackend(t *testing.T) {
	const (
		offlineBackendName = "snapNoBackBackend"
		scName             = "snapNoBackSC"
		volumeName         = "snapNoBackVolume"
		snapName           = "snapNoBackSnapshot"
		backendProtocol    = config.File
	)

	orchestrator := getOrchestrator(t, false)
	defer cleanup(t, orchestrator)
	addBackendStorageClass(t, orchestrator, offlineBackendName, scName, backendProtocol)
	_, err := orchestrator.AddVolume(
		ctx(), tu.GenerateVolumeConfig(
			volumeName, 50,
			scName, config.File,
		),
	)
	if err != nil {
		t.Fatal("Unable to create volume: ", err)
	}

	// For the full test, we create everything and recreate the AddSnapshot transaction.
	snapshotConfig := generateSnapshotConfig(snapName, volumeName, volumeName)
	if _, err := orchestrator.CreateSnapshot(ctx(), snapshotConfig); err != nil {
		t.Fatal("Unable to add snapshot: ", err)
	}

	// Simulate deleting the existing backend without going through Trident then bootstrapping
	backend, err := orchestrator.getBackendByBackendName(offlineBackendName)
	if err != nil {
		t.Fatalf("Unable to get backend from store: %v", err)
	}
	orchestrator.mutex.Lock()
	err = orchestrator.storeClient.DeleteBackend(ctx(), backend)
	if err != nil {
		t.Fatalf("Unable to delete volume from store: %v", err)
	}
	orchestrator.mutex.Unlock()

	newOrchestrator := getOrchestrator(t, false)
	bootstrappedSnapshot, err := newOrchestrator.GetSnapshot(ctx(), snapshotConfig.VolumeName, snapshotConfig.Name)
	if err != nil {
		t.Fatalf("error getting snapshot: %v", err)
	}
	if bootstrappedSnapshot == nil {
		t.Error("Volume not found during bootstrap.")
	}
	if !bootstrappedSnapshot.State.IsMissingBackend() {
		t.Error("Unexpected snapshot state.")
	}
	// Delete snapshot in missing_backend state
	err = newOrchestrator.DeleteSnapshot(ctx(), volumeName, snapName)
	if err != nil {
		t.Error("could not delete snapshot with missing backend")
	}
}

func TestBootstrapVolumeMissingBackend(t *testing.T) {
	const (
		offlineBackendName = "bootstrapVolBackend"
		scName             = "bootstrapVolSC"
		volumeName         = "bootstrapVolVolume"
		backendProtocol    = config.File
	)

	orchestrator := getOrchestrator(t, false)
	defer cleanup(t, orchestrator)
	addBackendStorageClass(t, orchestrator, offlineBackendName, scName, backendProtocol)
	_, err := orchestrator.AddVolume(
		ctx(), tu.GenerateVolumeConfig(
			volumeName, 50,
			scName, config.File,
		),
	)
	if err != nil {
		t.Fatal("Unable to create volume: ", err)
	}

	// Simulate deleting the existing backend without going through Trident then bootstrapping
	backend, err := orchestrator.getBackendByBackendName(offlineBackendName)
	if err != nil {
		t.Fatalf("Unable to get backend from store: %v", err)
	}
	orchestrator.mutex.Lock()
	err = orchestrator.storeClient.DeleteBackend(ctx(), backend)
	if err != nil {
		t.Fatalf("Unable to delete volume from store: %v", err)
	}
	orchestrator.mutex.Unlock()

	newOrchestrator := getOrchestrator(t, false)
	bootstrappedVolume, err := newOrchestrator.GetVolume(ctx(), volumeName)
	if err != nil {
		t.Fatalf("error getting volume: %v", err)
	}
	if bootstrappedVolume == nil {
		t.Error("volume not found during bootstrap")
	}
	if !bootstrappedVolume.State.IsMissingBackend() {
		t.Error("unexpected volume state")
	}

	// Delete volume in missing_backend state
	err = newOrchestrator.DeleteVolume(ctx(), volumeName)
	if err != nil {
		t.Error("could not delete volume with missing backend")
	}
}

func TestBackendCleanup(t *testing.T) {
	const (
		offlineBackendName = "cleanupBackend"
		onlineBackendName  = "onlineBackend"
		scName             = "cleanupBackendTest"
		volumeName         = "cleanupVolume"
		backendProtocol    = config.File
	)

	orchestrator := getOrchestrator(t, false)
	addBackendStorageClass(t, orchestrator, offlineBackendName, scName, backendProtocol)
	_, err := orchestrator.AddVolume(
		ctx(), tu.GenerateVolumeConfig(
			volumeName, 50,
			scName, config.File,
		),
	)
	if err != nil {
		t.Fatal("Unable to create volume: ", err)
	}

	// This needs to go after the volume addition to ensure that the volume
	// ends up on the backend to be offlined.
	addBackend(t, orchestrator, onlineBackendName, backendProtocol)

	err = orchestrator.DeleteBackend(ctx(), offlineBackendName)
	if err != nil {
		t.Fatalf("Unable to delete backend %s: %v", offlineBackendName, err)
	}
	// Simulate deleting the existing volume and then bootstrapping
	orchestrator.mutex.Lock()
	vol, ok := orchestrator.volumes[volumeName]
	if !ok {
		t.Fatalf("Unable to find volume %s in backend.", volumeName)
	}
	err = orchestrator.storeClient.DeleteVolume(ctx(), vol)
	if err != nil {
		t.Fatalf("Unable to delete volume from store: %v", err)
	}
	orchestrator.mutex.Unlock()

	newOrchestrator := getOrchestrator(t, false)
	if bootstrappedBackend, _ := newOrchestrator.GetBackend(ctx(), offlineBackendName); bootstrappedBackend != nil {
		t.Error("Empty offline backend not deleted during bootstrap.")
	}
	if bootstrappedBackend, _ := newOrchestrator.GetBackend(ctx(), onlineBackendName); bootstrappedBackend == nil {
		t.Error("Empty online backend deleted during bootstrap.")
	}
}

func TestLoadBackend(t *testing.T) {
	const (
		backendName = "load-backend-test"
	)
	// volumes must be nil in order to satisfy reflect.DeepEqual comparison. It isn't recommended to compare slices with deepEqual
	var volumes []fake.Volume
	orchestrator := getOrchestrator(t, false)
	configJSON, err := fakedriver.NewFakeStorageDriverConfigJSON(
		backendName,
		config.File,
		map[string]*fake.StoragePool{
			"primary": {
				Attrs: map[string]sa.Offer{
					sa.Media:            sa.NewStringOffer("hdd"),
					sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
					sa.TestingAttribute: sa.NewBoolOffer(true),
				},
				Bytes: 100 * 1024 * 1024 * 1024,
			},
		},
		volumes,
	)
	originalBackend, err := orchestrator.AddBackend(ctx(), configJSON, "")
	if err != nil {
		t.Fatal("Unable to initially add backend: ", err)
	}
	persistentBackend, err := orchestrator.storeClient.GetBackend(ctx(), backendName)
	if err != nil {
		t.Fatal("Unable to retrieve backend from store client: ", err)
	}
	// Note that this will register as an update, but it should be close enough
	newConfig, err := persistentBackend.MarshalConfig()
	if err != nil {
		t.Fatal("Unable to marshal config from stored backend: ", err)
	}
	newBackend, err := orchestrator.AddBackend(ctx(), newConfig, "")
	if err != nil {
		t.Error("Unable to update backend from config: ", err)
	} else if !reflect.DeepEqual(newBackend, originalBackend) {
		t.Error("Newly loaded backend differs.")
	}

	newOrchestrator := getOrchestrator(t, false)
	if bootstrappedBackend, _ := newOrchestrator.GetBackend(ctx(), backendName); bootstrappedBackend == nil {
		t.Error("Unable to find backend after bootstrapping.")
	} else if !reflect.DeepEqual(bootstrappedBackend, originalBackend) {
		t.Errorf("External backends differ.")
		diffExternalBackends(t, originalBackend, bootstrappedBackend)
	}
	cleanup(t, orchestrator)
}

func prepRecoveryTest(
	t *testing.T, orchestrator *TridentOrchestrator, backendName, scName string,
) {
	configJSON, err := fakedriver.NewFakeStorageDriverConfigJSON(
		backendName,
		config.File,
		map[string]*fake.StoragePool{
			"primary": {
				Attrs: map[string]sa.Offer{
					sa.Media:            sa.NewStringOffer("hdd"),
					sa.ProvisioningType: sa.NewStringOffer("thick", "thin"),
					sa.RecoveryTest:     sa.NewBoolOffer(true),
				},
				Bytes: 100 * 1024 * 1024 * 1024,
			},
		},
		[]fake.Volume{},
	)
	_, err = orchestrator.AddBackend(ctx(), configJSON, "")
	if err != nil {
		t.Fatal("Unable to initialize backend: ", err)
	}
	_, err = orchestrator.AddStorageClass(
		ctx(), &storageclass.Config{
			Name: scName,
			Attributes: map[string]sa.Request{
				sa.Media:            sa.NewStringRequest("hdd"),
				sa.ProvisioningType: sa.NewStringRequest("thick"),
				sa.RecoveryTest:     sa.NewBoolRequest(true),
			},
		},
	)
	if err != nil {
		t.Fatal("Unable to add storage class: ", err)
	}
}

func runRecoveryTests(
	t *testing.T,
	orchestrator *TridentOrchestrator,
	backendName string,
	op storage.VolumeOperation,
	testCases []recoveryTest,
) {
	for _, c := range testCases {
		// Manipulate the persistent store directly, since it's
		// easier to store the results of a partially completed volume addition
		// than to actually inject a failure.
		volTxn := &storage.VolumeTransaction{
			Config: c.volumeConfig,
			Op:     op,
		}
		err := orchestrator.storeClient.AddVolumeTransaction(ctx(), volTxn)
		if err != nil {
			t.Fatalf("%s: Unable to create volume transaction: %v", c.name, err)
		}
		newOrchestrator := getOrchestrator(t, false)
		newOrchestrator.mutex.Lock()
		if _, ok := newOrchestrator.volumes[c.volumeConfig.Name]; ok {
			t.Errorf("%s: volume still present in orchestrator.", c.name)
			// Note: assume that if the volume's still present in the
			// top-level map, it's present everywhere else and that, if it's
			// absent there, it's absent everywhere else in memory
		}
		backend, err := newOrchestrator.getBackendByBackendName(backendName)
		if backend == nil || err != nil {
			t.Fatalf("%s: Backend not found after bootstrapping.", c.name)
		}
		f, ok := backend.Driver().(*fakedriver.StorageDriver)
		if !ok {
			t.Fatalf("%e", utils.TypeAssertionError("backend.Driver().(*fakedriver.StorageDriver)"))
		}
		// Destroy should be always called on the backend
		if _, ok := f.DestroyedVolumes[f.GetInternalVolumeName(ctx(), c.volumeConfig.Name)]; !ok && c.expectDestroy {
			t.Errorf("%s: Destroy not called on volume.", c.name)
		}
		_, err = newOrchestrator.storeClient.GetVolume(ctx(), c.volumeConfig.Name)
		if err != nil {
			if !persistentstore.MatchKeyNotFoundErr(err) {
				t.Errorf("%s: unable to communicate with backing store: %v", c.name, err)
			}
		} else {
			t.Errorf("%s: Found VolumeConfig still stored in store.", c.name)
		}
		if txns, err := newOrchestrator.storeClient.GetVolumeTransactions(ctx()); err != nil {
			t.Errorf("%s: Unable to retrieve transactions from backing store: %v", c.name, err)
		} else if len(txns) > 0 {
			t.Errorf("%s: Transaction not cleared from the backing store", c.name)
		}
		newOrchestrator.mutex.Unlock()
	}
}

func TestAddVolumeRecovery(t *testing.T) {
	const (
		backendName      = "addRecoveryBackend"
		scName           = "addRecoveryBackendSC"
		fullVolumeName   = "addRecoveryVolumeFull"
		txOnlyVolumeName = "addRecoveryVolumeTxOnly"
	)
	orchestrator := getOrchestrator(t, false)
	prepRecoveryTest(t, orchestrator, backendName, scName)
	// It's easier to add the volume and then reinject the transaction begin
	// afterwards
	fullVolumeConfig := tu.GenerateVolumeConfig(fullVolumeName, 50, scName, config.File)
	_, err := orchestrator.AddVolume(ctx(), fullVolumeConfig)
	if err != nil {
		t.Fatal("Unable to add volume: ", err)
	}
	txOnlyVolumeConfig := tu.GenerateVolumeConfig(txOnlyVolumeName, 50, scName, config.File)
	// BEGIN actual test
	runRecoveryTests(
		t, orchestrator, backendName, storage.AddVolume,
		[]recoveryTest{
			{name: "full", volumeConfig: fullVolumeConfig, expectDestroy: true},
			{name: "txOnly", volumeConfig: txOnlyVolumeConfig, expectDestroy: true},
		},
	)
	cleanup(t, orchestrator)
}

func TestAddVolumeWithTMRNonONTAPNAS(t *testing.T) {
	// Add a single backend of fake
	// create volume with relationship annotation added
	// witness failure
	const (
		backendName    = "addRecoveryBackend"
		scName         = "addRecoveryBackendSC"
		fullVolumeName = "addRecoveryVolumeFull"
	)
	orchestrator := getOrchestrator(t, false)
	prepRecoveryTest(t, orchestrator, backendName, scName)
	// It's easier to add the volume and then reinject the transaction begin
	// afterwards
	fullVolumeConfig := tu.GenerateVolumeConfig(
		fullVolumeName, 50, scName,
		config.File,
	)
	fullVolumeConfig.PeerVolumeHandle = "fakesvm:fakevolume"
	fullVolumeConfig.IsMirrorDestination = true
	_, err := orchestrator.AddVolume(ctx(), fullVolumeConfig)
	if err == nil || !strings.Contains(err.Error(), "no suitable") {
		t.Fatal("Unexpected failure")
	}
	cleanup(t, orchestrator)
}

func TestDeleteVolumeRecovery(t *testing.T) {
	const (
		backendName      = "deleteRecoveryBackend"
		scName           = "deleteRecoveryBackendSC"
		fullVolumeName   = "deleteRecoveryVolumeFull"
		txOnlyVolumeName = "deleteRecoveryVolumeTxOnly"
	)
	orchestrator := getOrchestrator(t, false)
	prepRecoveryTest(t, orchestrator, backendName, scName)

	// For the full test, we delete everything but the ending transaction.
	fullVolumeConfig := tu.GenerateVolumeConfig(fullVolumeName, 50, scName, config.File)
	if _, err := orchestrator.AddVolume(ctx(), fullVolumeConfig); err != nil {
		t.Fatal("Unable to add volume: ", err)
	}
	if err := orchestrator.DeleteVolume(ctx(), fullVolumeName); err != nil {
		t.Fatal("Unable to remove full volume: ", err)
	}

	txOnlyVolumeConfig := tu.GenerateVolumeConfig(txOnlyVolumeName, 50, scName, config.File)
	if _, err := orchestrator.AddVolume(ctx(), txOnlyVolumeConfig); err != nil {
		t.Fatal("Unable to add tx only volume: ", err)
	}

	// BEGIN actual test
	runRecoveryTests(
		t, orchestrator, backendName,
		storage.DeleteVolume, []recoveryTest{
			{name: "full", volumeConfig: fullVolumeConfig, expectDestroy: false},
			{name: "txOnly", volumeConfig: txOnlyVolumeConfig, expectDestroy: true},
		},
	)
	cleanup(t, orchestrator)
}

func generateSnapshotConfig(
	name, volumeName, volumeInternalName string,
) *storage.SnapshotConfig {
	return &storage.SnapshotConfig{
		Version:            config.OrchestratorAPIVersion,
		Name:               name,
		VolumeName:         volumeName,
		VolumeInternalName: volumeInternalName,
	}
}

func runSnapshotRecoveryTests(
	t *testing.T,
	orchestrator *TridentOrchestrator,
	backendName string,
	op storage.VolumeOperation,
	testCases []recoveryTest,
) {
	for _, c := range testCases {
		// Manipulate the persistent store directly, since it's
		// easier to store the results of a partially completed snapshot addition
		// than to actually inject a failure.
		volTxn := &storage.VolumeTransaction{
			Config:         c.volumeConfig,
			SnapshotConfig: c.snapshotConfig,
			Op:             op,
		}
		if err := orchestrator.storeClient.AddVolumeTransaction(ctx(), volTxn); err != nil {
			t.Fatalf("%s: Unable to create volume transaction: %v", c.name, err)
		}
		newOrchestrator := getOrchestrator(t, false)
		newOrchestrator.mutex.Lock()
		if _, ok := newOrchestrator.snapshots[c.snapshotConfig.ID()]; ok {
			t.Errorf("%s: snapshot still present in orchestrator.", c.name)
			// Note: assume that if the snapshot's still present in the
			// top-level map, it's present everywhere else and that, if it's
			// absent there, it's absent everywhere else in memory
		}
		backend, err := newOrchestrator.getBackendByBackendName(backendName)
		if err != nil {
			t.Fatalf("%s: Backend not found after bootstrapping.", c.name)
		}
		f, ok := backend.Driver().(*fakedriver.StorageDriver)
		if !ok {
			t.Fatalf("%e", utils.TypeAssertionError("backend.Driver().(*fakedriver.StorageDriver)"))
		}

		_, ok = f.DestroyedSnapshots[c.snapshotConfig.ID()]
		if !ok && c.expectDestroy {
			t.Errorf("%s: Destroy not called on snapshot.", c.name)
		} else if ok && !c.expectDestroy {
			t.Errorf("%s: Destroy should not have been called on snapshot.", c.name)
		}

		_, err = newOrchestrator.storeClient.GetSnapshot(ctx(), c.snapshotConfig.VolumeName, c.snapshotConfig.Name)
		if err != nil {
			if !persistentstore.MatchKeyNotFoundErr(err) {
				t.Errorf("%s: unable to communicate with backing store: %v", c.name, err)
			}
		} else {
			t.Errorf("%s: Found SnapshotConfig still stored in store.", c.name)
		}
		if txns, err := newOrchestrator.storeClient.GetVolumeTransactions(ctx()); err != nil {
			t.Errorf("%s: Unable to retrieve transactions from backing store: %v", c.name, err)
		} else if len(txns) > 0 {
			t.Errorf("%s: Transaction not cleared from the backing store", c.name)
		}
		newOrchestrator.mutex.Unlock()
	}
}

func TestAddSnapshotRecovery(t *testing.T) {
	const (
		backendName        = "addSnapshotRecoveryBackend"
		scName             = "addSnapshotRecoveryBackendSC"
		volumeName         = "addSnapshotRecoveryVolume"
		fullSnapshotName   = "addSnapshotRecoverySnapshotFull"
		txOnlySnapshotName = "addSnapshotRecoverySnapshotTxOnly"
	)
	orchestrator := getOrchestrator(t, false)
	prepRecoveryTest(t, orchestrator, backendName, scName)

	// It's easier to add the volume/snapshot and then reinject the transaction again afterwards.
	volumeConfig := tu.GenerateVolumeConfig(volumeName, 50, scName, config.File)
	if _, err := orchestrator.AddVolume(ctx(), volumeConfig); err != nil {
		t.Fatal("Unable to add volume: ", err)
	}

	// For the full test, we create everything and recreate the AddSnapshot transaction.
	fullSnapshotConfig := generateSnapshotConfig(fullSnapshotName, volumeName, volumeName)
	if _, err := orchestrator.CreateSnapshot(ctx(), fullSnapshotConfig); err != nil {
		t.Fatal("Unable to add snapshot: ", err)
	}

	// For the partial test, we add only the AddSnapshot transaction.
	txOnlySnapshotConfig := generateSnapshotConfig(txOnlySnapshotName, volumeName, volumeName)

	// BEGIN actual test.  Note that the delete idempotency is handled at the backend layer
	// (above the driver), so if the snapshot doesn't exist after bootstrapping, the driver
	// will not be called to delete the snapshot.
	runSnapshotRecoveryTests(
		t, orchestrator, backendName, storage.AddSnapshot,
		[]recoveryTest{
			{name: "full", volumeConfig: volumeConfig, snapshotConfig: fullSnapshotConfig, expectDestroy: true},
			{name: "txOnly", volumeConfig: volumeConfig, snapshotConfig: txOnlySnapshotConfig, expectDestroy: false},
		},
	)
	cleanup(t, orchestrator)
}

func TestDeleteSnapshotRecovery(t *testing.T) {
	const (
		backendName        = "deleteSnapshotRecoveryBackend"
		scName             = "deleteSnapshotRecoveryBackendSC"
		volumeName         = "deleteSnapshotRecoveryVolume"
		fullSnapshotName   = "deleteSnapshotRecoverySnapshotFull"
		txOnlySnapshotName = "deleteSnapshotRecoverySnapshotTxOnly"
	)
	orchestrator := getOrchestrator(t, false)
	prepRecoveryTest(t, orchestrator, backendName, scName)

	// For the full test, we delete everything and recreate the delete transaction.
	volumeConfig := tu.GenerateVolumeConfig(volumeName, 50, scName, config.File)
	if _, err := orchestrator.AddVolume(ctx(), volumeConfig); err != nil {
		t.Fatal("Unable to add volume: ", err)
	}
	fullSnapshotConfig := generateSnapshotConfig(fullSnapshotName, volumeName, volumeName)
	if _, err := orchestrator.CreateSnapshot(ctx(), fullSnapshotConfig); err != nil {
		t.Fatal("Unable to add snapshot: ", err)
	}
	if err := orchestrator.DeleteSnapshot(ctx(), volumeName, fullSnapshotName); err != nil {
		t.Fatal("Unable to remove full snapshot: ", err)
	}

	// For the partial test, we ensure the snapshot will be restored during bootstrapping,
	// and the delete transaction will ensure everything is deleted.
	txOnlySnapshotConfig := generateSnapshotConfig(txOnlySnapshotName, volumeName, volumeName)
	if _, err := orchestrator.CreateSnapshot(ctx(), txOnlySnapshotConfig); err != nil {
		t.Fatal("Unable to add snapshot: ", err)
	}

	// BEGIN actual test.  Note that the delete idempotency is handled at the backend layer
	// (above the driver), so if the snapshot doesn't exist after bootstrapping, the driver
	// will not be called to delete the snapshot.
	runSnapshotRecoveryTests(
		t, orchestrator, backendName, storage.DeleteSnapshot,
		[]recoveryTest{
			{name: "full", snapshotConfig: fullSnapshotConfig, expectDestroy: false},
			{name: "txOnly", snapshotConfig: txOnlySnapshotConfig, expectDestroy: true},
		},
	)
	cleanup(t, orchestrator)
}

// The next series of tests test that bootstrap doesn't exit early if it
// encounters a key error for one of the main types of entries.
func TestStorageClassOnlyBootstrap(t *testing.T) {
	const scName = "storageclass-only"

	orchestrator := getOrchestrator(t, false)
	originalSC, err := orchestrator.AddStorageClass(
		ctx(), &storageclass.Config{
			Name: scName,
			Attributes: map[string]sa.Request{
				sa.Media:            sa.NewStringRequest("hdd"),
				sa.ProvisioningType: sa.NewStringRequest("thick"),
				sa.RecoveryTest:     sa.NewBoolRequest(true),
			},
		},
	)
	if err != nil {
		t.Fatal("Unable to add storage class: ", err)
	}
	newOrchestrator := getOrchestrator(t, false)
	bootstrappedSC, err := newOrchestrator.GetStorageClass(ctx(), scName)
	if bootstrappedSC == nil || err != nil {
		t.Error("Unable to find storage class after bootstrapping.")
	} else if !reflect.DeepEqual(bootstrappedSC, originalSC) {
		t.Errorf("External storage classs differ:\n\tOriginal: %v\n\tBootstrapped: %v", originalSC, bootstrappedSC)
	}
	cleanup(t, orchestrator)
}

func TestFirstVolumeRecovery(t *testing.T) {
	const (
		backendName      = "firstRecoveryBackend"
		scName           = "firstRecoveryBackendSC"
		txOnlyVolumeName = "firstRecoveryVolumeTxOnly"
	)
	orchestrator := getOrchestrator(t, false)
	prepRecoveryTest(t, orchestrator, backendName, scName)
	txOnlyVolumeConfig := tu.GenerateVolumeConfig(txOnlyVolumeName, 50, scName, config.File)
	// BEGIN actual test
	runRecoveryTests(
		t, orchestrator, backendName, storage.AddVolume, []recoveryTest{
			{
				name: "firstTXOnly", volumeConfig: txOnlyVolumeConfig,
				expectDestroy: true,
			},
		},
	)
	cleanup(t, orchestrator)
}

func TestOrchestratorNotReady(t *testing.T) {
	var (
		err            error
		backend        *storage.BackendExternal
		backends       []*storage.BackendExternal
		volume         *storage.VolumeExternal
		volumes        []*storage.VolumeExternal
		snapshot       *storage.SnapshotExternal
		snapshots      []*storage.SnapshotExternal
		storageClass   *storageclass.External
		storageClasses []*storageclass.External
	)

	orchestrator := getOrchestrator(t, false)
	orchestrator.bootstrapped = false
	orchestrator.bootstrapError = utils.NotReadyError()

	backend, err = orchestrator.AddBackend(ctx(), "", "")
	if backend != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected AddBackend to return an error.")
	}

	backend, err = orchestrator.GetBackend(ctx(), "")
	if backend != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected GetBackend to return an error.")
	}

	backends, err = orchestrator.ListBackends(ctx())
	if backends != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected ListBackends to return an error.")
	}

	err = orchestrator.DeleteBackend(ctx(), "")
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected DeleteBackend to return an error.")
	}

	volume, err = orchestrator.AddVolume(ctx(), nil)
	if volume != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected AddVolume to return an error.")
	}

	volume, err = orchestrator.CloneVolume(ctx(), nil)
	if volume != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected CloneVolume to return an error.")
	}

	volume, err = orchestrator.GetVolume(ctx(), "")
	if volume != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected GetVolume to return an error.")
	}

	volumes, err = orchestrator.ListVolumes(ctx())
	if volumes != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected ListVolumes to return an error.")
	}

	err = orchestrator.DeleteVolume(ctx(), "")
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected DeleteVolume to return an error.")
	}

	err = orchestrator.AttachVolume(ctx(), "", "", nil)
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected AttachVolume to return an error.")
	}

	err = orchestrator.DetachVolume(ctx(), "", "")
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected DetachVolume to return an error.")
	}

	snapshot, err = orchestrator.CreateSnapshot(ctx(), nil)
	if snapshot != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected CreateSnapshot to return an error.")
	}

	snapshot, err = orchestrator.GetSnapshot(ctx(), "", "")
	if snapshot != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected GetSnapshot to return an error.")
	}

	snapshots, err = orchestrator.ListSnapshots(ctx())
	if snapshots != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected ListSnapshots to return an error.")
	}

	snapshots, err = orchestrator.ReadSnapshotsForVolume(ctx(), "")
	if snapshots != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected ReadSnapshotsForVolume to return an error.")
	}

	err = orchestrator.DeleteSnapshot(ctx(), "", "")
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected DeleteSnapshot to return an error.")
	}

	err = orchestrator.ReloadVolumes(ctx())
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected ReloadVolumes to return an error.")
	}

	storageClass, err = orchestrator.AddStorageClass(ctx(), nil)
	if storageClass != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected AddStorageClass to return an error.")
	}

	storageClass, err = orchestrator.GetStorageClass(ctx(), "")
	if storageClass != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected GetStorageClass to return an error.")
	}

	storageClasses, err = orchestrator.ListStorageClasses(ctx())
	if storageClasses != nil || !utils.IsNotReadyError(err) {
		t.Errorf("Expected ListStorageClasses to return an error.")
	}

	err = orchestrator.DeleteStorageClass(ctx(), "")
	if !utils.IsNotReadyError(err) {
		t.Errorf("Expected DeleteStorageClass to return an error.")
	}
}

func importVolumeSetup(
	t *testing.T, backendName, scName, volumeName, importOriginalName string,
	backendProtocol config.Protocol,
) (*TridentOrchestrator, *storage.VolumeConfig) {
	// Object setup
	orchestrator := getOrchestrator(t, false)
	addBackendStorageClass(t, orchestrator, backendName, scName, backendProtocol)

	orchestrator.mutex.Lock()
	_, ok := orchestrator.storageClasses[scName]
	if !ok {
		t.Fatal("Storageclass not found in orchestrator map")
	}
	orchestrator.mutex.Unlock()

	backendUUID := ""
	for _, b := range orchestrator.backends {
		if b.Name() == backendName {
			backendUUID = b.BackendUUID()
			break
		}
	}

	if backendUUID == "" {
		t.Fatal("BackendUUID not found")
	}

	volumeConfig := tu.GenerateVolumeConfig(volumeName, 50, scName, backendProtocol)
	volumeConfig.ImportOriginalName = importOriginalName
	volumeConfig.ImportBackendUUID = backendUUID
	return orchestrator, volumeConfig
}

func TestImportVolumeFailures(t *testing.T) {
	const (
		backendName     = "backend82"
		scName          = "sc01"
		volumeName      = "volume82"
		originalName01  = "origVolume01"
		backendProtocol = config.File
	)

	createPVandPVCError := func(volExternal *storage.VolumeExternal, driverType string) error {
		return fmt.Errorf("failed to create PV")
	}

	orchestrator, volumeConfig := importVolumeSetup(t, backendName, scName, volumeName, originalName01, backendProtocol)

	_, err := orchestrator.LegacyImportVolume(ctx(), volumeConfig, backendName, false, createPVandPVCError)

	// verify that importVolumeCleanup renamed volume to originalName
	backend, _ := orchestrator.getBackendByBackendName(backendName)
	volExternal, err := backend.Driver().GetVolumeExternal(ctx(), originalName01)
	if err != nil {
		t.Fatalf("failed to get volumeExternal for %s", originalName01)
	}
	if volExternal.Config.Size != "1000000000" {
		t.Errorf("falied to verify %s size %s", originalName01, volExternal.Config.Size)
	}

	// verify that we cleaned up the persisted state
	if _, ok := orchestrator.volumes[volumeConfig.Name]; ok {
		t.Errorf("volume %s should not exist in orchestrator's volume cache", volumeConfig.Name)
	}
	persistedVolume, err := orchestrator.storeClient.GetVolume(ctx(), volumeConfig.Name)
	if persistedVolume != nil {
		t.Errorf("volume %s should not be persisted", volumeConfig.Name)
	}

	cleanup(t, orchestrator)
}

func TestLegacyImportVolume(t *testing.T) {
	const (
		backendName     = "backend02"
		scName          = "sc01"
		volumeName01    = "volume01"
		volumeName02    = "volume02"
		originalName01  = "origVolume01"
		originalName02  = "origVolume02"
		backendProtocol = config.File
	)

	createPVandPVCNoOp := func(volExternal *storage.VolumeExternal, driverType string) error {
		return nil
	}

	orchestrator, volumeConfig := importVolumeSetup(t, backendName, scName, volumeName01, originalName01,
		backendProtocol)

	// The volume exists on the backend with the original name.
	// Set volumeConfig.InternalName to the expected volumeName post import.
	volumeConfig.InternalName = volumeName01

	notManagedVolConfig := volumeConfig.ConstructClone()
	notManagedVolConfig.Name = volumeName02
	notManagedVolConfig.InternalName = volumeName02
	notManagedVolConfig.ImportOriginalName = originalName02
	notManagedVolConfig.ImportNotManaged = true

	// Test configuration
	for _, c := range []struct {
		name                 string
		volumeConfig         *storage.VolumeConfig
		notManaged           bool
		createFunc           VolumeCallback
		expectedInternalName string
	}{
		{
			name:                 "managed",
			volumeConfig:         volumeConfig,
			notManaged:           false,
			createFunc:           createPVandPVCNoOp,
			expectedInternalName: volumeConfig.InternalName,
		},
		{
			name:                 "notManaged",
			volumeConfig:         notManagedVolConfig,
			notManaged:           true,
			createFunc:           createPVandPVCNoOp,
			expectedInternalName: originalName02,
		},
	} {
		// The test code
		volExternal, err := orchestrator.LegacyImportVolume(ctx(), c.volumeConfig, backendName, c.notManaged,
			c.createFunc)
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
		} else {
			if volExternal.Config.InternalName != c.expectedInternalName {
				t.Errorf(
					"%s: expected matching internal names %s - %s",
					c.name, c.expectedInternalName, volExternal.Config.InternalName,
				)
			}
			if _, ok := orchestrator.volumes[volExternal.Config.Name]; ok {
				if c.notManaged {
					t.Errorf("%s: notManaged volume %s should not be persisted", c.name, volExternal.Config.Name)
				}
			} else if !c.notManaged {
				t.Errorf("%s: managed volume %s should be persisted", c.name, volExternal.Config.Name)
			}

		}

	}
	cleanup(t, orchestrator)
}

func TestImportVolume(t *testing.T) {
	const (
		backendName     = "backend02"
		scName          = "sc01"
		volumeName01    = "volume01"
		volumeName02    = "volume02"
		originalName01  = "origVolume01"
		originalName02  = "origVolume02"
		backendProtocol = config.File
	)

	orchestrator, volumeConfig := importVolumeSetup(
		t, backendName, scName, volumeName01, originalName01, backendProtocol,
	)

	notManagedVolConfig := volumeConfig.ConstructClone()
	notManagedVolConfig.Name = volumeName02
	notManagedVolConfig.ImportOriginalName = originalName02
	notManagedVolConfig.ImportNotManaged = true

	// Test configuration
	for _, c := range []struct {
		name                 string
		volumeConfig         *storage.VolumeConfig
		expectedInternalName string
	}{
		{
			name:                 "managed",
			volumeConfig:         volumeConfig,
			expectedInternalName: volumeName01,
		},
		{
			name:                 "notManaged",
			volumeConfig:         notManagedVolConfig,
			expectedInternalName: originalName02,
		},
	} {
		// The test code
		volExternal, err := orchestrator.ImportVolume(ctx(), c.volumeConfig)
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
		} else {
			if volExternal.Config.InternalName != c.expectedInternalName {
				t.Errorf(
					"%s: expected matching internal names %s - %s",
					c.name, c.expectedInternalName, volExternal.Config.InternalName,
				)
			}
			if _, ok := orchestrator.volumes[volExternal.Config.Name]; !ok {
				t.Errorf("%s: managed volume %s should be persisted", c.name, volExternal.Config.Name)
			}
		}

	}
	cleanup(t, orchestrator)
}

func TestValidateImportVolumeNasBackend(t *testing.T) {
	const (
		backendName     = "backend01"
		scName          = "sc01"
		volumeName      = "volume01"
		originalName    = "origVolume01"
		backendProtocol = config.File
	)

	orchestrator, volumeConfig := importVolumeSetup(t, backendName, scName, volumeName, originalName, backendProtocol)

	_, err := orchestrator.AddVolume(ctx(), volumeConfig)
	if err != nil {
		t.Fatal("Unable to add volume: ", err)
	}

	// The volume exists on the backend with the original name since we added it above.
	// Set volumeConfig.InternalName to the expected volumeName post import.
	volumeConfig.InternalName = volumeName

	// Create VolumeConfig objects for the remaining error conditions
	pvExistsVolConfig := volumeConfig.ConstructClone()
	pvExistsVolConfig.ImportOriginalName = volumeName

	unknownSCVolConfig := volumeConfig.ConstructClone()
	unknownSCVolConfig.StorageClass = "sc02"

	missingVolConfig := volumeConfig.ConstructClone()
	missingVolConfig.ImportOriginalName = "noVol"

	accessModeVolConfig := volumeConfig.ConstructClone()
	accessModeVolConfig.AccessMode = config.ReadWriteMany
	accessModeVolConfig.Protocol = config.Block

	volumeModeVolConfig := volumeConfig.ConstructClone()
	volumeModeVolConfig.VolumeMode = config.RawBlock
	volumeModeVolConfig.Protocol = config.File

	protocolVolConfig := volumeConfig.ConstructClone()
	protocolVolConfig.Protocol = config.Block

	for _, c := range []struct {
		name         string
		volumeConfig *storage.VolumeConfig
		valid        bool
		error        string
	}{
		{name: "volumeConfig", volumeConfig: volumeConfig, valid: true, error: ""},
		{name: "pvExists", volumeConfig: pvExistsVolConfig, valid: false, error: "already exists"},
		{name: "unknownSC", volumeConfig: unknownSCVolConfig, valid: false, error: "unknown storage class"},
		{name: "missingVolume", volumeConfig: missingVolConfig, valid: false, error: "volume noVol was not found"},
		{
			name:         "accessMode",
			volumeConfig: accessModeVolConfig,
			valid:        false,
			error:        "incompatible",
		},
		{name: "volumeMode", volumeConfig: volumeModeVolConfig, valid: false, error: "incompatible volume mode "},
		{name: "protocol", volumeConfig: protocolVolConfig, valid: false, error: "incompatible with the backend"},
	} {
		// The test code
		err = orchestrator.validateImportVolume(ctx(), c.volumeConfig)
		if err != nil {
			if c.valid {
				t.Errorf("%s: unexpected error %v", c.name, err)
			} else {
				if !strings.Contains(err.Error(), c.error) {
					t.Errorf("%s: expected %s but received error %v", c.name, c.error, err)
				}
			}
		} else if !c.valid {
			t.Errorf("%s: expected error but passed test", c.name)
		}
	}

	cleanup(t, orchestrator)
}

func TestValidateImportVolumeSanBackend(t *testing.T) {
	const (
		backendName     = "backend01"
		scName          = "sc01"
		volumeName      = "volume01"
		originalName    = "origVolume01"
		backendProtocol = config.Block
	)

	orchestrator, volumeConfig := importVolumeSetup(t, backendName, scName, volumeName, originalName, backendProtocol)

	_, err := orchestrator.AddVolume(ctx(), volumeConfig)
	if err != nil {
		t.Fatal("Unable to add volume: ", err)
	}

	// The volume exists on the backend with the original name since we added it above.
	// Set volumeConfig.InternalName to the expected volumeName post import.
	volumeConfig.InternalName = volumeName

	// Create VolumeConfig objects for the remaining error conditions
	protocolVolConfig := volumeConfig.ConstructClone()
	protocolVolConfig.Protocol = config.File

	ext4RawBlockFSVolConfig := volumeConfig.ConstructClone()
	ext4RawBlockFSVolConfig.VolumeMode = config.RawBlock
	ext4RawBlockFSVolConfig.Protocol = config.Block
	ext4RawBlockFSVolConfig.FileSystem = "ext4"

	for _, c := range []struct {
		name         string
		volumeConfig *storage.VolumeConfig
		valid        bool
		error        string
	}{
		{name: "protocol", volumeConfig: protocolVolConfig, valid: false, error: "incompatible with the backend"},
		{
			name:         "invalidFS",
			volumeConfig: ext4RawBlockFSVolConfig,
			valid:        false,
			error:        "cannot create raw-block volume",
		},
	} {
		// The test code
		err = orchestrator.validateImportVolume(ctx(), c.volumeConfig)
		if err != nil {
			if c.valid {
				t.Errorf("%s: unexpected error %v", c.name, err)
			} else {
				if !strings.Contains(err.Error(), c.error) {
					t.Errorf("%s: expected %s but received error %v", c.name, c.error, err)
				}
			}
		} else if !c.valid {
			t.Errorf("%s: expected error but passed test", c.name)
		}
	}

	cleanup(t, orchestrator)
}

func TestAddVolumePublication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Verify that the core calls the store client with the correct object, returning success
	mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), fakePub).Return(nil)

	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient

	err := orchestrator.AddVolumePublication(context.Background(), fakePub)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error adding volume publication: %v", err))
	assert.Contains(t, orchestrator.volumePublications.Map(), fakePub.VolumeName,
		"volume publication missing from orchestrator's cache")
	assert.NotNil(t, orchestrator.volumePublications.Get(fakePub.VolumeName, fakePub.NodeName),
		"volume publication missing from orchestrator's cache")
	assert.Equal(t, fakePub, orchestrator.volumePublications.Get(fakePub.VolumeName, fakePub.NodeName),
		"volume publication was not correctly added")
}

func TestAddVolumePublicationError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Verify that the core calls the store client with the correct object, but return an error
	mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), fakePub).Return(fmt.Errorf("fake error"))

	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient

	err := orchestrator.AddVolumePublication(context.Background(), fakePub)
	assert.NotNilf(t, err, "add volume publication did not return an error")
	assert.NotContains(t, orchestrator.volumePublications.Map(), fakePub.VolumeName,
		"volume publication was added orchestrator's cache")
}

func TestUpdateVolumePublication_FailsWithNoCacheEntry(t *testing.T) {
	// Set up mocks.
	mockCtrl := gomock.NewController(t)
	mockBackend := mockstorage.NewMockBackend(mockCtrl)

	// Set up test variables.
	var notSafeToAttach *bool
	volumeName := "foo"
	nodeName := "bar"

	// BackendUUID may be calls a number of times, mock it out.
	mockBackend.EXPECT().BackendUUID().Return("12345").AnyTimes() // Always return the right UUID

	// Initialize the orchestrator.
	orchestrator := getOrchestrator(t, false)

	// Make the update publication call.
	err := orchestrator.UpdateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)

	// Try to get the cached publication from the publications map.
	cachedPub, found := orchestrator.volumePublications.TryGet(volumeName, nodeName)

	assert.Error(t, err, "expected error")
	assert.False(t, found, "expected false value")
	assert.Nil(t, cachedPub, "expected non-nil value")
}

func TestUpdateVolumePublication_NoUpdate(t *testing.T) {
	// Set up mocks.
	mockCtrl := gomock.NewController(t)
	mockBackend := mockstorage.NewMockBackend(mockCtrl)

	// Set up test variables.
	var notSafeToAttach *bool
	volumeName := "foo"
	nodeName := "bar"
	vol := &storage.Volume{
		Config:      &storage.VolumeConfig{Name: volumeName},
		BackendUUID: "12345",
	}
	pub := &utils.VolumePublication{
		Name:       utils.GenerateVolumePublishName(volumeName, nodeName),
		VolumeName: volumeName,
		NodeName:   nodeName,
	}

	// BackendUUID may be calls a number of times, mock it out.
	mockBackend.EXPECT().BackendUUID().Return(vol.BackendUUID).AnyTimes() // Always return the right UUID

	// Initialize the orchestrator.
	orchestrator := getOrchestrator(t, false)

	// The publication must exist prior to calling update, add it to the rw cache.
	if err := orchestrator.volumePublications.Set(volumeName, nodeName, pub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Get the publication from the cache before calling update.
	oldPub := orchestrator.volumePublications.Get(volumeName, nodeName)

	// Make the update publication call.
	err := orchestrator.UpdateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)

	// Get the cached publication from the publications map.
	cachedPub := orchestrator.volumePublications.Get(volumeName, nodeName)

	assert.NoError(t, err, "expected no error")
	assert.NotNil(t, cachedPub, "expected non-nil value")
	assert.Equal(t, oldPub.NotSafeToAttach, cachedPub.NotSafeToAttach, "expected equal values")
}

func TestUpdateVolumePublication_UnpublishFails(t *testing.T) {
	// Set up mocks.
	mockCtrl := gomock.NewController(t)
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

	// Set up test variables.
	notSafeToAttach := new(bool)
	*notSafeToAttach = true
	volumeName := "foo"
	nodeName := "bar"
	vol := &storage.Volume{
		Config:      &storage.VolumeConfig{Name: volumeName},
		BackendUUID: "12345",
	}
	node := &utils.Node{Name: nodeName}
	pub := &utils.VolumePublication{
		Name:            utils.GenerateVolumePublishName(volumeName, nodeName),
		VolumeName:      volumeName,
		NodeName:        nodeName,
		NotSafeToAttach: !*notSafeToAttach, // make the publication false on the pub.
		Unpublished:     true,
	}

	// BackendUUID may be calls a number of times, mock it out.
	mockBackend.EXPECT().BackendUUID().Return(vol.BackendUUID).AnyTimes() // Always return the right UUID

	// Initialize the orchestrator and its caches.
	orchestrator := getOrchestrator(t, false)
	orchestrator.storeClient = mockStoreClient
	orchestrator.backends[vol.BackendUUID] = mockBackend
	orchestrator.volumes[volumeName] = vol
	orchestrator.nodes[nodeName] = node
	// The publication must exist prior to calling update, add it to the rw cache.
	if err := orchestrator.volumePublications.Set(volumeName, nodeName, pub); err != nil {
		t.Fatal("unable to set cache value")
	}

	mockBackend.EXPECT().UnpublishVolume(gomock.Any(), vol.Config, gomock.Any()).Return(fmt.Errorf("error"))
	err := orchestrator.updateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)

	// Try to get the cached publication from the publications map. It should still exist because the delete call failed.
	cachedPub, found := orchestrator.volumePublications.TryGet(volumeName, nodeName)

	assert.Error(t, err, "expected error")
	assert.True(t, found, "expected true value")          // True indicates the publication still exists.
	assert.NotNil(t, cachedPub, "expected non-nil value") // Non-nil indicates the publication still exists.
	assert.EqualValues(t, pub, cachedPub, "expected equal values")
}

func TestUpdateVolumePublication_CleanToDirty(t *testing.T) {
	// Set up mocks.
	mockCtrl := gomock.NewController(t)
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

	// Set up test variables.
	notSafeToAttach := new(bool)
	*notSafeToAttach = true
	volumeName := "foo"
	nodeName := "bar"
	vol := &storage.Volume{
		Config:      &storage.VolumeConfig{Name: volumeName},
		BackendUUID: "12345",
	}
	node := &utils.Node{Name: nodeName}
	pub := &utils.VolumePublication{
		Name:            utils.GenerateVolumePublishName(volumeName, nodeName),
		VolumeName:      volumeName,
		NodeName:        nodeName,
		NotSafeToAttach: !*notSafeToAttach, // make the publication false on the pub.
		Unpublished:     true,
	}

	// BackendUUID may be calls a number of times, mock it out.
	mockBackend.EXPECT().BackendUUID().Return(vol.BackendUUID).AnyTimes() // Always return the right UUID

	// Initialize the orchestrator and its caches.
	orchestrator := getOrchestrator(t, false)
	orchestrator.storeClient = mockStoreClient
	orchestrator.backends[vol.BackendUUID] = mockBackend
	orchestrator.volumes[volumeName] = vol
	orchestrator.nodes[nodeName] = node
	// The publication must exist prior to calling update, add it to the rw cache.
	if err := orchestrator.volumePublications.Set(volumeName, nodeName, pub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Get the cached publication before update to ensure NotSafeToAttach changes.
	oldPub := orchestrator.volumePublications.Get(volumeName, nodeName)

	// Mock out any clients and make the update publication call.
	mockStoreClient.EXPECT().UpdateVolumePublication(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBackend.EXPECT().UnpublishVolume(gomock.Any(), vol.Config, gomock.Any()).Return(nil).Times(1)
	err := orchestrator.updateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)
	cachedPub := orchestrator.volumePublications.Get(volumeName, nodeName)

	assert.NoError(t, err, "expected no error")
	assert.Equal(t, cachedPub.NotSafeToAttach, *notSafeToAttach, "expected equal values")
	assert.False(t, oldPub.NotSafeToAttach, "expected false value")
	assert.NotEqual(t, cachedPub.NotSafeToAttach, oldPub.NotSafeToAttach, "expected unequal values")
}

func TestUpdateVolumePublication_CleanToDirtyFails(t *testing.T) {
	// Set up mocks.
	mockCtrl := gomock.NewController(t)
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

	// Set up test variables.
	notSafeToAttach := new(bool)
	*notSafeToAttach = true
	volumeName := "foo"
	nodeName := "bar"
	vol := &storage.Volume{
		Config:      &storage.VolumeConfig{Name: volumeName},
		BackendUUID: "12345",
	}
	node := &utils.Node{Name: nodeName}
	pub := &utils.VolumePublication{
		Name:            utils.GenerateVolumePublishName(volumeName, nodeName),
		VolumeName:      volumeName,
		NodeName:        nodeName,
		NotSafeToAttach: !*notSafeToAttach, // make the publication false on the pub.
		Unpublished:     true,
	}

	// BackendUUID may be calls a number of times, mock it out.
	mockBackend.EXPECT().BackendUUID().Return(vol.BackendUUID).AnyTimes() // Always return the right UUID

	// Initialize the orchestrator and its caches.
	orchestrator := getOrchestrator(t, false)
	orchestrator.storeClient = mockStoreClient
	orchestrator.backends[vol.BackendUUID] = mockBackend
	orchestrator.volumes[volumeName] = vol
	orchestrator.nodes[nodeName] = node
	// The publication must exist prior to calling update, add it to the rw cache.
	if err := orchestrator.volumePublications.Set(volumeName, nodeName, pub); err != nil {
		t.Fatal("unable to set cache value")
	}

	mockBackend.EXPECT().UnpublishVolume(gomock.Any(), vol.Config, gomock.Any()).Return(nil)
	mockStoreClient.EXPECT().UpdateVolumePublication(gomock.Any(), gomock.Any()).Return(errors.New("update publication failed"))
	err := orchestrator.updateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)

	// Try to get the cached publication from the publications map. It should still exist because the update call failed.
	cachedPub, found := orchestrator.volumePublications.TryGet(volumeName, nodeName)

	assert.Error(t, err, "expected error")
	assert.True(t, found, "expected true value")          // True indicates the publication still exists.
	assert.NotNil(t, cachedPub, "expected non-nil value") // Non-nil indicates the publication still exists.
	assert.EqualValues(t, pub, cachedPub, "expected equal values")
}

func TestUpdateVolumePublication_DirtyToClean(t *testing.T) {
	// Set up mocks.
	mockCtrl := gomock.NewController(t)
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

	// Set up test variables.
	notSafeToAttach := new(bool)
	volumeName := "foo"
	nodeName := "bar"
	vol := &storage.Volume{
		Config:      &storage.VolumeConfig{Name: volumeName},
		BackendUUID: "12345",
	}
	node := &utils.Node{Name: nodeName}
	pub := &utils.VolumePublication{
		Name:            utils.GenerateVolumePublishName(volumeName, nodeName),
		VolumeName:      volumeName,
		NodeName:        nodeName,
		NotSafeToAttach: !*notSafeToAttach, // make the publication false on the pub.
		Unpublished:     true,
	}

	// BackendUUID may be calls a number of times, mock it out.
	mockBackend.EXPECT().BackendUUID().Return(vol.BackendUUID).AnyTimes() // Always return the right UUID

	// Initialize the orchestrator and its caches.
	orchestrator := getOrchestrator(t, false)
	orchestrator.storeClient = mockStoreClient
	orchestrator.backends[vol.BackendUUID] = mockBackend
	orchestrator.volumes[volumeName] = vol
	orchestrator.nodes[nodeName] = node
	// The publication must exist prior to calling update, add it to the rw cache.
	if err := orchestrator.volumePublications.Set(volumeName, nodeName, pub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Test if DeleteVolumePublication fails, that the entire operation fails.
	mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), gomock.Any()).Return(errors.New("pub not found")).Times(1)
	err := orchestrator.updateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)

	// Try to get the cached publication from the publications map. It should still exist because the delete call failed.
	cachedPub, found := orchestrator.volumePublications.TryGet(volumeName, nodeName)

	assert.Error(t, err, "expected error")
	assert.True(t, found, "expected true value")          // True indicates the publication still exists.
	assert.NotNil(t, cachedPub, "expected non-nil value") // Non-nil indicates the publication still exists.
	assert.EqualValues(t, pub, cachedPub, "expected equal values")

	// Test the happy path for cleaning a publication works and the publication is deleted.
	mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err = orchestrator.updateVolumePublication(context.Background(), volumeName, nodeName, notSafeToAttach)

	// Try to get the cached publication from the publications map.
	cachedPub, found = orchestrator.volumePublications.TryGet(volumeName, nodeName)

	assert.NoError(t, err, "expected no error")
	assert.False(t, found, "expected false value") // False indicates the publication was removed.
	assert.Nil(t, cachedPub, "expected nil value") // Nil indicates the publication was removed.
}

func TestGetVolumePublication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	actualPub, err := orchestrator.GetVolumePublication(context.Background(), fakePub.VolumeName, fakePub.NodeName)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error getting volume publication: %v", err))
	assert.Equal(t, fakePub, actualPub, "volume publication was not correctly retrieved")
}

func TestGetVolumePublicationNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient

	actualPub, err := orchestrator.GetVolumePublication(context.Background(), "NotFound", "NotFound")
	assert.NotNilf(t, err, fmt.Sprintf("unexpected success getting volume publication: %v", err))
	assert.True(t, utils.IsNotFoundError(err), "incorrect error type returned")
	assert.Empty(t, actualPub, "non-empty publication returned")
}

func TestGetVolumePublicationError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Simulate a bootstrap error
	orchestrator.bootstrapError = fmt.Errorf("some error")

	actualPub, err := orchestrator.GetVolumePublication(context.Background(), fakePub.VolumeName, fakePub.NodeName)
	assert.NotNilf(t, err, fmt.Sprintf("unexpected success getting volume publication: %v", err))
	assert.False(t, utils.IsNotFoundError(err), "incorrect error type returned")
	assert.Empty(t, actualPub, "non-empty publication returned")
}

func TestListVolumePublications(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub1 := &utils.VolumePublication{
		Name:            "foo/bar",
		NodeName:        "bar",
		VolumeName:      "foo",
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: false,
	}
	fakePub2 := &utils.VolumePublication{
		Name:            "baz/biz",
		NodeName:        "biz",
		VolumeName:      "baz",
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: false,
	}
	fakePub3 := &utils.VolumePublication{
		Name:            fmt.Sprintf("%s/buz", fakePub1.VolumeName),
		NodeName:        "buz",
		VolumeName:      fakePub1.VolumeName,
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: true,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	err := orchestrator.volumePublications.Set(fakePub1.VolumeName, fakePub1.NodeName, fakePub1)
	err = orchestrator.volumePublications.Set(fakePub2.VolumeName, fakePub2.NodeName, fakePub2)
	err = orchestrator.volumePublications.Set(fakePub3.VolumeName, fakePub3.NodeName, fakePub3)
	if err != nil {
		t.Fatal("unable to set cache value")
	}

	expectedAllPubs := []*utils.VolumePublicationExternal{
		fakePub1.ConstructExternal(),
		fakePub2.ConstructExternal(),
		fakePub3.ConstructExternal(),
	}
	expectedCleanPubs := []*utils.VolumePublicationExternal{
		fakePub1.ConstructExternal(),
		fakePub2.ConstructExternal(),
	}
	expectedDirtyPubs := []*utils.VolumePublicationExternal{
		fakePub3.ConstructExternal(),
	}

	actualPubs, err := orchestrator.ListVolumePublications(context.Background(), nil)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedAllPubs, actualPubs, "incorrect publication list returned")

	actualPubs, err = orchestrator.ListVolumePublications(context.Background(), utils.Ptr(false))
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedCleanPubs, actualPubs, "incorrect publication list returned")

	actualPubs, err = orchestrator.ListVolumePublications(context.Background(), utils.Ptr(true))
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedDirtyPubs, actualPubs, "incorrect publication list returned")
}

func TestListVolumePublicationsNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient

	actualPubs, err := orchestrator.ListVolumePublications(context.Background(), nil)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.Empty(t, actualPubs, "non-empty publication list returned")
}

func TestListVolumePublicationsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Simulate a bootstrap error
	orchestrator.bootstrapError = fmt.Errorf("some error")

	actualPubs, err := orchestrator.ListVolumePublications(context.Background(), nil)
	assert.NotNil(t, err, fmt.Sprintf("unexpected success listing volume publications"))
	assert.Empty(t, actualPubs, "non-empty publication list returned")
}

func TestListVolumePublicationsForVolume(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub1 := &utils.VolumePublication{
		Name:            "foo/bar",
		NodeName:        "bar",
		VolumeName:      "foo",
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: false,
	}
	fakePub2 := &utils.VolumePublication{
		Name:            "baz/biz",
		NodeName:        "biz",
		VolumeName:      "baz",
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: false,
	}
	fakePub3 := &utils.VolumePublication{
		Name:            fmt.Sprintf("%s/buz", fakePub1.VolumeName),
		NodeName:        "buz",
		VolumeName:      fakePub1.VolumeName,
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: true,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	err := orchestrator.volumePublications.Set(fakePub1.VolumeName, fakePub1.NodeName, fakePub1)
	err = orchestrator.volumePublications.Set(fakePub2.VolumeName, fakePub2.NodeName, fakePub2)
	err = orchestrator.volumePublications.Set(fakePub3.VolumeName, fakePub3.NodeName, fakePub3)
	if err != nil {
		t.Fatal("unable to set cache value")
	}

	expectedAllPubs := []*utils.VolumePublicationExternal{fakePub1.ConstructExternal(), fakePub3.ConstructExternal()}
	expectedCleanPubs := []*utils.VolumePublicationExternal{fakePub1.ConstructExternal()}
	expectedDirtyPubs := []*utils.VolumePublicationExternal{fakePub3.ConstructExternal()}

	actualPubs, err := orchestrator.ListVolumePublicationsForVolume(context.Background(),
		fakePub1.VolumeName, nil)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedAllPubs, actualPubs, "incorrect publication list returned")

	actualPubs, err = orchestrator.ListVolumePublicationsForVolume(context.Background(),
		fakePub1.VolumeName, utils.Ptr(false))
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedCleanPubs, actualPubs, "incorrect publication list returned")

	actualPubs, err = orchestrator.ListVolumePublicationsForVolume(context.Background(),
		fakePub1.VolumeName, utils.Ptr(true))
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedDirtyPubs, actualPubs, "incorrect publication list returned")
}

func TestListVolumePublicationsForVolumeNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	actualPubs, err := orchestrator.ListVolumePublicationsForVolume(context.Background(), "NotFound", nil)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.Empty(t, actualPubs, "non-empty publication list returned")
}

func TestListVolumePublicationsForVolumeError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Simulate a bootstrap error
	orchestrator.bootstrapError = fmt.Errorf("some error")

	actualPubs, err := orchestrator.ListVolumePublicationsForVolume(context.Background(), fakePub.VolumeName, nil)
	assert.NotNil(t, err, fmt.Sprintf("unexpected success listing volume publications"))
	assert.Empty(t, actualPubs, "non-empty publication list returned")
}

func TestListVolumePublicationsForNode(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub1 := &utils.VolumePublication{
		Name:            "foo/bar",
		NodeName:        "bar",
		VolumeName:      "foo",
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: false,
	}
	fakePub2 := &utils.VolumePublication{
		Name:            "baz/biz",
		NodeName:        "biz",
		VolumeName:      "baz",
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: false,
	}
	fakePub3 := &utils.VolumePublication{
		Name:            fmt.Sprintf("%s/buz", fakePub1.VolumeName),
		NodeName:        "buz",
		VolumeName:      fakePub1.VolumeName,
		ReadOnly:        true,
		AccessMode:      1,
		NotSafeToAttach: true,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	err := orchestrator.volumePublications.Set(fakePub1.VolumeName, fakePub1.NodeName, fakePub1)
	err = orchestrator.volumePublications.Set(fakePub2.VolumeName, fakePub2.NodeName, fakePub2)
	err = orchestrator.volumePublications.Set(fakePub3.VolumeName, fakePub3.NodeName, fakePub3)
	if err != nil {
		t.Fatal("unable to set cache value")
	}

	expectedAllPubs := []*utils.VolumePublicationExternal{fakePub2.ConstructExternal()}
	expectedCleanPubs := []*utils.VolumePublicationExternal{fakePub2.ConstructExternal()}
	expectedDirtyPubs := []*utils.VolumePublicationExternal{}

	actualPubs, err := orchestrator.ListVolumePublicationsForNode(context.Background(), fakePub2.NodeName, nil)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedAllPubs, actualPubs, "incorrect publication list returned")

	actualPubs, err = orchestrator.ListVolumePublicationsForNode(context.Background(), fakePub2.NodeName,
		utils.Ptr(false))
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedCleanPubs, actualPubs, "incorrect publication list returned")

	actualPubs, err = orchestrator.ListVolumePublicationsForNode(context.Background(), fakePub2.NodeName,
		utils.Ptr(true))
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.ElementsMatch(t, expectedDirtyPubs, actualPubs, "incorrect publication list returned")
}

func TestListVolumePublicationsForNodeNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	actualPubs, err := orchestrator.ListVolumePublicationsForNode(context.Background(), "NotFound", nil)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error listing volume publications: %v", err))
	assert.Empty(t, actualPubs, "non-empty publication list returned")
}

func TestListVolumePublicationsForNodeError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	// Simulate a bootstrap error
	orchestrator.bootstrapError = fmt.Errorf("some error")

	actualPubs, err := orchestrator.ListVolumePublicationsForVolume(context.Background(), fakePub.NodeName, nil)
	assert.NotNil(t, err, fmt.Sprintf("unexpected success listing volume publications"))
	assert.Empty(t, actualPubs, "non-empty publication list returned")
}

func TestDeleteVolumePublication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Create a fake VolumePublication
	fakePub1 := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	fakePub2 := &utils.VolumePublication{
		Name:       "baz/biz",
		NodeName:   "biz",
		VolumeName: "baz",
		ReadOnly:   true,
		AccessMode: 1,
	}
	fakePub3 := &utils.VolumePublication{
		Name:       fmt.Sprintf("%s/buz", fakePub1.VolumeName),
		NodeName:   "buz",
		VolumeName: fakePub1.VolumeName,
		ReadOnly:   true,
		AccessMode: 1,
	}
	fakeNode := &utils.Node{Name: "biz"}
	fakeNode2 := &utils.Node{Name: "buz"}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	err := orchestrator.volumePublications.Set(fakePub1.VolumeName, fakePub1.NodeName, fakePub1)
	err = orchestrator.volumePublications.Set(fakePub2.VolumeName, fakePub2.NodeName, fakePub2)
	err = orchestrator.volumePublications.Set(fakePub3.VolumeName, fakePub3.NodeName, fakePub3)
	if err != nil {
		t.Fatal("unable to set cache value")
	}

	orchestrator.nodes = map[string]*utils.Node{fakeNode.Name: fakeNode, fakeNode2.Name: fakeNode2}

	// Verify if this is the last nodeID for a given volume the volume entry is completely removed from the cache
	mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), fakePub2).Return(nil)
	err = orchestrator.DeleteVolumePublication(context.Background(), fakePub2.VolumeName, fakePub2.NodeName)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error deleting volume publication: %v", err))

	cachedPub, ok := orchestrator.volumePublications.TryGet(fakePub2.VolumeName, fakePub2.NodeName)
	assert.False(t, ok, "publication not removed from the cache")
	assert.Nil(t, cachedPub, "publication not removed from the cache")

	// Verify if this is not the last nodeID for a given volume the volume entry is not removed from the cache
	mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), fakePub3).Return(nil)
	err = orchestrator.DeleteVolumePublication(context.Background(), fakePub3.VolumeName, fakePub3.NodeName)
	assert.NoError(t, err, fmt.Sprintf("unexpected error deleting volume publication: %v", err))

	cachedPub, ok = orchestrator.volumePublications.TryGet(fakePub3.VolumeName, fakePub3.NodeName)
	assert.False(t, ok, "publication not removed from the cache")
	assert.Nil(t, cachedPub, "publication not removed from the cache")

	cachedPub, ok = orchestrator.volumePublications.TryGet(fakePub1.VolumeName, fakePub1.NodeName)
	assert.True(t, ok, "publication improperly removed from the cache")
	assert.NotNil(t, cachedPub, "publication improperly removed from the cache")
}

func TestDeleteVolumePublicationNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client.
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase.
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication.
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	fakeNode := &utils.Node{Name: fakePub.NodeName}

	// Create an instance of the orchestrator for this test.
	orchestrator := getOrchestrator(t, false)
	orchestrator.nodes = map[string]*utils.Node{fakePub.NodeName: fakeNode}
	// Add the mocked objects to the orchestrator.
	orchestrator.storeClient = mockStoreClient
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set publication in the cache.")
	}

	// When Trident can't find the publication in the cache, it should ask the persistent store.
	ctx := context.Background()
	mockStoreClient.EXPECT().DeleteVolumePublication(ctx, fakePub).Return(nil).Times(1)
	err := orchestrator.DeleteVolumePublication(ctx, fakePub.VolumeName, fakePub.NodeName)
	cachedPubs := orchestrator.volumePublications.Map()
	assert.NoError(t, err, fmt.Sprintf("unexpected error deleting volume publication"))
	assert.Nil(t, cachedPubs[fakePub.VolumeName][fakePub.NodeName], "expected no cache entry")
}

func TestDeleteVolumePublicationNotFoundPersistence(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	fakeNode := &utils.Node{Name: fakePub.NodeName}

	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}
	orchestrator.nodes = map[string]*utils.Node{fakeNode.Name: fakeNode}

	// Verify delete is idempotent when the persistence object is missing
	mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), fakePub).Return(utils.NotFoundError("not found"))
	err := orchestrator.DeleteVolumePublication(context.Background(), fakePub.VolumeName, fakePub.NodeName)
	assert.Nilf(t, err, fmt.Sprintf("unexpected error deleting volume publication: %v", err))
	assert.NotContains(t, orchestrator.volumePublications.Map(), fakePub.VolumeName,
		"publication not properly removed from cache")
}

func TestDeleteVolumePublicationError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Create a fake VolumePublication
	fakePub := &utils.VolumePublication{
		Name:       "foo/bar",
		NodeName:   "bar",
		VolumeName: "foo",
		ReadOnly:   true,
		AccessMode: 1,
	}
	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	// Populate volume publications
	if err := orchestrator.volumePublications.Set(fakePub.VolumeName, fakePub.NodeName, fakePub); err != nil {
		t.Fatal("unable to set cache value")
	}

	mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), fakePub).Return(fmt.Errorf("some error"))

	err := orchestrator.DeleteVolumePublication(context.Background(), fakePub.VolumeName, fakePub.NodeName)
	assert.NotNil(t, err, fmt.Sprintf("unexpected success deleting volume publication"))
	assert.False(t, utils.IsNotFoundError(err), "incorrect error type returned")
	assert.Equal(t, fakePub, orchestrator.volumePublications.Get(fakePub.VolumeName, fakePub.NodeName),
		"publication improperly removed/updated in cache")
}

func TestAddNode(t *testing.T) {
	node := &utils.Node{
		Name:           "testNode",
		IQN:            "myIQN",
		IPs:            []string{"1.1.1.1", "2.2.2.2"},
		TopologyLabels: map[string]string{"topology.kubernetes.io/region": "Region1"},
		Deleted:        false,
	}
	orchestrator := getOrchestrator(t, false)
	if err := orchestrator.AddNode(ctx(), node, nil); err != nil {
		t.Errorf("adding node failed; %v", err)
	}
}

func TestGetNode(t *testing.T) {
	orchestrator := getOrchestrator(t, false)
	expectedNode := &utils.Node{
		Name:           "testNode",
		IQN:            "myIQN",
		IPs:            []string{"1.1.1.1", "2.2.2.2"},
		TopologyLabels: map[string]string{"topology.kubernetes.io/region": "Region1"},
		Deleted:        false,
	}
	unexpectedNode := &utils.Node{
		Name:           "testNode2",
		IQN:            "myOtherIQN",
		IPs:            []string{"3.3.3.3", "4.4.4.4"},
		TopologyLabels: map[string]string{"topology.kubernetes.io/region": "Region2"},
		Deleted:        false,
	}
	initialNodes := map[string]*utils.Node{}
	initialNodes[expectedNode.Name] = expectedNode
	initialNodes[unexpectedNode.Name] = unexpectedNode
	orchestrator.nodes = initialNodes

	actualNode, err := orchestrator.GetNode(ctx(), expectedNode.Name)
	if err != nil {
		t.Errorf("error getting node; %v", err)
	}

	if actualNode != expectedNode {
		t.Errorf("Did not get expected node back; expected %+v, got %+v", expectedNode, actualNode)
	}
}

func TestListNodes(t *testing.T) {
	orchestrator := getOrchestrator(t, false)
	expectedNode1 := &utils.Node{
		Name:    "testNode",
		IQN:     "myIQN",
		IPs:     []string{"1.1.1.1", "2.2.2.2"},
		Deleted: false,
	}
	expectedNode2 := &utils.Node{
		Name:    "testNode2",
		IQN:     "myOtherIQN",
		IPs:     []string{"3.3.3.3", "4.4.4.4"},
		Deleted: false,
	}
	initialNodes := map[string]*utils.Node{}
	initialNodes[expectedNode1.Name] = expectedNode1
	initialNodes[expectedNode2.Name] = expectedNode2
	orchestrator.nodes = initialNodes
	expectedNodes := []*utils.Node{expectedNode1, expectedNode2}

	actualNodes, err := orchestrator.ListNodes(ctx())
	if err != nil {
		t.Errorf("error listing nodes; %v", err)
	}

	if !unorderedNodeSlicesEqual(actualNodes, expectedNodes) {
		t.Errorf("node list values do not match; expected %v, found %v", expectedNodes, actualNodes)
	}
}

func unorderedNodeSlicesEqual(x, y []*utils.Node) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of node pointers -> int
	diff := make(map[*utils.Node]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the node _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

func TestDeleteNode(t *testing.T) {
	orchestrator := getOrchestrator(t, false)
	initialNode := &utils.Node{
		Name:    "testNode",
		IQN:     "myIQN",
		IPs:     []string{"1.1.1.1", "2.2.2.2"},
		Deleted: false,
	}
	initialNodes := map[string]*utils.Node{}
	initialNodes[initialNode.Name] = initialNode
	orchestrator.nodes = initialNodes

	if err := orchestrator.DeleteNode(ctx(), initialNode.Name); err != nil {
		t.Errorf("error deleting node; %v", err)
	}

	if _, ok := orchestrator.nodes[initialNode.Name]; ok {
		t.Errorf("node was not properly deleted")
	}
}

func TestSnapshotVolumes(t *testing.T) {
	mockPools := tu.GetFakePools()
	orchestrator := getOrchestrator(t, false)

	errored := false
	for _, c := range []struct {
		name      string
		protocol  config.Protocol
		poolNames []string
	}{
		{
			name:      "fast-a",
			protocol:  config.File,
			poolNames: []string{tu.FastSmall, tu.FastThinOnly},
		},
		{
			name:      "fast-b",
			protocol:  config.File,
			poolNames: []string{tu.FastThinOnly, tu.FastUniqueAttr},
		},
		{
			name:      "slow-file",
			protocol:  config.File,
			poolNames: []string{tu.SlowNoSnapshots, tu.SlowSnapshots},
		},
		{
			name:      "slow-block",
			protocol:  config.Block,
			poolNames: []string{tu.SlowNoSnapshots, tu.SlowSnapshots, tu.MediumOverlap},
		},
	} {
		pools := make(map[string]*fake.StoragePool, len(c.poolNames))
		for _, poolName := range c.poolNames {
			pools[poolName] = mockPools[poolName]
		}
		cfg, err := fakedriver.NewFakeStorageDriverConfigJSON(c.name, c.protocol, pools, make([]fake.Volume, 0))
		if err != nil {
			t.Fatalf("Unable to generate cfg JSON for %s: %v", c.name, err)
		}
		_, err = orchestrator.AddBackend(ctx(), cfg, "")
		if err != nil {
			t.Errorf("Unable to add backend %s: %v", c.name, err)
			errored = true
		}
		orchestrator.mutex.Lock()
		backend, err := orchestrator.getBackendByBackendName(c.name)
		if err != nil {
			t.Fatalf("Backend %s not stored in orchestrator", c.name)
		}
		persistentBackend, err := orchestrator.storeClient.GetBackend(ctx(), c.name)
		if err != nil {
			t.Fatalf("Unable to get backend %s from persistent store: %v", c.name, err)
		} else if !reflect.DeepEqual(
			backend.ConstructPersistent(ctx()),
			persistentBackend,
		) {
			t.Error("Wrong data stored for backend ", c.name)
		}
		orchestrator.mutex.Unlock()
	}
	if errored {
		t.Fatal("Failed to add all backends; aborting remaining tests.")
	}

	// Add storage classes
	storageClasses := []storageClassTest{
		{
			config: &storageclass.Config{
				Name: "slow",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(40),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "slow-file", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
			},
		},
		{
			config: &storageclass.Config{
				Name: "fast",
				Attributes: map[string]sa.Request{
					sa.IOPS:             sa.NewIntRequest(2000),
					sa.Snapshots:        sa.NewBoolRequest(true),
					sa.ProvisioningType: sa.NewStringRequest("thin"),
				},
			},
			expected: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
	}
	for _, s := range storageClasses {
		_, err := orchestrator.AddStorageClass(ctx(), s.config)
		if err != nil {
			t.Errorf("Unable to add storage class %s: %v", s.config.Name, err)
			continue
		}
		validateStorageClass(t, orchestrator, s.config.Name, s.expected)
	}

	for _, s := range []struct {
		name            string
		config          *storage.VolumeConfig
		expectedSuccess bool
		expectedMatches []*tu.PoolMatch
	}{
		{
			name:            "file",
			config:          tu.GenerateVolumeConfig("file", 1, "fast", config.File),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "fast-a", Pool: tu.FastSmall},
				{Backend: "fast-a", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastThinOnly},
				{Backend: "fast-b", Pool: tu.FastUniqueAttr},
			},
		},
		{
			name:            "block",
			config:          tu.GenerateVolumeConfig("block", 1, "slow", config.Block),
			expectedSuccess: true,
			expectedMatches: []*tu.PoolMatch{
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
				{Backend: "slow-block", Pool: tu.SlowSnapshots},
			},
		},
	} {
		// Create the source volume
		_, err := orchestrator.AddVolume(ctx(), s.config)
		if err != nil {
			t.Errorf("%s: could not add volume: %v", s.name, err)
			continue
		}

		orchestrator.mutex.Lock()
		volume, found := orchestrator.volumes[s.config.Name]
		if s.expectedSuccess && !found {
			t.Errorf("%s: did not get volume where expected.", s.name)
			continue
		}
		orchestrator.mutex.Unlock()

		// Now take a snapshot and ensure everything looks fine
		snapshotName := "snapshot-" + uuid.New().String()
		snapshotConfig := &storage.SnapshotConfig{
			Version:    config.OrchestratorAPIVersion,
			Name:       snapshotName,
			VolumeName: volume.Config.Name,
		}
		snapshotExternal, err := orchestrator.CreateSnapshot(ctx(), snapshotConfig)
		if err != nil {
			t.Fatalf("%s: got unexpected error creating snapshot: %v", s.name, err)
		}

		orchestrator.mutex.Lock()
		// Snapshot should be registered in the store
		persistentSnapshot, err := orchestrator.storeClient.GetSnapshot(ctx(), volume.Config.Name, snapshotName)
		if err != nil {
			t.Errorf("%s: unable to communicate with backing store: %v", snapshotName, err)
		}
		persistentSnapshotExternal := persistentSnapshot.ConstructExternal()
		if !reflect.DeepEqual(persistentSnapshotExternal, snapshotExternal) {
			t.Errorf(
				"%s: external snapshot %s stored in backend does not match created snapshot.",
				snapshotName, persistentSnapshot.Config.Name,
			)
			externalSnapshotJSON, err := json.Marshal(persistentSnapshotExternal)
			if err != nil {
				t.Fatal("Unable to remarshal JSON: ", err)
			}
			origSnapshotJSON, err := json.Marshal(snapshotExternal)
			if err != nil {
				t.Fatal("Unable to remarshal JSON: ", err)
			}
			t.Logf("\tExpected: %s\n\tGot: %s\n", string(externalSnapshotJSON), string(origSnapshotJSON))
		}
		orchestrator.mutex.Unlock()

		err = orchestrator.DeleteSnapshot(ctx(), volume.Config.Name, snapshotName)
		if err != nil {
			t.Fatalf("%s: got unexpected error deleting snapshot: %v", s.name, err)
		}

		orchestrator.mutex.Lock()
		// Snapshot should not be registered in the store
		persistentSnapshot, err = orchestrator.storeClient.GetSnapshot(ctx(), volume.Config.Name, snapshotName)
		if err != nil && !persistentstore.MatchKeyNotFoundErr(err) {
			t.Errorf("%s: unable to communicate with backing store: %v", snapshotName, err)
		}
		if persistentSnapshot != nil {
			t.Errorf("%s: got snapshot when not expected.", snapshotName)
			continue
		}
		orchestrator.mutex.Unlock()
	}

	cleanup(t, orchestrator)
}

func TestGetProtocol(t *testing.T) {
	orchestrator := getOrchestrator(t, false)

	type accessVariables struct {
		volumeMode config.VolumeMode
		accessMode config.AccessMode
		protocol   config.Protocol
		expected   config.Protocol
	}

	accessModesPositiveTests := []accessVariables{
		{config.Filesystem, config.ModeAny, config.ProtocolAny, config.ProtocolAny},
		{config.Filesystem, config.ModeAny, config.File, config.File},
		{config.Filesystem, config.ModeAny, config.Block, config.Block},
		{config.Filesystem, config.ModeAny, config.BlockOnFile, config.BlockOnFile},
		{config.Filesystem, config.ReadWriteOnce, config.ProtocolAny, config.ProtocolAny},
		{config.Filesystem, config.ReadWriteOnce, config.File, config.File},
		{config.Filesystem, config.ReadWriteOnce, config.Block, config.Block},
		{config.Filesystem, config.ReadWriteOnce, config.BlockOnFile, config.BlockOnFile},
		{config.Filesystem, config.ReadOnlyMany, config.Block, config.Block},
		{config.Filesystem, config.ReadOnlyMany, config.ProtocolAny, config.ProtocolAny},
		{config.Filesystem, config.ReadOnlyMany, config.File, config.File},
		{config.Filesystem, config.ReadWriteMany, config.ProtocolAny, config.File},
		{config.Filesystem, config.ReadWriteMany, config.File, config.File},
		// {config.Filesystem, config.ReadWriteMany, config.Block, config.ProtocolAny},
		// {config.Filesystem, config.ReadWriteMany, config.BlockOnFile, config.ProtocolAny},
		{config.RawBlock, config.ModeAny, config.ProtocolAny, config.Block},
		// {config.RawBlock, config.ModeAny, config.File, config.ProtocolAny},
		{config.RawBlock, config.ModeAny, config.Block, config.Block},
		{config.RawBlock, config.ReadWriteOnce, config.ProtocolAny, config.Block},
		// {config.RawBlock, config.ReadWriteOnce, config.File, config.ProtocolAny},
		{config.RawBlock, config.ReadWriteOnce, config.Block, config.Block},
		{config.RawBlock, config.ReadOnlyMany, config.ProtocolAny, config.Block},
		// {config.RawBlock, config.ReadOnlyMany, config.File, config.ProtocolAny},
		{config.RawBlock, config.ReadOnlyMany, config.Block, config.Block},
		{config.RawBlock, config.ReadWriteMany, config.ProtocolAny, config.Block},
		// {config.RawBlock, config.ReadWriteMany, config.File, config.ProtocolAny},
		{config.RawBlock, config.ReadWriteMany, config.Block, config.Block},
	}

	accessModesNegativeTests := []accessVariables{
		{config.Filesystem, config.ReadOnlyMany, config.BlockOnFile, config.ProtocolAny},
		{config.Filesystem, config.ReadWriteMany, config.Block, config.ProtocolAny},
		{config.Filesystem, config.ReadWriteMany, config.BlockOnFile, config.ProtocolAny},
		{config.RawBlock, config.ModeAny, config.File, config.ProtocolAny},
		{config.RawBlock, config.ModeAny, config.BlockOnFile, config.ProtocolAny},
		{config.RawBlock, config.ReadWriteOnce, config.File, config.ProtocolAny},
		{config.RawBlock, config.ReadWriteOnce, config.BlockOnFile, config.ProtocolAny},

		{config.RawBlock, config.ReadOnlyMany, config.File, config.ProtocolAny},
		{config.RawBlock, config.ReadWriteMany, config.File, config.ProtocolAny},

		{config.RawBlock, config.ReadOnlyMany, config.BlockOnFile, config.ProtocolAny},
		{config.RawBlock, config.ReadWriteMany, config.BlockOnFile, config.ProtocolAny},
	}

	for _, tc := range accessModesPositiveTests {
		protocolLocal, err := orchestrator.getProtocol(ctx(), tc.volumeMode, tc.accessMode, tc.protocol)
		assert.Nil(t, err, nil)
		assert.Equal(t, tc.expected, protocolLocal, "expected both the protocols to be equal!")
	}

	for _, tc := range accessModesNegativeTests {
		protocolLocal, err := orchestrator.getProtocol(ctx(), tc.volumeMode, tc.accessMode, tc.protocol)
		assert.NotNil(t, err)
		assert.Equal(t, tc.expected, protocolLocal, "expected both the protocols to be equal!")
	}
}

func TestGetBackend(t *testing.T) {
	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)

	// Set fake values
	backendName := "foobar"
	backendUUID := "1234"
	// Create the expected return object
	expectedBackendExternal := &storage.BackendExternal{
		Name:        backendName,
		BackendUUID: backendUUID,
	}

	// Create a mocked backend
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	// Set backend behavior we don't care about for this testcase
	mockBackend.EXPECT().Name().Return(backendName).AnyTimes()        // Always return the fake name
	mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes() // Always return the fake uuid
	// Set backend behavior we do care about for this testcase
	mockBackend.EXPECT().ConstructExternal(gomock.Any()).Return(expectedBackendExternal) // Return the expected object

	// Create an instance of the orchestrator
	orchestrator := getOrchestrator(t, false)
	// Add the mocked backend to the orchestrator
	orchestrator.backends[backendUUID] = mockBackend

	// Run the test
	actualBackendExternal, err := orchestrator.GetBackend(ctx(), backendName)

	// Verify the results
	assert.Nilf(t, err, "Error getting backend; %v", err)
	assert.Equal(t, expectedBackendExternal, actualBackendExternal, "Did not get the expected backend object")
}

func TestGetBackendByBackendUUID(t *testing.T) {
	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)

	// Set fake values
	backendName := "foobar"
	backendUUID := "1234"
	// Create the expected return object
	expectedBackendExternal := &storage.BackendExternal{
		Name:        backendName,
		BackendUUID: backendUUID,
	}

	// Create mocked backend that returns the expected object
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().ConstructExternal(gomock.Any()).Times(1).Return(expectedBackendExternal)

	// Create an instance of the orchestrator
	orchestrator := getOrchestrator(t, false)
	// Add the mocked backend to the orchestrator
	orchestrator.backends[backendUUID] = mockBackend

	// Run the test
	actualBackendExternal, err := orchestrator.GetBackendByBackendUUID(ctx(), backendUUID)

	// Verify the results
	assert.Nilf(t, err, "Error getting backend; %v", err)
	assert.Equal(t, expectedBackendExternal, actualBackendExternal, "Did not get the expected backend object")
}

func TestListBackends(t *testing.T) {
	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create list of 2 fake objects that we expect to be returned
	expectedBackendExternal1 := &storage.BackendExternal{
		Name:        "foo",
		BackendUUID: "12345",
	}
	expectedBackendExternal2 := &storage.BackendExternal{
		Name:        "bar",
		BackendUUID: "67890",
	}
	expectedBackendList := []*storage.BackendExternal{expectedBackendExternal1, expectedBackendExternal2}

	// Create 2 mocked backends that each return one of the expected fake objects when called
	mockBackend1 := mockstorage.NewMockBackend(mockCtrl)
	mockBackend1.EXPECT().ConstructExternal(gomock.Any()).Return(expectedBackendExternal1)
	mockBackend2 := mockstorage.NewMockBackend(mockCtrl)
	mockBackend2.EXPECT().ConstructExternal(gomock.Any()).Return(expectedBackendExternal2)

	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked backends to the orchestrator
	orchestrator.backends[expectedBackendExternal1.BackendUUID] = mockBackend1
	orchestrator.backends[expectedBackendExternal2.BackendUUID] = mockBackend2

	// Perform the test
	actualBackendList, err := orchestrator.ListBackends(ctx())

	// Verify the results
	assert.Nilf(t, err, "Error listing backends; %v", err)
	assert.ElementsMatch(t, expectedBackendList, actualBackendList, "Did not get expected list of backends")
}

func TestDeleteBackend(t *testing.T) {
	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)

	// Set fake values
	backendName := "foobar"
	backendUUID := "1234"

	// Create a mocked storage backend
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	// Set backend behavior we don't care about for this testcase
	mockBackend.EXPECT().Name().Return(backendName).AnyTimes()                  // Always return the fake name
	mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()           // Always return the fake UUID
	mockBackend.EXPECT().ConfigRef().Return("").AnyTimes()                      // Always return an empty configRef
	mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()               // Always return a fake driver name
	mockBackend.EXPECT().Storage().Return(map[string]storage.Pool{}).AnyTimes() // Always return an empty storage list
	mockBackend.EXPECT().HasVolumes().Return(false).AnyTimes()                  // Always return no volumes
	// Set the backend behavior we do care about for this testcase
	mockBackend.EXPECT().SetState(storage.Deleting) // The backend should be set to deleting
	mockBackend.EXPECT().SetOnline(false)           // The backend should be set offline
	mockBackend.EXPECT().Terminate(gomock.Any())    // The backend should be terminated

	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	// Set the store client behavior we don't care about for this testcase
	mockStoreClient.EXPECT().GetVolumeTransactions(gomock.Any()).Return([]*storage.VolumeTransaction{}, nil).AnyTimes()
	// Set the store client behavior we do care about for this testcase
	mockStoreClient.EXPECT().DeleteBackend(gomock.Any(), mockBackend).Return(nil)

	// Create an instance of the orchestrator for this test
	orchestrator := getOrchestrator(t, false)
	// Add the mocked objects to the orchestrator
	orchestrator.storeClient = mockStoreClient
	orchestrator.backends[backendUUID] = mockBackend

	// Perform the test
	err := orchestrator.DeleteBackend(ctx(), backendName)

	// Verify the results
	assert.Nilf(t, err, "Error getting backend; %v", err)
	_, ok := orchestrator.backends[backendUUID]
	assert.False(t, ok, "Backend was not properly deleted")
}

func TestPublishVolumeFailedToUpdatePersistentStore(t *testing.T) {
	config.CurrentDriverContext = config.ContextCSI
	defer func() { config.CurrentDriverContext = "" }()

	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)

	// Set fake values
	backendUUID := "1234"
	expectedError := fmt.Errorf("failure")

	// Create mocked backend that returns the expected object
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	// Create a mocked persistent store client
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(expectedError)

	// Create an instance of the orchestrator
	orchestrator := getOrchestrator(t, false)
	orchestrator.storeClient = mockStoreClient
	// Add the mocked backend to the orchestrator
	orchestrator.backends[backendUUID] = mockBackend
	volConfig := tu.GenerateVolumeConfig("fake-volume", 1, "fast", config.File)
	orchestrator.volumes["fake-volume"] = &storage.Volume{BackendUUID: backendUUID, Config: volConfig}

	// Run the test
	err := orchestrator.PublishVolume(ctx(), "fake-volume", &utils.VolumePublishInfo{})
	assert.Error(t, err, "Unexpected success publishing volume.")
}

func TestGetCHAP(t *testing.T) {
	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)

	// Set fake values
	backendUUID := "1234"
	volumeName := "foobar"
	volume := &storage.Volume{
		BackendUUID: backendUUID,
	}
	nodeName := "foobar"
	expectedChapInfo := &utils.IscsiChapInfo{
		UseCHAP:              true,
		IscsiUsername:        "foo",
		IscsiInitiatorSecret: "bar",
		IscsiTargetUsername:  "baz",
		IscsiTargetSecret:    "biz",
	}

	// Create mocked backend that returns the expected object
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().GetChapInfo(gomock.Any(), volumeName, nodeName).Return(expectedChapInfo, nil)
	// Create an instance of the orchestrator
	orchestrator := getOrchestrator(t, false)
	// Add the mocked backend and fake volume to the orchestrator
	orchestrator.backends[backendUUID] = mockBackend
	orchestrator.volumes[volumeName] = volume
	actualChapInfo, err := orchestrator.GetCHAP(ctx(), volumeName, nodeName)
	assert.Nil(t, err, "Unexpected error")
	assert.Equal(t, expectedChapInfo, actualChapInfo, "Unexpected chap info returned.")
}

func TestGetCHAPFailure(t *testing.T) {
	// Boilerplate mocking code
	mockCtrl := gomock.NewController(t)

	// Set fake values
	backendUUID := "1234"
	volumeName := "foobar"
	volume := &storage.Volume{
		BackendUUID: backendUUID,
	}
	nodeName := "foobar"
	expectedError := fmt.Errorf("some error")

	// Create mocked backend that returns the expected object
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().GetChapInfo(gomock.Any(), volumeName, nodeName).Return(nil, expectedError)
	// Create an instance of the orchestrator
	orchestrator := getOrchestrator(t, false)
	// Add the mocked backend and fake volume to the orchestrator
	orchestrator.backends[backendUUID] = mockBackend
	orchestrator.volumes[volumeName] = volume
	actualChapInfo, actualErr := orchestrator.GetCHAP(ctx(), volumeName, nodeName)
	assert.Nil(t, actualChapInfo, "Unexpected CHAP info")
	assert.Equal(t, expectedError, actualErr, "Unexpected error")
}

func TestPublishVolume(t *testing.T) {
	var (
		backendUUID        = "1234"
		nodeName           = "foo"
		volumeName         = "bar"
		subordinateVolName = "subvol"
		volume             = &storage.Volume{
			BackendUUID: backendUUID,
			Config:      &storage.VolumeConfig{AccessInfo: utils.VolumeAccessInfo{}},
		}
		subordinatevolume = &storage.Volume{
			BackendUUID: backendUUID,
			Config:      &storage.VolumeConfig{ShareSourceVolume: volumeName},
		}
		node = &utils.Node{Deleted: false}
	)
	tt := []struct {
		name               string
		volumeName         string
		shareSourceName    string
		subordinateVolName string
		subordinateVolumes map[string]*storage.Volume
		volumes            map[string]*storage.Volume
		subVolConfig       *storage.VolumeConfig
		nodes              map[string]*utils.Node
		pubsSynced         bool
		lastPub            time.Time
		volumeEnforceable  bool
		mocks              func(
			mockBackend *mockstorage.MockBackend,
			mockStoreClient *mockpersistentstore.MockStoreClient, volume *storage.Volume,
		)
		wantErr     assert.ErrorAssertionFunc
		pubTime     assert.ValueAssertionFunc
		pubEnforced assert.BoolAssertionFunc
		synced      assert.BoolAssertionFunc
	}{
		{
			name:              "LegacyVolumePubsNotSyncedNoPublicationsYet",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:              "LegacyVolumePubsNotSyncedTooSoon",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Second), // Ensure that "some" time has passed
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:              "LegacyVolumePubsNotSynced2MinutesPassed",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(nil)
				mockBackend.EXPECT().EnablePublishEnforcement(ctx(), volume).DoAndReturn(
					func(ctx context.Context, volume *storage.Volume) error {
						volume.Config.AccessInfo.PublishEnforcement = true
						return nil
					})
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:              "LegacyVolumePubsSynced",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        true,
			volumeEnforceable: false,
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().EnablePublishEnforcement(ctx(), volume).DoAndReturn(
					func(ctx context.Context, volume *storage.Volume) error {
						volume.Config.AccessInfo.PublishEnforcement = true
						return nil
					})
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:              "EnforcedVolumePubsNotSyncedNoPubsYet",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: true,
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.True,
			synced:      assert.False,
		},
		{
			name:              "EnforcedVolumePubsNotSyncedTooSoon",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: true,
			lastPub:           time.Now().Add(-time.Second), // Ensure that "some" time has passed
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.True,
			synced:      assert.False,
		},
		{
			name:              "EnforcedVolumePubsNotSynced2MinutesPassed",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: true,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:              "EnforcedVolumePubsSynced",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        true,
			volumeEnforceable: true,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:              "VolumeNotFound",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{},
			pubsSynced:        false,
			volumeEnforceable: false,
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:              "VolumeIsDeleting",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				volume.State = storage.VolumeStateDeleting
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:              "ErrorGettingVersion",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(nil, fmt.Errorf("some error"))
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:              "ErrorSettingVersion",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(fmt.Errorf("some error"))
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:              "ErrorEnablingEnforcement",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(nil)
				mockBackend.EXPECT().EnablePublishEnforcement(ctx(), gomock.Any()).Return(fmt.Errorf("some error"))
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.False,
			synced:      assert.True,
		},
		{
			name:              "ErrorReconcilingNodeAccessEnforcement",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(nil)
				mockBackend.EXPECT().EnablePublishEnforcement(ctx(), volume).DoAndReturn(
					func(ctx context.Context, volume *storage.Volume) error {
						volume.Config.AccessInfo.PublishEnforcement = true
						return nil
					})
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(fmt.Errorf("some error"))
				mockBackend.EXPECT().Name().Return("").AnyTimes()
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:              "ErrorPublishingVolume",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(nil)
				mockBackend.EXPECT().EnablePublishEnforcement(ctx(), volume).DoAndReturn(
					func(ctx context.Context, volume *storage.Volume) error {
						volume.Config.AccessInfo.PublishEnforcement = true
						return nil
					})
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:              "ErrorUpdatingVolume",
			volumeName:        volumeName,
			volumes:           map[string]*storage.Volume{volumeName: volume},
			nodes:             map[string]*utils.Node{nodeName: node},
			pubsSynced:        false,
			volumeEnforceable: false,
			lastPub:           time.Now().Add(-time.Minute * 2),
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				version := &config.PersistentStateVersion{}
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().GetVersion(ctx()).Return(version, nil)
				mockStoreClient.EXPECT().SetVersion(ctx(), version).Return(nil)
				mockBackend.EXPECT().EnablePublishEnforcement(ctx(), volume).DoAndReturn(
					func(ctx context.Context, volume *storage.Volume) error {
						volume.Config.AccessInfo.PublishEnforcement = true
						return nil
					})
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(fmt.Errorf("some error"))
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.True,
			synced:      assert.True,
		},
		{
			name:               "SubordinateVolumeTest",
			shareSourceName:    volumeName,
			volumeName:         subordinateVolName,
			subordinateVolumes: map[string]*storage.Volume{subordinateVolName: subordinatevolume},
			volumes:            map[string]*storage.Volume{volumeName: volume},
			subVolConfig:       tu.GenerateVolumeConfig("subvol", 1, "fakeSC", config.File),
			nodes:              map[string]*utils.Node{nodeName: node},
			pubsSynced:         false,
			volumeEnforceable:  false,
			lastPub:            time.Now().Add(-time.Second), // Ensure that "some" time has passed
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
				mockStoreClient.EXPECT().AddVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().ReconcileNodeAccess(ctx(), gomock.Any()).Return(nil)
				mockBackend.EXPECT().PublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(ctx(), volume).Return(nil)
			},
			wantErr:     assert.NoError,
			pubTime:     assert.IsIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
		{
			name:               "SubordinateVolumeTestFail",
			shareSourceName:    volumeName,
			volumeName:         subordinateVolName,
			subordinateVolumes: map[string]*storage.Volume{subordinateVolName: subordinatevolume},
			volumes:            map[string]*storage.Volume{"newsrcvol": volume},
			subVolConfig:       tu.GenerateVolumeConfig("subvol", 1, "fakeSC", config.File),
			nodes:              map[string]*utils.Node{nodeName: node},
			pubsSynced:         false,
			volumeEnforceable:  false,
			lastPub:            time.Now().Add(-time.Second), // Ensure that "some" time has passed
			mocks: func(
				mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient,
				volume *storage.Volume,
			) {
			},
			wantErr:     assert.Error,
			pubTime:     assert.IsNonIncreasing,
			pubEnforced: assert.False,
			synced:      assert.False,
		},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			config.CurrentDriverContext = config.ContextCSI
			defer func() { config.CurrentDriverContext = "" }()
			volume.State = storage.VolumeStateOnline
			// Boilerplate mocking code
			mockCtrl := gomock.NewController(t)

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

			// Create an instance of the orchestrator
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.volumes = tr.volumes
			o.subordinateVolumes = tr.subordinateVolumes
			o.nodes = tr.nodes
			o.volumePublicationsSynced = tr.pubsSynced
			o.lastVolumePublication = tr.lastPub
			volume.Config.AccessInfo.PublishEnforcement = tr.volumeEnforceable
			subordinatevolume.Config.AccessInfo.PublishEnforcement = tr.volumeEnforceable
			subordinatevolume.Config.ShareSourceVolume = tr.shareSourceName
			tr.mocks(mockBackend, mockStoreClient, volume)

			// Run the test
			err := o.publishVolume(ctx(), tr.volumeName, &utils.VolumePublishInfo{HostName: nodeName})
			if !tr.wantErr(t, err, "Unexpected Result") {
				return
			}
			if !tr.pubTime(t, []int64{tr.lastPub.UnixNano(), o.lastVolumePublication.UnixNano()}) {
				return
			}
			if !tr.pubEnforced(t, volume.Config.AccessInfo.PublishEnforcement) {
				return
			}
			if !tr.synced(t, o.volumePublicationsSynced) {
				return
			}
		})
	}
}

func TestPublishVolume_DirtyPublication(t *testing.T) {
	config.CurrentDriverContext = config.ContextCSI
	defer func() { config.CurrentDriverContext = "" }()

	var (
		mockCtrl        = gomock.NewController(t)
		mockBackend     = mockstorage.NewMockBackend(mockCtrl)
		mockStoreClient = mockpersistentstore.NewMockStoreClient(mockCtrl)
		backendUUID     = "12345"
		volumeName      = "foo"
		nodeName        = "bar"
	)
	mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()

	node := &utils.Node{Name: nodeName}
	pub := &utils.VolumePublication{
		Name:            utils.GenerateVolumePublishName(volumeName, nodeName),
		NodeName:        nodeName,
		VolumeName:      volumeName,
		NotSafeToAttach: true,
	}
	vol := &storage.Volume{BackendUUID: backendUUID, State: storage.VolumeStateOnline}

	o := getOrchestrator(t, false)
	o.storeClient = mockStoreClient
	o.backends = map[string]storage.Backend{backendUUID: mockBackend}
	o.nodes = map[string]*utils.Node{node.Name: node}
	o.volumes = map[string]*storage.Volume{volumeName: vol}
	_ = o.volumePublications.Set(pub.VolumeName, pub.NodeName, pub)

	err := o.publishVolume(ctx(), volumeName, &utils.VolumePublishInfo{HostName: nodeName})
	assert.Error(t, err, "Unexpected success publishing dirty publication")
}

func TestUnpublishVolume(t *testing.T) {
	var (
		backendUUID     = "1234"
		nodeName        = "foo"
		otherNodeName   = "fiz"
		volumeName      = "bar"
		otherVolumeName = "baz"
		parentSubVols   = map[string]interface{}{"dummy": nil}
		volConfig       = &storage.VolumeConfig{Name: volumeName, SubordinateVolumes: parentSubVols}
		volume          = &storage.Volume{BackendUUID: backendUUID, Config: volConfig}
		subordinateVol  = &storage.Volume{Config: &storage.VolumeConfig{Name: "dummy", ShareSourceVolume: "abc"}}
		volumeNoBackend = &storage.Volume{Config: volConfig}
		node            = &utils.Node{Deleted: false}
		deletedNode     = &utils.Node{Deleted: true}
		publication     = &utils.VolumePublication{
			NodeName:   nodeName,
			VolumeName: volumeName,
		}
		dirtyPublication = &utils.VolumePublication{
			NodeName:        nodeName,
			VolumeName:      volumeName,
			NotSafeToAttach: true,
		}
		otherVolPublication = &utils.VolumePublication{
			NodeName:   nodeName,
			VolumeName: otherVolumeName,
		}
		otherNodePublication = &utils.VolumePublication{
			NodeName:   otherNodeName,
			VolumeName: volumeName,
		}
		otherNodeAndVolPublication = &utils.VolumePublication{
			NodeName:   otherNodeName,
			VolumeName: otherVolumeName,
		}
	)
	tt := []struct {
		name            string
		volumeName      string
		nodeName        string
		dirty           bool
		driverContext   config.DriverContext
		volumes         map[string]*storage.Volume
		nodes           map[string]*utils.Node
		publications    map[string]map[string]*utils.VolumePublication
		mocks           func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient)
		wantErr         assert.ErrorAssertionFunc
		wantUnpublished assert.BoolAssertionFunc
	}{
		{
			name:          "NoOtherPublications",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.False,
		},
		{
			name:          "DirtyUnpublish",
			volumeName:    volumeName,
			nodeName:      nodeName,
			dirty:         true,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				// We should not delete the publication if this is a dirty unpublish call
				mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.False,
		},
		{
			name:          "DirtyPublication",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: dirtyPublication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				// We've already unpublished so make sure we don't call it again
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
				mockStoreClient.EXPECT().UpdateVolumePublication(ctx(), gomock.Any()).Return(nil)
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.True,
		},
		{
			name:          "OtherPublications",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: deletedNode},
			publications: map[string]map[string]*utils.VolumePublication{
				volumeName: {
					nodeName: publication,
				},
				otherVolumeName: {
					nodeName: otherVolPublication,
				},
			},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.False,
		},
		{
			name:          "OtherPublicationsDifferentNode",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes: map[string]*utils.Node{
				nodeName:      node,
				otherNodeName: node,
			},
			publications: map[string]map[string]*utils.VolumePublication{
				volumeName: {
					nodeName:      publication,
					otherNodeName: otherNodePublication,
				},
			},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().DeleteVolumePublication(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.False,
		},
		{
			name:          "VolumeNotFound",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
			},
			wantErr:         assert.Error,
			wantUnpublished: assert.False,
		},
		{
			name:          "BackendNotFound",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volumeNoBackend},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
			},
			wantErr:         assert.Error,
			wantUnpublished: assert.False,
		},

		{
			name:          "PublicationNotFound_CSI",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				// There is no publication, so there is nothing to unpublish, so no calls should be made
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.False,
		},
		{
			name:          "PublicationNotFound_Docker",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextDocker,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				// There is no publication, but there might still be a volume published, so call it anyway
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr:         assert.NoError,
			wantUnpublished: assert.False,
		},
		{
			name:          "BackendUnpublishError",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr:         assert.Error,
			wantUnpublished: assert.False,
		},
		{
			name:          "UpdatePublishedVolumeError",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: dirtyPublication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockStoreClient.EXPECT().UpdateVolumePublication(ctx(), gomock.Any()).Return(errors.New("update failed"))
			},
			wantErr:         assert.Error,
			wantUnpublished: assert.False,
		},
		{
			name:          "SubordinateVolumeParentNotFound",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
			},
			wantErr:         assert.Error,
			wantUnpublished: assert.False,
		},
		{
			name:          "NodeNotFoundWarning",
			volumeName:    volumeName,
			nodeName:      nodeName,
			driverContext: config.ContextCSI,
			volumes:       map[string]*storage.Volume{volumeName: volume},
			nodes:         map[string]*utils.Node{nodeName: node},
			publications:  map[string]map[string]*utils.VolumePublication{volumeName: {nodeName: publication}, "dummy": {"dummy": otherNodeAndVolPublication}},
			mocks: func(mockBackend *mockstorage.MockBackend, mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockBackend.EXPECT().UnpublishVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().DeleteVolumePublication(ctx(), gomock.Any()).Return(errors.New("failed to delete"))
			},
			wantErr:         assert.Error,
			wantUnpublished: assert.False,
		},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			config.CurrentDriverContext = tr.driverContext
			defer func() { config.CurrentDriverContext = "" }()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend

			o.volumes = tr.volumes
			o.nodes = tr.nodes
			if len(tr.publications) != 0 {
				o.volumePublications.SetMap(tr.publications)
			}
			o.subordinateVolumes["dummy"] = subordinateVol
			if tr.name == "SubordinateVolumeParentNotFound" {
				o.subordinateVolumes[tr.volumeName] = subordinateVol
			}

			// Resetting to false as some test cases make it true
			publication.Unpublished = false
			dirtyPublication.Unpublished = false

			tr.mocks(mockBackend, mockStoreClient)

			err := o.unpublishVolume(ctx(), tr.volumeName, tr.nodeName, tr.dirty)
			if !tr.wantErr(t, err, "Unexpected Result") {
				return
			}
			pub, ok := o.volumePublications.TryGet(volumeName, nodeName)
			if ok {
				if pub != nil && !tr.wantUnpublished(t, pub.Unpublished, "Unpublished flag is incorrect") {
					return
				}
			}
		})
	}

	// Tests for Public Unpublish Volume
	config.CurrentDriverContext = config.ContextDocker
	defer func() { config.CurrentDriverContext = "" }()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()

	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

	// Create orchestrator with fake backend and initial values
	o := getOrchestrator(t, false)
	o.storeClient = mockStoreClient
	o.backends[backendUUID] = mockBackend

	o.bootstrapError = errors.New("bootstrap error")
	err := o.UnpublishVolume(ctx(), volumeName, nodeName)
	assert.Error(t, err, "bootstrap error")

	o.bootstrapError = nil
	err = o.UnpublishVolume(ctx(), volumeName, nodeName)
	assert.Error(t, err, "volume not found")
}

func TestBootstrapSubordinateVolumes(t *testing.T) {
	var (
		backendUUID       = "1234"
		subVolumeName     = "sub_abc"
		sourceVolumeName  = "source_abc"
		sourceVolConfig   = &storage.VolumeConfig{Name: sourceVolumeName}
		sourceVolume      = &storage.Volume{Config: sourceVolConfig}
		subVolConfig      = &storage.VolumeConfig{Name: subVolumeName, ShareSourceVolume: sourceVolumeName}
		subVolume         = &storage.Volume{Config: subVolConfig}
		subVolConfig_fail = &storage.VolumeConfig{Name: subVolumeName}
		subVolume_fail    = &storage.Volume{Config: subVolConfig_fail}
	)

	tests := []struct {
		name               string
		sourceVolumeName   string
		subVolumeName      string
		volumes            map[string]*storage.Volume
		subordinateVolumes map[string]*storage.Volume
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name:               "BootStrapSubordinateVolumes",
			sourceVolumeName:   sourceVolumeName,
			subVolumeName:      subVolumeName,
			volumes:            map[string]*storage.Volume{sourceVolumeName: sourceVolume},
			subordinateVolumes: map[string]*storage.Volume{subVolumeName: subVolume},
			wantErr:            assert.NoError,
		},
		{
			name:               "BootStrapSubordinateVolumesNotFound",
			sourceVolumeName:   sourceVolumeName,
			subVolumeName:      subVolumeName,
			volumes:            map[string]*storage.Volume{sourceVolumeName: sourceVolume},
			subordinateVolumes: map[string]*storage.Volume{subVolumeName: subVolume_fail},
			wantErr:            assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentDriverContext = config.ContextCSI
			defer func() { config.CurrentDriverContext = "" }()
			mockCtrl := gomock.NewController(t)

			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.subordinateVolumes = tt.subordinateVolumes
			o.volumes = tt.volumes

			err := o.bootstrapSubordinateVolumes(ctx())
			if !tt.wantErr(t, err, "Unexpected Result") {
				return
			}
		})
	}
}

func TestAddSubordinateVolume(t *testing.T) {
	backendUUID := "1234"
	tests := []struct {
		name                  string
		sourceVolumeName      string
		subVolumeName         string
		shareSourceName       string
		ImportNotManaged      bool
		Orphaned              bool
		IsMirrorDestination   bool
		SourceVolStorageClass string
		SubVolStorageClass    string
		NFSPath               string
		backendId             string
		CloneSourceVolume     string
		ImportOriginalName    string
		State                 storage.VolumeState
		sourceVolConfig       *storage.VolumeConfig
		subVolConfig          *storage.VolumeConfig
		AddVolErr             error
		wantErr               assert.ErrorAssertionFunc
		wantVolumes           *storage.VolumeExternal
	}{
		{
			name:             "AddSubordinateVolumesNotRegularVolume",
			sourceVolumeName: "fake_vol",
			subVolumeName:    "fake_vol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fake_vol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fake_vol", 1, "fakeSC", config.File),
			wantErr:          assert.Error,
		},
		{
			name:             "AddSubordinateVolumesNotAlreadySubordinate",
			sourceVolumeName: "fake_vol",
			subVolumeName:    "fake_sub_vol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			wantErr:          assert.Error,
		},
		{
			name:             "AddSubordinateVolumesSourceDoesNotExists",
			sourceVolumeName: "fake_vol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			wantErr:          assert.Error,
		},
		{
			name:             "AddSubordinateVolumesImportNotManaged",
			sourceVolumeName: "fake_vol",
			shareSourceName:  "fake_vol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			ImportNotManaged: true,
			wantErr:          assert.Error,
		},
		{
			name:             "AddSubordinateVolumesOrphaned",
			sourceVolumeName: "fake_vol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			shareSourceName:  "fake_vol",
			Orphaned:         true,
			wantErr:          assert.Error,
		},
		{
			name:             "AddSubordinateVolumesStateNotOnline",
			sourceVolumeName: "fake_vol",
			shareSourceName:  "fake_vol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:            storage.VolumeStateDeleting,
			wantErr:          assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesSCNotSame",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "sourceSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "subSC", config.File),
			State:                 storage.VolumeStateOnline,
			SourceVolStorageClass: "sourceSC",
			SubVolStorageClass:    "subSC",
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesNoBackendFound",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			backendId:             "fakebackend",
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesSourceNotNFS",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			backendId:             backendUUID,
			NFSPath:               "",
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesSubVolSizeInvalid",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 99999999999999, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			backendId:             backendUUID,
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesSrcVolSizeInvalid",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_vol", 999999999999999, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			backendId:             backendUUID,
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesSubVolSizeLarger",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 10, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			backendId:             backendUUID,
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesCloneTest",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			CloneSourceVolume:     "fakeclone",
			backendId:             backendUUID,
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesMirrorTest",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			IsMirrorDestination:   true,
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			CloneSourceVolume:     "",
			backendId:             backendUUID,
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumesImportNameTest",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			IsMirrorDestination:   false,
			ImportOriginalName:    "fakeImportName",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			CloneSourceVolume:     "",
			backendId:             backendUUID,
			wantErr:               assert.Error,
		},
		{
			name:                  "AddSubordinateVolumes",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			IsMirrorDestination:   false,
			ImportOriginalName:    "",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			CloneSourceVolume:     "",
			backendId:             backendUUID,
			AddVolErr:             nil,
			wantErr:               assert.NoError,
		},
		{
			name:                  "AddSubordinateVolumesFail",
			sourceVolumeName:      "fake_vol",
			shareSourceName:       "fake_vol",
			sourceVolConfig:       tu.GenerateVolumeConfig("fake_Source_vol", 1, "fakeSC", config.File),
			subVolConfig:          tu.GenerateVolumeConfig("fake_sub_vol", 1, "fakeSC", config.File),
			State:                 storage.VolumeStateOnline,
			NFSPath:               "fakepath",
			IsMirrorDestination:   false,
			ImportOriginalName:    "",
			SourceVolStorageClass: "fakeSC",
			SubVolStorageClass:    "fakeSC",
			CloneSourceVolume:     "",
			backendId:             backendUUID,
			wantErr:               assert.Error,
			AddVolErr:             errors.New("failed to add volume"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentDriverContext = config.ContextCSI
			defer func() { config.CurrentDriverContext = "" }()
			mockCtrl := gomock.NewController(t)

			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetProtocol(gomock.Any()).Return(config.File).AnyTimes()
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().AddVolume(gomock.Any(), gomock.Any()).AnyTimes().Return(tt.AddVolErr).AnyTimes()

			sourceVolume := &storage.Volume{Config: tt.sourceVolConfig}
			subVolume := &storage.Volume{Config: tt.subVolConfig}
			tt.subVolConfig.ShareSourceVolume = tt.shareSourceName
			tt.sourceVolConfig.ImportNotManaged = tt.ImportNotManaged
			tt.sourceVolConfig.StorageClass = tt.SourceVolStorageClass
			tt.sourceVolConfig.AccessInfo.NfsPath = tt.NFSPath
			tt.subVolConfig.StorageClass = tt.SubVolStorageClass
			tt.subVolConfig.CloneSourceVolume = tt.CloneSourceVolume
			tt.subVolConfig.IsMirrorDestination = tt.IsMirrorDestination
			tt.subVolConfig.ImportOriginalName = tt.ImportOriginalName
			sourceVolume.Orphaned = tt.Orphaned
			sourceVolume.State = tt.State
			subVolume.State = storage.VolumeStateSubordinate
			subVolume.BackendUUID = tt.backendId
			sourceVolume.BackendUUID = tt.backendId

			volumes := map[string]*storage.Volume{tt.sourceVolumeName: sourceVolume}
			subordinateVolumes := map[string]*storage.Volume{tt.subVolumeName: subVolume}

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.subordinateVolumes = subordinateVolumes
			o.volumes = volumes

			gotVolumes, err := o.addSubordinateVolume(ctx(), tt.subVolConfig)
			if err == nil {
				tt.wantVolumes = subVolume.ConstructExternal()
			} else {
				tt.wantVolumes = nil
			}
			if !tt.wantErr(t, err, "Unexpected Result") {
				return
			}
			if !reflect.DeepEqual(gotVolumes, tt.wantVolumes) {
				t.Errorf("TridentOrchestrator.ListSubordinateVolumes() = %v, expected %v", gotVolumes, tt.wantVolumes)
			}
		})
	}
}

func TestListSubordinateVolumes(t *testing.T) {
	backendUUID := "1234"
	tests := []struct {
		name                string
		sourceVolumeName    string
		sourceVolConfigName string
		wrongSrcVolName     string
		subVolumeName       string
		bootstrapError      error
		sourceVolConfig     *storage.VolumeConfig
		subVolConfig        *storage.VolumeConfig
		SubordinateVolumes  map[string]interface{}
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name:                "ListSubordinateVolumesError",
			sourceVolumeName:    "fakeSrcVol",
			sourceVolConfigName: "fakeSrcVol",
			subVolumeName:       "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:        tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      fmt.Errorf("fake error"),
			wantErr:             assert.Error,
		},
		{
			name:                "ListSubordinateVolumesNoError",
			sourceVolumeName:    "fakeSrcVol",
			sourceVolConfigName: "fakeSrcVol",
			subVolumeName:       "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:        tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      nil,
			wantErr:             assert.NoError,
		},
		{
			name:                "ListSubordinateVolumesNoVolPassed",
			sourceVolumeName:    "",
			sourceVolConfigName: "",
			subVolumeName:       "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:        tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      nil,
			wantErr:             assert.NoError,
		},
		{
			name:                "ListSubordinateVolumesWrongVolPassed",
			sourceVolumeName:    "fakeSrcVol",
			sourceVolConfigName: "WrongfakeSrcVol",
			wrongSrcVolName:     "fakeWrongName",
			subVolumeName:       "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:        tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      nil,
			wantErr:             assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentDriverContext = config.ContextCSI
			defer func() { config.CurrentDriverContext = "" }()
			mockCtrl := gomock.NewController(t)

			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetProtocol(gomock.Any()).Return(config.File).AnyTimes()
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

			sourceVolume := &storage.Volume{Config: tt.sourceVolConfig}
			subVolume := &storage.Volume{Config: tt.subVolConfig}
			tt.sourceVolConfig.SubordinateVolumes = make(map[string]interface{})
			sourceVolume.Config.SubordinateVolumes[tt.subVolumeName] = nil
			volumes := map[string]*storage.Volume{tt.sourceVolumeName: sourceVolume}
			subordinateVolumes := map[string]*storage.Volume{tt.subVolumeName: subVolume}
			wantVolumes := make([]*storage.VolumeExternal, 0, len(volumes))

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.subordinateVolumes = subordinateVolumes
			o.volumes = volumes
			o.bootstrapError = tt.bootstrapError

			gotVolumes, err := o.ListSubordinateVolumes(ctx(), tt.sourceVolConfigName)
			if err == nil {
				wantVolumes = append(wantVolumes, subVolume.ConstructExternal())
			} else {
				wantVolumes = nil
			}
			if !tt.wantErr(t, err, "Unexpected Result") {
				return
			}
			if !reflect.DeepEqual(gotVolumes, wantVolumes) {
				t.Errorf("TridentOrchestrator.ListSubordinateVolumes() = %v, expected %v", gotVolumes, wantVolumes)
			}
		})
	}
}

func TestGetSubordinateSourceVolume(t *testing.T) {
	backendUUID := "1234"
	tests := []struct {
		name                string
		sourceVolumeName    string
		sourceVolConfigName string
		subordVolumeName    string
		bootstrapError      error
		sourceVolConfig     *storage.VolumeConfig
		subordVolConfig     *storage.VolumeConfig
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name:                "FindParentVolumeNoError",
			sourceVolumeName:    "fakeSrcVol",
			sourceVolConfigName: "fakeSrcVol",
			subordVolumeName:    "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:     tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      nil,
			wantErr:             assert.NoError,
		},
		{
			name:                "FindParentVolumeError",
			sourceVolumeName:    "fakeSrcVol",
			sourceVolConfigName: "fakeSrcVol",
			subordVolumeName:    "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:     tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      fmt.Errorf("fake error"),
			wantErr:             assert.Error,
		},
		{
			name:                "FindParentVolumeError_ParentVolumeNotFound",
			sourceVolumeName:    "fakeSrcVol1",
			sourceVolConfigName: "fakeSrcVol",
			subordVolumeName:    "fakeSubordinateVol",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:     tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      nil,
			wantErr:             assert.Error,
		},
		{
			name:                "FindParentVolumeError_SubordinateVolumeNotFound",
			sourceVolumeName:    "fakeSrcVol1",
			sourceVolConfigName: "fakeSrcVol",
			subordVolumeName:    "fakeSubordinateVol1",
			sourceVolConfig:     tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:     tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:      nil,
			wantErr:             assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentDriverContext = config.ContextCSI
			defer func() { config.CurrentDriverContext = "" }()
			mockCtrl := gomock.NewController(t)

			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetProtocol(gomock.Any()).Return(config.File).AnyTimes()
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

			sourceVolume := &storage.Volume{Config: tt.sourceVolConfig}
			subordVolume := &storage.Volume{Config: tt.subordVolConfig}
			subordVolume.Config.ShareSourceVolume = sourceVolume.Config.Name
			tt.sourceVolConfig.SubordinateVolumes = make(map[string]interface{})
			volumes := map[string]*storage.Volume{tt.sourceVolumeName: sourceVolume}
			subordinateVolumes := map[string]*storage.Volume{tt.subordVolumeName: subordVolume}
			if tt.subordVolumeName == "fakeSubordinateVol1" {
				subordinateVolumes = map[string]*storage.Volume{}
			}
			var wantVolume *storage.VolumeExternal

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.subordinateVolumes = subordinateVolumes
			o.volumes = volumes
			o.bootstrapError = tt.bootstrapError

			gotVolume, err := o.GetSubordinateSourceVolume(ctx(), tt.subordVolumeName)
			if err == nil {
				wantVolume = sourceVolume.ConstructExternal()
			} else {
				wantVolume = nil
			}
			if !tt.wantErr(t, err, "Unexpected Result") {
				return
			}
			if !reflect.DeepEqual(gotVolume, wantVolume) {
				t.Errorf("TridentOrchestrator.FindParentVolume() = %v, expected %v", gotVolume, wantVolume)
			}
		})
	}
}

func TestListVolumePublicationsForVolumeAndSubordinates(t *testing.T) {
	backendUUID := "1234"
	tests := []struct {
		name                   string
		sourceVolumeName       string
		listVolName            string
		subVolumeName          string
		subVolumeState         storage.VolumeState
		bootstrapError         error
		sourceVolConfig        *storage.VolumeConfig
		subVolConfig           *storage.VolumeConfig
		SubordinateVolumes     map[string]interface{}
		wantVolumePublications []*utils.VolumePublication
	}{
		{
			name:                   "TestListVolumePublicationsForVolumeAndSubordinatesWrongVolume",
			sourceVolumeName:       "fakeSrcVol",
			listVolName:            "WrongSrcVolName",
			subVolumeName:          "fakeSubordinateVol",
			subVolumeState:         storage.VolumeStateSubordinate,
			sourceVolConfig:        tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:           tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			wantVolumePublications: []*utils.VolumePublication{},
		},
		{
			name:                   "TestListVolumePublicationsForVolumeAndSubordinatesSubVol",
			sourceVolumeName:       "fakeSrcVol",
			listVolName:            "fakeSubordinateVol",
			subVolumeName:          "fakeSubordinateVol",
			subVolumeState:         storage.VolumeStateSubordinate,
			sourceVolConfig:        tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:           tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			wantVolumePublications: []*utils.VolumePublication{},
		},
		{
			name:             "TestListVolumePublicationsForVolumeAndSubordinates",
			sourceVolumeName: "fakeSrcVol",
			listVolName:      "fakeSrcVol",
			subVolumeName:    "fakeSubordinateVol",
			subVolumeState:   storage.VolumeStateSubordinate,
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subVolConfig:     tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			wantVolumePublications: []*utils.VolumePublication{
				{
					Name:       "foo1/bar1",
					NodeName:   "bar1",
					VolumeName: "fakeSrcVol",
					ReadOnly:   true,
					AccessMode: 1,
				},
				{
					Name:       "foo2/bar2",
					NodeName:   "bar2",
					VolumeName: "fakeSubordinateVol",
					ReadOnly:   true,
					AccessMode: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentDriverContext = config.ContextCSI
			defer func() { config.CurrentDriverContext = "" }()
			mockCtrl := gomock.NewController(t)

			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetProtocol(gomock.Any()).Return(config.File).AnyTimes()
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

			// create source and subordinate volumes
			sourceVolume := &storage.Volume{Config: tt.sourceVolConfig}
			subVolume := &storage.Volume{Config: tt.subVolConfig}
			subVolume.State = tt.subVolumeState
			tt.sourceVolConfig.SubordinateVolumes = make(map[string]interface{})
			sourceVolume.Config.SubordinateVolumes[tt.subVolumeName] = nil
			volumes := map[string]*storage.Volume{tt.sourceVolumeName: sourceVolume}
			subordinateVolumes := map[string]*storage.Volume{tt.subVolumeName: subVolume}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.subordinateVolumes = subordinateVolumes
			o.volumes = volumes

			// add the fake publications to orchestrator
			if len(tt.wantVolumePublications) > 0 {
				_ = o.volumePublications.Set(tt.wantVolumePublications[0].VolumeName,
					tt.wantVolumePublications[0].NodeName, tt.wantVolumePublications[0])
				_ = o.volumePublications.Set(tt.wantVolumePublications[1].VolumeName,
					tt.wantVolumePublications[1].NodeName, tt.wantVolumePublications[1])
			}

			gotVolumesPublications := o.listVolumePublicationsForVolumeAndSubordinates(ctx(), tt.listVolName)

			if !reflect.DeepEqual(gotVolumesPublications, tt.wantVolumePublications) {
				t.Errorf("gotVolumesPublications = %v, expected %v", gotVolumesPublications, tt.wantVolumePublications)
			}
		})
	}
}

func TestAddVolumeWithSubordinateVolume(t *testing.T) {
	const (
		backendName    = "fakeBackend"
		scName         = "fakeSC"
		fullVolumeName = "fakeVolume"
	)
	var wantVolume *storage.VolumeExternal
	var wantErr assert.ErrorAssertionFunc
	wantErr = assert.Error
	orchestrator := getOrchestrator(t, false)

	fullVolumeConfig := tu.GenerateVolumeConfig(
		fullVolumeName, 50, scName,
		config.File,
	)

	fullVolumeConfig.ShareSourceVolume = "fakeSourceVolume"
	Volume := &storage.Volume{Config: fullVolumeConfig}
	gotVolume, err := orchestrator.AddVolume(ctx(), fullVolumeConfig)

	if err == nil {
		wantVolume = Volume.ConstructExternal()
	} else {
		wantVolume = nil
	}
	if !wantErr(t, err, "Unexpected Result") {
		return
	}
	if !reflect.DeepEqual(gotVolume, wantVolume) {
		t.Errorf("GotVolume = %v, expected %v", gotVolume, wantVolume)
	}
	cleanup(t, orchestrator)
}

func TestAddVolumeWhenVolumeExistsAsSubordinateVolume(t *testing.T) {
	const (
		backendName    = "fakeBackend"
		scName         = "fakeSC"
		fullVolumeName = "fakeVolume"
	)
	var wantVolume *storage.VolumeExternal
	var wantErr assert.ErrorAssertionFunc
	wantErr = assert.Error
	orchestrator := getOrchestrator(t, false)

	fullVolumeConfig := tu.GenerateVolumeConfig(
		fullVolumeName, 50, scName,
		config.File,
	)

	Volume := &storage.Volume{Config: fullVolumeConfig}
	orchestrator.subordinateVolumes[fullVolumeConfig.Name] = Volume
	gotVolume, err := orchestrator.AddVolume(ctx(), fullVolumeConfig)

	if err == nil {
		wantVolume = Volume.ConstructExternal()
	} else {
		wantVolume = nil
	}
	if !wantErr(t, err, "Unexpected Result") {
		return
	}
	if !reflect.DeepEqual(gotVolume, wantVolume) {
		t.Errorf("GotVolume = %v, expected %v", gotVolume, wantVolume)
	}
	cleanup(t, orchestrator)
}

func TestResizeSubordinateVolume(t *testing.T) {
	backendUUID := "abcd"
	tests := []struct {
		name             string
		resizeVal        string
		sourceVolumeName string
		subordVolumeName string
		bootstrapError   error
		sourceVolConfig  *storage.VolumeConfig
		subordVolConfig  *storage.VolumeConfig
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "BootstrapError",
			resizeVal:        "1gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   errors.New("bootstrap error"),
			wantErr:          assert.Error,
		},
		{
			name:             "NilSubordinateVolumeError",
			resizeVal:        "1gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "NilParentVolumeError",
			resizeVal:        "1gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "UnsupportedResizeValue1",
			resizeVal:        "2xi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "UnsupportedResizeValue2",
			resizeVal:        "-2gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "UnsupportedSourceSizeValue1",
			resizeVal:        "2gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  &storage.VolumeConfig{Name: "fakeSrcVol", InternalName: "fakeSrcVol", Size: "1xi", Protocol: config.File, StorageClass: "fakeSC", SnapshotPolicy: "none", SnapshotDir: "none", UnixPermissions: "", VolumeMode: config.Filesystem},
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "UnsupportedSourceSizeValue2",
			resizeVal:        "2gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  &storage.VolumeConfig{Name: "fakeSrcVol", InternalName: "fakeSrcVol", Size: "-1gi", Protocol: config.File, StorageClass: "fakeSC", SnapshotPolicy: "none", SnapshotDir: "none", UnixPermissions: "", VolumeMode: config.Filesystem},
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "SubordinateSizeGreaterError",
			resizeVal:        "3Gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "PersistentStoreUpdateError",
			resizeVal:        "2Gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 2, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.Error,
		},
		{
			name:             "PersistentStoreUpdateSuccess",
			resizeVal:        "2Gi",
			sourceVolumeName: "fakeSrcVol",
			subordVolumeName: "fakeSubordinateVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 2, "fakeSC", config.File),
			subordVolConfig:  tu.GenerateVolumeConfig("fakeSubordinateVol", 1, "fakeSC", config.File),
			bootstrapError:   nil,
			wantErr:          assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceVolume := &storage.Volume{Config: tt.sourceVolConfig}
			subordVolume := &storage.Volume{Config: tt.subordVolConfig}
			subordVolume.Config.ShareSourceVolume = tt.sourceVolumeName
			tt.sourceVolConfig.SubordinateVolumes = make(map[string]interface{})
			volumes := map[string]*storage.Volume{tt.sourceVolumeName: sourceVolume}
			subordinateVolumes := map[string]*storage.Volume{tt.subordVolumeName: subordVolume}
			if tt.name == "NilSubordinateVolumeError" {
				subordinateVolumes = map[string]*storage.Volume{tt.subordVolumeName: nil}
			}
			if tt.name == "NilParentVolumeError" {
				volumes = map[string]*storage.Volume{tt.sourceVolumeName: nil}
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetProtocol(gomock.Any()).Return(config.File).AnyTimes()
			mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()
			mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
			mockBackend.EXPECT().Name().Return("mockBackend").AnyTimes()

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			if tt.name == "PersistentStoreUpdateError" {
				mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(errors.New("persistent store update failed"))
			} else {
				mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.subordinateVolumes = subordinateVolumes
			o.volumes = volumes
			o.bootstrapError = tt.bootstrapError

			err := o.ResizeVolume(ctx(), tt.subordVolumeName, tt.resizeVal)
			tt.wantErr(t, err, "Unexpected Result")
		})
	}
}

func TestResizeVolume(t *testing.T) {
	backendUUID := "abcd"
	volName := "fakeSrcVol"
	volConfig := tu.GenerateVolumeConfig(volName, 1, "fakeSC", config.File)
	volConfig.SubordinateVolumes = make(map[string]interface{})

	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "VolumeNotFoundError",
			wantErr: assert.Error,
		},
		{
			name:    "VolumeStateIsDeletingError",
			wantErr: assert.Error,
		},
		{
			name:    "AddVolumeTransactionError",
			wantErr: assert.Error,
		},
		{
			name:    "BackendNotFoundError",
			wantErr: assert.Error,
		},
		{
			name:    "BackendResizeVolumeError",
			wantErr: assert.Error,
		},
		{
			name:    "PersistentStoreUpdateError",
			wantErr: assert.Error,
		},
		{
			name:    "DeleteVolumeTranxError",
			wantErr: assert.Error,
		},
		{
			name:    "PersistentStoreUpdateSuccess",
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol := &storage.Volume{Config: volConfig}
			vol.Orphaned = true
			volumes := map[string]*storage.Volume{}
			if tt.name != "VolumeNotFoundError" {
				volumes[volName] = vol
			}
			if tt.name == "VolumeStateIsDeletingError" {
				vol.State = storage.VolumeStateDeleting
			}
			if tt.name != "BackendNotFoundError" {
				vol.BackendUUID = backendUUID
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetProtocol(gomock.Any()).Return(config.File).AnyTimes()
			mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()
			mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
			mockBackend.EXPECT().Name().Return("mockBackend").AnyTimes()
			if tt.name == "BackendResizeVolumeError" {
				mockBackend.EXPECT().ResizeVolume(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("unable to resize"))
			} else {
				mockBackend.EXPECT().ResizeVolume(ctx(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			}

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, nil).AnyTimes()
			if tt.name == "DeleteVolumeTranxError" {
				mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("delete failed"))
			} else {
				mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}
			if tt.name == "AddVolumeTransactionError" {
				mockStoreClient.EXPECT().AddVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to add to transaction"))
			} else {
				mockStoreClient.EXPECT().AddVolumeTransaction(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}
			if tt.name == "PersistentStoreUpdateError" {
				mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(errors.New("persistent store update failed"))
			} else {
				mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}

			// Create orchestrator with fake backend and initial values
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.volumes = volumes

			err := o.ResizeVolume(ctx(), volName, "2gi")
			tt.wantErr(t, err, "Unexpected Result")
		})
	}
}

func TestHandleFailedTranxResizeVolume(t *testing.T) {
	backendUUID := "abcd"
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "VolumeNotPresent",
			wantErr: assert.Error,
		},
		{
			name:    "BackendNotFoundError",
			wantErr: assert.NoError,
		},
		{
			name:    "ResizeVolumeSuccess",
			wantErr: assert.NoError,
		},
		{
			name:    "ResizeVolumeFail",
			wantErr: assert.Error,
		},
	}
	vc := tu.GenerateVolumeConfig("fakeVol", 1, "fakeSC", config.File)
	svt := &storage.VolumeTransaction{Op: storage.ResizeVolume, Config: vc}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol := &storage.Volume{Config: vc}
			if tt.name == "VolumeNotPresent" || tt.name == "ResizeVolumeSuccess" {
				vol.BackendUUID = backendUUID
			}
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			// Create a fake backend with UUID
			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(nil).AnyTimes()
			if tt.name == "VolumeNotPresent" {
				mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("delete failed"))
			} else {
				mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}
			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			if tt.name == "ResizeVolumeSuccess" || tt.name == "ResizeVolumeFail" {
				o.volumes = map[string]*storage.Volume{"fakeVol": vol}
			}
			err := o.handleFailedTransaction(ctx(), svt)
			tt.wantErr(t, err, "Unexpected Result")
		})
	}
}

func TestDeleteSubordinateVolume(t *testing.T) {
	backendUUID := "abcd"
	sc := "fakeSC"
	srcVolName := "fakeSrcVol"
	subVolName := "fakeSubordinateVol"

	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "VolumeNotFound",
			wantErr: assert.Error,
		},
		{
			name:    "FailedToDeleteVolumeFromPersistentStore",
			wantErr: assert.Error,
		},
		{
			name:    "SharedSourceVolumeNotFound",
			wantErr: assert.NoError,
		},
		{
			name:    "ErrorDeletingSourceVolume",
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceVolume := &storage.Volume{State: storage.VolumeStateDeleting, BackendUUID: backendUUID}
			sourceVolume.Config = tu.GenerateVolumeConfig(srcVolName, 1, sc, config.File)
			sourceVolume.Config.SubordinateVolumes = map[string]interface{}{subVolName: nil}

			subordVolume := &storage.Volume{Config: tu.GenerateVolumeConfig(subVolName, 1, sc, config.File)}
			if tt.name != "SharedSourceVolumeNotFound" {
				subordVolume.Config.ShareSourceVolume = srcVolName
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(errors.New("error")).AnyTimes()

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			if tt.name == "FailedToDeleteVolumeFromPersistentStore" {
				mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("delete failed"))
			} else {
				mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.volumes[srcVolName] = sourceVolume
			if tt.name != "VolumeNotFound" {
				o.subordinateVolumes[subVolName] = subordVolume
			}

			err := o.deleteSubordinateVolume(ctx(), subVolName)
			tt.wantErr(t, err, "Unexpected result")
		})
	}

	// Calling deleteSubordinateVolume from DeleteVolume function
	t.Run("DeleteSubVolFromDeleteVol", func(t *testing.T) {
		srcVolConfig := tu.GenerateVolumeConfig(srcVolName, 1, sc, config.File)
		srcVol := &storage.Volume{Config: srcVolConfig, BackendUUID: backendUUID}
		srcVol.Config.SubordinateVolumes = map[string]interface{}{subVolName: nil}

		subVolConfig := tu.GenerateVolumeConfig(subVolName, 1, sc, config.File)
		subVol := &storage.Volume{Config: subVolConfig}
		subVol.Config.ShareSourceVolume = srcVolName

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockBackend := mockstorage.NewMockBackend(mockCtrl)
		mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(errors.New("error")).AnyTimes()
		mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
		mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()
		mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
		mockBackend.EXPECT().Name().Return("mockBackend").AnyTimes()

		mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
		mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("delete failed"))

		o := getOrchestrator(t, false)
		o.storeClient = mockStoreClient
		o.backends[backendUUID] = mockBackend
		o.subordinateVolumes[subVolName] = subVol
		o.volumes[srcVolName] = srcVol

		err := o.DeleteVolume(ctx(), subVolName)
		assert.Error(t, err, "delete failed")
	})
}

func TestDeleteVolume(t *testing.T) {
	backendUUID := "abcd"
	srcVolName := "fakeSrcVol"
	cloneVolName := "cloneVol"
	sc := "fakeSC"

	publicDeleteVolumeTests := []struct {
		name             string
		sourceVolumeName string
		sourceVolConfig  *storage.VolumeConfig
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "FailedToAddTranx",
			sourceVolumeName: "fakeSrcVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			wantErr:          assert.Error,
		},
		{
			name:             "FailedToDeleteTranx",
			sourceVolumeName: "fakeSrcVol",
			sourceVolConfig:  tu.GenerateVolumeConfig("fakeSrcVol", 1, "fakeSC", config.File),
			wantErr:          assert.Error,
		},
	}

	for _, tt := range publicDeleteVolumeTests {
		t.Run(tt.name, func(t *testing.T) {
			sourceVolume := &storage.Volume{Config: tt.sourceVolConfig, BackendUUID: backendUUID, Orphaned: true}
			sourceVolume.State = storage.VolumeStateDeleting

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(errors.New("error")).AnyTimes()
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()
			mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
			mockBackend.EXPECT().Name().Return("mockBackend").AnyTimes()

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to delete tranx")).AnyTimes()
			if tt.name == "FailedToAddTranx" {
				mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, errors.New("falied to get tranx"))
			} else {
				mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, nil)
				mockStoreClient.EXPECT().AddVolumeTransaction(ctx(), gomock.Any()).Return(nil)
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.volumes[tt.sourceVolumeName] = sourceVolume

			err := o.DeleteVolume(ctx(), tt.sourceVolumeName)
			tt.wantErr(t, err, "Unexpected result")
		})
	}

	privateDeleteVolumeTests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "HasSubordinateVolume",
			wantErr: assert.NoError,
		},
		{
			name:    "HasSubordinateVolumeError",
			wantErr: assert.Error,
		},
		{
			name:    "BackendNilError",
			wantErr: assert.Error,
		},
		{
			name:    "DeleteFromPersistentStoreError",
			wantErr: assert.Error,
		},
		{
			name:    "DeleteClone",
			wantErr: assert.NoError,
		},
		{
			name:    "DeleteBackendError",
			wantErr: assert.Error,
		},
	}

	for _, tt := range privateDeleteVolumeTests {
		t.Run(tt.name, func(t *testing.T) {
			sourceVolume := &storage.Volume{State: storage.VolumeStateDeleting, BackendUUID: backendUUID, Orphaned: true}
			sourceVolume.Config = tu.GenerateVolumeConfig(srcVolName, 1, sc, config.File)
			sourceVolume.Config.CloneSourceVolume = cloneVolName
			if tt.name == "HasSubordinateVolume" || tt.name == "HasSubordinateVolumeError" {
				sourceVolume.Config.SubordinateVolumes = map[string]interface{}{"dummyVol": nil}
			}
			if tt.name == "BackendNilError" {
				sourceVolume.BackendUUID = "dummyBackend"
			}

			var cloneVolume *storage.Volume
			if tt.name == "DeleteClone" {
				cloneVolume = &storage.Volume{Config: tu.GenerateVolumeConfig(cloneVolName, 1, sc, config.File)}
				cloneVolume.State = storage.VolumeStateDeleting
				cloneVolume.Config.SubordinateVolumes = map[string]interface{}{"dummyVol": nil}
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()
			mockBackend.EXPECT().Name().Return("mockBackend").AnyTimes()
			if tt.name == "DeleteBackendError" {
				mockBackend.EXPECT().HasVolumes().Return(false)
				mockBackend.EXPECT().State().Return(storage.Deleting)
			} else {
				mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
			}
			if tt.name == "DeleteFromPersistentStoreError" {
				mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(&storage.NotManagedError{})
			} else {
				mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to delete tranx")).AnyTimes()
			if tt.name == "HasSubordinateVolumeError" || tt.name == "DeleteClone" {
				mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(errors.New("error updating volume"))
			} else {
				mockStoreClient.EXPECT().UpdateVolume(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}
			if tt.name == "BackendNilError" || tt.name == "DeleteFromPersistentStoreError" {
				mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("unable to delete"))
			} else {
				mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}
			if tt.name == "DeleteBackendError" {
				mockStoreClient.EXPECT().DeleteBackend(ctx(), gomock.Any()).Return(errors.New("delete backend failed"))
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.volumes[srcVolName] = sourceVolume
			if tt.name == "DeleteClone" {
				o.volumes[cloneVolName] = cloneVolume
			}

			err := o.deleteVolume(ctx(), srcVolName)
			tt.wantErr(t, err, "Unexpected result")
		})
	}
}

func TestHandleFailedDeleteVolumeError(t *testing.T) {
	backendUUID := "abcd"
	srcVolName := "fakeSrcVol"
	sc := "fakeSC"
	volConfig := tu.GenerateVolumeConfig(srcVolName, 1, sc, config.File)
	tnx := &storage.VolumeTransaction{Op: storage.DeleteVolume, Config: volConfig}
	sourceVolume := &storage.Volume{Config: volConfig, BackendUUID: backendUUID}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
	mockBackend.EXPECT().GetDriverName().Return("baz").AnyTimes()
	mockBackend.EXPECT().Name().Return("mockBackend").AnyTimes()
	mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(nil)

	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
	mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("unable to delete volume"))
	mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("unable to delete transaction"))

	o := getOrchestrator(t, false)
	o.storeClient = mockStoreClient
	o.backends[backendUUID] = mockBackend
	o.volumes[srcVolName] = sourceVolume

	err := o.handleFailedTransaction(ctx(), tnx)
	assert.Error(t, err, "failed to delete volume transaction")
}

func TestCreateSnapshotError(t *testing.T) {
	backendUUID := "abcd"
	// snap that already exists in o.snapshots
	sc1 := &storage.SnapshotConfig{Name: "snap1", VolumeName: "vol1"}
	id1 := sc1.ID()
	// snap of volume in deleting state
	sc2 := &storage.SnapshotConfig{Name: "snap2", VolumeName: "vol2"}
	vol2 := &storage.Volume{Config: &storage.VolumeConfig{Name: "vol2"}, State: storage.VolumeStateDeleting}
	// snap and vol used in other tests
	sc3 := &storage.SnapshotConfig{Name: "snap3", VolumeName: "vol3", Version: "1", InternalName: "snap3", VolumeInternalName: "snap3"}
	vol3 := &storage.Volume{Config: &storage.VolumeConfig{Name: "vol3", InternalName: "vol3"}, BackendUUID: backendUUID}
	snap3 := &storage.Snapshot{Config: sc3, Created: "1pm", SizeBytes: 100, State: storage.SnapshotStateCreating}

	tests := []struct {
		name       string
		snapConfig *storage.SnapshotConfig
		vol        *storage.Volume
	}{
		{
			name:       "SnapshotExists",
			snapConfig: sc1,
			vol:        vol2,
		},
		{
			name:       "VolumeNotFound",
			snapConfig: &storage.SnapshotConfig{Name: "dummy", VolumeName: "dummy"},
			vol:        vol2,
		},
		{
			name:       "VolumeInDeletingState",
			snapConfig: sc2,
			vol:        vol2,
		},
		{
			name:       "BackendNotFound",
			snapConfig: &storage.SnapshotConfig{Name: "snap3", VolumeName: "vol3"},
			vol:        &storage.Volume{Config: &storage.VolumeConfig{Name: "vol3"}, BackendUUID: "dummy"},
		},
		{
			name:       "SnapshotNotPossible",
			snapConfig: sc3,
			vol:        vol3,
		},
		{
			name:       "AddVolumeTranxError",
			snapConfig: sc3,
			vol:        vol3,
		},
		{
			name:       "MaxLimitError",
			snapConfig: sc3,
			vol:        vol3,
		},
		{
			name:       "CreateSnapshotError",
			snapConfig: sc3,
			vol:        vol3,
		},
		{
			name:       "AddSnapshotError",
			snapConfig: sc3,
			vol:        vol3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().Name().Return("backend").AnyTimes()
			mockBackend.EXPECT().GetDriverName().Return("driver").AnyTimes()
			mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
			if tt.name == "SnapshotNotPossible" {
				mockBackend.EXPECT().CanSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("cannot take snapshot"))
			} else {
				mockBackend.EXPECT().CanSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			}

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			if tt.name == "AddVolumeTranxError" {
				mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, errors.New("error getting transaction"))
			} else {
				mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, nil).AnyTimes()
				mockStoreClient.EXPECT().AddVolumeTransaction(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}
			if tt.name == "MaxLimitError" {
				mockBackend.EXPECT().CreateSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(nil, utils.MaxLimitReachedError("error"))
				mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to delete transaction"))
			}
			if tt.name == "CreateSnapshotError" {
				mockBackend.EXPECT().CreateSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to create snapshot"))
				mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(nil)
			}

			if tt.name == "AddSnapshotError" {
				mockBackend.EXPECT().CreateSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(snap3, nil)
				mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("cleanup error"))
				mockStoreClient.EXPECT().AddSnapshot(ctx(), gomock.Any()).Return(errors.New("failed to add snapshot"))
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.backends[backendUUID] = mockBackend
			o.snapshots[id1] = &storage.Snapshot{Config: sc1}
			o.volumes[tt.vol.Config.Name] = tt.vol

			_, err := o.CreateSnapshot(ctx(), tt.snapConfig)
			assert.Error(t, err, "unexpected error")
		})
	}
}

func TestDeleteSnapshotError(t *testing.T) {
	backendUUID := "abcd"
	volName := "vol"
	snapName := "snap"
	snapID := storage.MakeSnapshotID(volName, snapName)
	snapConfig := &storage.SnapshotConfig{Name: snapName, VolumeName: volName}

	publicTests := []struct {
		name     string
		volume   *storage.Volume
		snapshot *storage.Snapshot
	}{
		{
			name:     "SnapshotNotFound",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
		{
			name:     "VolumeNotFound",
			volume:   nil,
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
		{
			name:     "NilVolumeFound",
			volume:   nil,
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateMissingVolume},
		},
		{
			name:     "BackendNotFound",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
		{
			name:     "NilBackendFound",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateMissingBackend},
		},
		{
			name:     "AddVolumeTranxFail",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID, Orphaned: true},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
		{
			name:     "DeleteSnapshotFail",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID, Orphaned: true},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
	}

	for _, tt := range publicTests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().BackendUUID().Return(backendUUID).AnyTimes()
			mockBackend.EXPECT().Name().Return("backend").AnyTimes()
			mockBackend.EXPECT().GetDriverName().Return("driver").AnyTimes()
			mockBackend.EXPECT().State().Return(storage.Online).AnyTimes()
			if tt.name == "DeleteSnapshotFail" {
				mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("delete snapshot failed"))
			} else {
				mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			}

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to delete transaction")).AnyTimes()
			if tt.name == "NilVolumeFound" || tt.name == "NilBackendFound" {
				mockStoreClient.EXPECT().DeleteSnapshotIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("delete failed"))
			}
			if tt.name == "AddVolumeTranxFail" {
				mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, errors.New("failed to get transaction"))
			} else {
				mockStoreClient.EXPECT().GetExistingVolumeTransaction(ctx(), gomock.Any()).Return(nil, nil).AnyTimes()
				mockStoreClient.EXPECT().AddVolumeTransaction(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			if tt.name != "BackendNotFound" {
				o.backends[backendUUID] = mockBackend
			}
			if tt.name == "NilBackendFound" {
				o.backends[backendUUID] = nil
			}
			if tt.name != "SnapshotNotFound" {
				o.snapshots[snapID] = tt.snapshot
			}
			if tt.name != "VolumeNotFound" {
				o.volumes[volName] = tt.volume
			}

			err := o.DeleteSnapshot(ctx(), volName, snapName)
			assert.Error(t, err, "Unexpected error")
		})
	}

	privateTests := []struct {
		name     string
		volume   *storage.Volume
		snapshot *storage.Snapshot
	}{
		{
			name:     "VolumeNotFound2",
			volume:   nil,
			snapshot: nil,
		},
		{
			name:     "BackendNotFound2",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID},
			snapshot: nil,
		},
		{
			name:     "DeleteFromPersistentStoreFail",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
		{
			name:     "EmptyVolumeSnapshots",
			volume:   &storage.Volume{Config: &storage.VolumeConfig{Name: volName}, BackendUUID: backendUUID, State: storage.VolumeStateDeleting},
			snapshot: &storage.Snapshot{Config: &storage.SnapshotConfig{Name: snapName, VolumeName: volName}, State: storage.SnapshotStateCreating},
		},
	}

	for _, tt := range privateTests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBackend := mockstorage.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mockBackend.EXPECT().RemoveVolume(ctx(), gomock.Any()).Return(nil).AnyTimes()

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().DeleteVolumeIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("delete failed")).AnyTimes()
			if tt.name == "DeleteFromPersistentStoreFail" {
				mockStoreClient.EXPECT().DeleteSnapshotIgnoreNotFound(ctx(), gomock.Any()).Return(errors.New("delete failed"))
			} else {
				mockStoreClient.EXPECT().DeleteSnapshotIgnoreNotFound(ctx(), gomock.Any()).Return(nil).AnyTimes()
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			o.snapshots[snapID] = tt.snapshot
			if tt.name != "BackendNotFound2" {
				o.backends[backendUUID] = mockBackend
			}
			if tt.name != "VolumeNotFound2" {
				o.volumes[volName] = tt.volume
			}

			err := o.deleteSnapshot(ctx(), snapConfig)
			assert.Error(t, err, "Unexpected error")
		})
	}
}

func TestHandleFailedSnapshot(t *testing.T) {
	backendUUID := "abcd"
	snapName := "snap"
	volName := "vol"
	snapID := storage.MakeSnapshotID(volName, snapName)
	snapConfig := &storage.SnapshotConfig{Name: snapName, VolumeName: volName}
	volConfig := &storage.VolumeConfig{Name: volName}
	vt := &storage.VolumeTransaction{Config: volConfig, SnapshotConfig: snapConfig}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBackend := mockstorage.NewMockBackend(mockCtrl)
	mockBackend2 := mockstorage.NewMockBackend(mockCtrl)
	mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)

	o := getOrchestrator(t, false)
	o.storeClient = mockStoreClient
	o.snapshots[snapID] = &storage.Snapshot{Config: snapConfig}
	o.backends["xyz"] = mockBackend2
	o.backends[backendUUID] = mockBackend
	o.volumes[volName] = &storage.Volume{Config: volConfig, BackendUUID: backendUUID}

	// storage.AddSnapshot switch case tests
	vt.Op = storage.AddSnapshot

	mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delete snapshot"))
	mockBackend.EXPECT().Name().Return("abc")
	err := o.handleFailedTransaction(ctx(), vt)
	assert.Error(t, err, "Delete volume error")

	delete(o.snapshots, snapID)
	// As sequence of iteration in a map is not fixed, mockBackend2.State() may or may not get called
	mockBackend2.EXPECT().State().Return(storage.Unknown).AnyTimes()
	mockBackend.EXPECT().State().Return(storage.Online)
	mockBackend.EXPECT().Name().Return("abc")
	mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delete snapshot"))
	err = o.handleFailedTransaction(ctx(), vt)
	assert.Error(t, err, "Delete snapshot error")

	delete(o.backends, "xyz")
	mockBackend.EXPECT().State().Return(storage.Online)
	mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(nil)
	mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to delete transaction"))
	err = o.handleFailedTransaction(ctx(), vt)
	assert.Error(t, err, "Delete volume transaction error")

	// storage.DeleteSnapshot switch case tests
	vt.Op = storage.DeleteSnapshot

	o.snapshots[snapID] = &storage.Snapshot{Config: snapConfig}
	mockBackend.EXPECT().DeleteSnapshot(ctx(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delete snapshot"))
	mockBackend.EXPECT().Name().Return("abc")
	mockStoreClient.EXPECT().DeleteVolumeTransaction(ctx(), gomock.Any()).Return(errors.New("failed to delete transaction"))
	err = o.handleFailedTransaction(ctx(), vt)
	assert.Error(t, err, "Delete volume transaction error")
}

func TestUpdateBackendByBackendUUID(t *testing.T) {
	bName := "fake-backend"
	bConfig := map[string]interface{}{
		"version":           1,
		"storageDriverName": "fake",
		"backendName":       bName,
		"protocol":          config.File,
	}

	tests := []struct {
		name             string
		bootstrapErr     error
		backendName      string
		newBackendConfig map[string]interface{}
		contextValue     string
		callingConfigRef string
		mocks            func(mockStoreClient *mockpersistentstore.MockStoreClient)
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "BootstrapError",
			bootstrapErr:     errors.New("bootstrap error"),
			backendName:      bName,
			newBackendConfig: bConfig,
			mocks:            func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr:          assert.Error,
		},
		{
			name:             "BackendNotFound",
			backendName:      bName,
			newBackendConfig: bConfig,
			mocks:            func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr:          assert.Error,
		},
		{
			name:             "BackendCRError",
			backendName:      bName,
			newBackendConfig: bConfig,
			mocks:            func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr:          assert.Error,
		},
		{
			name:             "InvalidCallingConfig",
			backendName:      bName,
			newBackendConfig: bConfig,
			callingConfigRef: "test",
			mocks:            func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr:          assert.Error,
		},
		{
			name:        "BadCredentials",
			backendName: bName,
			newBackendConfig: map[string]interface{}{
				"version": 1, "storageDriverName": "fake", "backendName": bName,
				"username": "", "protocol": config.File,
			},
			contextValue:     ContextSourceCRD,
			callingConfigRef: "test",
			mocks:            func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr:          assert.Error,
		},
		{
			name:        "UpdateStoragePrefixError",
			backendName: bName,
			newBackendConfig: map[string]interface{}{
				"version": 1, "storageDriverName": "fake", "backendName": bName,
				"storagePrefix": "new", "protocol": config.File,
			},
			mocks:   func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr: assert.Error,
		},
		{
			name:        "BackendRenameWithExistingNameError",
			backendName: bName,
			newBackendConfig: map[string]interface{}{
				"version": 1, "storageDriverName": "fake",
				"backendName": "new", "protocol": config.File,
			},
			mocks:   func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr: assert.Error,
		},
		{
			name:        "BackendRenameError",
			backendName: bName,
			newBackendConfig: map[string]interface{}{
				"version": 1, "storageDriverName": "fake",
				"backendName": "new", "protocol": config.File,
			},
			mocks: func(mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockStoreClient.EXPECT().ReplaceBackendAndUpdateVolumes(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("rename error"))
			},
			wantErr: assert.Error,
		},
		{
			name:        "InvalidUpdateError",
			backendName: bName,
			newBackendConfig: map[string]interface{}{
				"version": 1, "storageDriverName": "fake",
				"backendName": bName, "protocol": config.Block,
			},
			mocks:   func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr: assert.Error,
		},
		{
			name:             "DefaultUpdateError",
			backendName:      bName,
			newBackendConfig: bConfig,
			mocks: func(mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockStoreClient.EXPECT().UpdateBackend(gomock.Any(), gomock.Any()).
					Return(errors.New("error updating backend"))
			},
			wantErr: assert.Error,
		},
		{
			name:        "UpdateVolumeAccessError",
			backendName: bName,
			newBackendConfig: map[string]interface{}{
				"version": 1, "storageDriverName": "fake",
				"backendName": bName, "volumeAccess": "1.1.1.1", "protocol": config.File,
			},
			mocks:   func(mockStoreClient *mockpersistentstore.MockStoreClient) {},
			wantErr: assert.Error,
		},
		{
			name:             "UpdateNonOrphanVolumeError",
			backendName:      bName,
			newBackendConfig: bConfig,
			mocks: func(mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockStoreClient.EXPECT().UpdateBackend(gomock.Any(), gomock.Any()).Return(nil)
				mockStoreClient.EXPECT().UpdateVolume(gomock.Any(), gomock.Any()).Return(errors.New("error updating non-orphan volume"))
			},
			wantErr: assert.Error,
		},
		{
			name:             "BackendUpdateSuccess",
			backendName:      bName,
			newBackendConfig: bConfig,
			mocks: func(mockStoreClient *mockpersistentstore.MockStoreClient) {
				mockStoreClient.EXPECT().UpdateBackend(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var oldBackend storage.Backend
			var oldBackendExt *storage.BackendExternal
			var configJSON []byte
			var backendUUID string
			var err error

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStoreClient := mockpersistentstore.NewMockStoreClient(mockCtrl)
			mockStoreClient.EXPECT().AddBackend(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			configJSON, err = json.Marshal(bConfig)
			if err != nil {
				t.Fatal("failed to unmarshal", err)
			}

			o := getOrchestrator(t, false)
			o.storeClient = mockStoreClient
			if tt.name != "BackendNotFound" {
				if oldBackendExt, err = o.AddBackend(ctx(), string(configJSON), ""); err != nil {
					t.Fatal("unable to create mock backend: ", err)
				}
				backendUUID = oldBackendExt.BackendUUID
				if tt.name == "UpdateNonOrphanVolumeError" {
					o.volumes["vol1"] = &storage.Volume{
						Config:      &storage.VolumeConfig{InternalName: "vol1"},
						BackendUUID: backendUUID, Orphaned: false,
					}
				}
			}

			o.bootstrapError = tt.bootstrapErr
			if tt.name == "BackendCRError" {
				// Use case where Backend ConfigRef is non-empty
				oldBackend, _ = o.getBackendByBackendUUID(backendUUID)
				oldBackend.SetConfigRef("test")
			}

			configJSON, err = json.Marshal(tt.newBackendConfig)
			if err != nil {
				t.Fatal("failed to unmarshal newBackendConfig", err)
			}

			if tt.name == "BackendRenameWithExistingNameError" || tt.name == "UpdateVolumeAccessWithExistingNameError" {
				// Adding backend with the same name as newBackendConfig to get this error
				if _, err = o.AddBackend(ctx(), string(configJSON), ""); err != nil {
					t.Fatal("unable to create mock backend: ", err)
				}
			}

			tt.mocks(mockStoreClient)
			c := context.WithValue(ctx(), ContextKeyRequestSource, tt.contextValue)

			_, err = o.UpdateBackendByBackendUUID(c, bName, string(configJSON), backendUUID, tt.callingConfigRef)
			tt.wantErr(t, err, "Unexpected result")
		})
	}
}
