// Copyright 2022 NetApp, Inc. All Rights Reserved.

//go:build linux

package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/netapp/trident/mocks/mock_utils/mock_luks"
)

func TestLUKSDeviceStruct_Positive(t *testing.T) {
	execCmd = fakeExecCommand
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()
	execReturnValue = ""
	execReturnCode = 0

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Test getters
	luksDevice := LUKSDevice{mappingName: "pvc-test", rawDevicePath: "/dev/sdb"}

	assert.Equal(t, "/dev/mapper/pvc-test", luksDevice.MappedDevicePath())
	assert.Equal(t, "/dev/sdb", luksDevice.RawDevicePath())

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: device is already luks formatted
	execReturnValue = ""
	execReturnCode = 0
	// Return code of 0 means it is formatted
	isFormatted, err := luksDevice.IsLUKSFormatted(context.Background())
	assert.NoError(t, err)
	assert.True(t, isFormatted)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: device is not LUKS formatted
	execReturnValue = ""
	execReturnCode = 1
	// Return code of 0 means it is formatted
	isFormatted, err = luksDevice.IsLUKSFormatted(context.Background())
	assert.NoError(t, err)
	assert.False(t, isFormatted)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: format device with LUKS
	execReturnValue = ""
	execReturnCode = 0
	// Return code of 0 means it is formatted
	err = luksDevice.LUKSFormat(context.Background(), "passphrase")
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: open LUKS device
	execReturnValue = ""
	execReturnCode = 0
	// Return code of 0 means it opened
	err = luksDevice.Open(context.Background(), "passphrase")
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: LUKS device is open
	execReturnValue = ""
	execReturnCode = 0
	// Return code of 0 means it opened
	isOpen, err := luksDevice.IsOpen(context.Background())
	assert.NoError(t, err)
	assert.True(t, isOpen)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: close LUKS device
	execReturnValue = ""
	execReturnCode = 0
	// Return code of 0 means it opened
	err = luksDevice.Close(context.Background())
	assert.NoError(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Negative cases: Test ExitError (bad command) on running command for all LUKSDevice methods
func TestLUKSDeviceStruct_Negative_MissingDevicePath(t *testing.T) {
	execCmd = fakeExecCommand
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()

	execReturnValue = ""
	execReturnCode = 0

	luksDevice := LUKSDevice{mappingName: "pvc-test", rawDevicePath: ""}

	isFormatted, err := luksDevice.IsLUKSFormatted(context.Background())
	assert.Error(t, err)
	assert.False(t, isFormatted)

	err = luksDevice.LUKSFormat(context.Background(), "passphrase")
	assert.Error(t, err)

	err = luksDevice.Open(context.Background(), "passphrase")
	assert.Error(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Negative cases: Test ExitError (bad command) on running command for all LUKSDevice methods
func TestLUKSDeviceStruct_Negative_ExitError(t *testing.T) {
	execCmd = fakeExecCommandExitError
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()

	luksDevice := LUKSDevice{mappingName: "pvc-test", rawDevicePath: "/dev/sdb"}

	isFormatted, err := luksDevice.IsLUKSFormatted(context.Background())
	assert.Error(t, err)
	assert.False(t, isFormatted)

	err = luksDevice.LUKSFormat(context.Background(), "passphrase")
	assert.Error(t, err)

	err = luksDevice.Open(context.Background(), "passphrase")
	assert.Error(t, err)

	isOpen, err := luksDevice.IsOpen(context.Background())
	assert.Error(t, err)
	assert.False(t, isOpen)

	err = luksDevice.Close(context.Background())
	assert.Error(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Negative cases: Test non-zero exit code on running command for all LUKSDevice methods
func TestLUKSDeviceStruct_Negative_ExitCode1(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnCode = 1
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()

	luksDevice := LUKSDevice{mappingName: "pvc-test", rawDevicePath: "/dev/sdb"}

	isFormatted, err := luksDevice.IsLUKSFormatted(context.Background())
	assert.NoError(t, err)
	assert.False(t, isFormatted)

	err = luksDevice.LUKSFormat(context.Background(), "passphrase")
	assert.Error(t, err)

	err = luksDevice.Open(context.Background(), "passphrase")
	assert.Error(t, err)

	isOpen, err := luksDevice.IsOpen(context.Background())
	assert.NoError(t, err)
	assert.False(t, isOpen)

	err = luksDevice.Close(context.Background())
	assert.Error(t, err)
}

func TestEnsureLUKSDevice(t *testing.T) {
	execCmd = fakeExecCommand
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()
	execReturnValue = ""
	execReturnCode = 0
	luksDevice := LUKSDevice{mappingName: "luks-pvc-123", rawDevicePath: "/dev/sdb"}
	luksFormatted, err := luksDevice.EnsureFormattedAndOpen(context.Background(), "mysecretlukspassphrase")
	assert.Nil(t, err)
	assert.Equal(t, false, luksFormatted)
}

func TestEnsureLUKSDevice_Positive(t *testing.T) {
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommand
	execReturnValue = ""
	execReturnCode = 0

	// Reset exec command
	defer func() {
		execCmd = exec.CommandContext
	}()
	fakePassphrase := "mysecretlukspassphrase"
	mockCtrl := gomock.NewController(t)
	mockLUKSDevice := mock_luks.NewMockLUKSDeviceInterface(mockCtrl)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Test LUKS device already open
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(true, nil).Times(1)

	luksFormatted, err := ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Nil(t, err)
	assert.Equal(t, false, luksFormatted)
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Test device is luks but not open
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(true, nil).Times(1)
	mockLUKSDevice.EXPECT().Open(gomock.Any(), fakePassphrase).Return(nil).Times(1)
	execCmd = fakeExecCommandPaddedOutput
	fakeData := ""
	// Return values for isDeviceUnformatted calls
	execReturnValue = string(fakeData)
	execReturnCode = 0
	execPadding = 2097152

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Nil(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Test device already has data
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().RawDevicePath().Return("/dev/sdb").Times(1)
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommandPaddedOutput
	// set non-zero bytes
	execReturnValue = "a"
	execReturnCode = 0
	execPadding = 2097152

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.NotNil(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Test device too small
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().RawDevicePath().Return("/dev/sdb").Times(1)
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommandPaddedOutput
	// set non-zero bytes
	execReturnValue = "a"
	execReturnCode = 0
	execPadding = 2097

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.NotNil(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Test device is empty
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().RawDevicePath().Return("/dev/sdb").Times(1)
	mockLUKSDevice.EXPECT().LUKSFormat(gomock.Any(), fakePassphrase).Return(nil).Times(1)
	mockLUKSDevice.EXPECT().Open(gomock.Any(), fakePassphrase).Return(nil).Times(1)
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommandPaddedOutput
	execReturnValue = ""
	execReturnCode = 0
	execPadding = 2097152

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	t.Logf("%v", err)
	assert.Nil(t, err)
	assert.Equal(t, true, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()
}

func TestEnsureLUKSDevice_Negative(t *testing.T) {
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()
	fakePassphrase := "mysecretlukspassphrase"
	mockCtrl := gomock.NewController(t)
	mockLUKSDevice := mock_luks.NewMockLUKSDeviceInterface(mockCtrl)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Test device already has data
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().RawDevicePath().Return("/dev/sdb").Times(1)
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommandPaddedOutput
	// set non-zero bytes
	execReturnValue = "a"
	execReturnCode = 0
	execPadding = 2097152

	luksFormatted, err := ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.NotNil(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Open with incorrect passphrase
	fakeError := fmt.Errorf("wrong passphrase")
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(true, nil).Times(1)
	mockLUKSDevice.EXPECT().Open(gomock.Any(), fakePassphrase).Return(fakeError).Times(1)

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Error(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Cannot check if device is already open
	fakeError = fmt.Errorf("error")
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, fakeError).Times(1)

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Error(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Cannot check if device is already luks formatted
	fakeError = fmt.Errorf("error")
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(true, fakeError).Times(1)

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Error(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Cannot check if device is formatted
	fakeError = fmt.Errorf("error")
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().RawDevicePath().Return("/dev/sdb").Times(1)
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommandPaddedOutput
	// set non-zero bytes
	execReturnValue = "fake error"
	execReturnCode = 1

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Error(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Fail to LUKS format device
	fakeError = fmt.Errorf("error")
	mockCtrl = gomock.NewController(t)
	mockLUKSDevice = mock_luks.NewMockLUKSDeviceInterface(mockCtrl)
	mockLUKSDevice.EXPECT().IsOpen(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().IsLUKSFormatted(gomock.Any()).Return(false, nil).Times(1)
	mockLUKSDevice.EXPECT().RawDevicePath().Return("/dev/sdb").Times(1)
	mockLUKSDevice.EXPECT().LUKSFormat(gomock.Any(), fakePassphrase).Return(fakeError).Times(1)
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommandPaddedOutput
	// set non-zero bytes
	execReturnValue = ""
	execReturnCode = 0
	execPadding = 2097152

	luksFormatted, err = ensureLUKSDevice(context.Background(), mockLUKSDevice, fakePassphrase)
	assert.Error(t, err)
	assert.Equal(t, false, luksFormatted)

	execCmd = exec.CommandContext
	mockCtrl.Finish()
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Negative cases: Test ExitError on running command for all LUKSDevice methods
func TestEnsureLUKSDeviceClosed_Negative(t *testing.T) {
	// Return values for isDeviceUnformatted calls
	execCmd = fakeExecCommand
	execReturnValue = ""
	execReturnCode = 0

	// Reset exec command and osFs after tests
	defer func() {
		execCmd = exec.CommandContext
		osFs = afero.NewOsFs()
	}()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Test device does not exist
	// Use a mem map fs to ensure the file does not exist
	osFs = afero.NewMemMapFs()
	err := EnsureLUKSDeviceClosed(context.Background(), "/dev/mapper/luks-test-dev")
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Test failure to stat file (filename too long)
	// afero MemMapFs does normalization of the filename, so we need to actually use osFs here
	osFs = afero.NewOsFs()
	var b strings.Builder
	b.Grow(1025)
	for i := 0; i < 1025; i++ {
		b.WriteByte('a')
	}
	s := b.String()
	err = EnsureLUKSDeviceClosed(context.Background(), "/dev/mapper/"+s)
	assert.Error(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Test luksClose works
	osFs = afero.NewMemMapFs()
	osFs.Create("/dev/mapper/luks-test-dev")
	err = EnsureLUKSDeviceClosed(context.Background(), "/dev/mapper/luks-test-dev")
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Test luksClose fails
	osFs = afero.NewMemMapFs()
	osFs.Create("/dev/mapper/luks-test-dev")
	// Return values for isDeviceUnformatted calls
	execReturnValue = "error"
	execReturnCode = 1
	err = EnsureLUKSDeviceClosed(context.Background(), "/dev/mapper/luks-test-dev")
	assert.Error(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Negative cases: Test exit code 2 on running command for Open
func TestLUKSDeviceStruct_Open_BadPassphrase(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnCode = 2
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()

	luksDevice := LUKSDevice{mappingName: "pvc-test", rawDevicePath: "/dev/sdb"}
	err := luksDevice.Open(context.Background(), "passphrase")
	assert.Error(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestRotateLUKSDevicePassphrase_Positive(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnValue = ""
	execReturnCode = 0

	// Reset exec command
	defer func() {
		execCmd = exec.CommandContext
	}()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Rotate no error
	luksDeviceName := "luks-pvc-test"
	luksDevice := &LUKSDevice{"/dev/sdb", luksDeviceName}
	err := luksDevice.RotatePassphrase(context.Background(), "pvc-test", "previous", "newpassphrase")
	assert.NoError(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestRotateLUKSDevicePassphrase_Negative(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnValue = ""
	execReturnCode = 0

	// Reset exec command
	defer func() {
		execCmd = exec.CommandContext
	}()

	luksDeviceName := "luks-pvc-test"
	luksDevice := &LUKSDevice{"/dev/sdb", luksDeviceName}

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Previous passphrase is empty
	err := luksDevice.RotatePassphrase(context.Background(), "pvc-test", "", "newpassphrase")
	assert.Error(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: New passphrase is empty
	err = luksDevice.RotatePassphrase(context.Background(), "pvc-test", "previouspassphrase", "")
	assert.Error(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Command error
	execReturnCode = 4
	err = luksDevice.RotatePassphrase(context.Background(), "pvc-test", "previous", "newpassphrase")
	assert.Error(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: No device path for LUKS device
	execReturnCode = 0
	luksDeviceName = "luks-pvc-test"
	luksDevice = &LUKSDevice{"", luksDeviceName}
	err = luksDevice.RotatePassphrase(context.Background(), "pvc-test", "previous", "newpassphrase")
	assert.Error(t, err)
}

func TestGetUnderlyingDevicePathForLUKSDevice(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnCode = 0
	execReturnValue = `/dev/mapper/luks-trident_pvc_0c6202cb_be41_46b7_bea9_7f2c5c2c4a41 is active and is in use.
  type:    LUKS2
  cipher:  aes-xts-plain64
  keysize: 512 bits
  key location: keyring
  device:  /dev/mapper/3600a09807770457a795d526950374c76
  sector size:  512
  offset:  32768 sectors
  size:    2064384 sectors
  mode:    read/write`

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case
	devicePath, err := GetUnderlyingDevicePathForLUKSDevice(context.Background(), "")
	assert.NoError(t, err)
	assert.Equal(t, "/dev/mapper/3600a09807770457a795d526950374c76", devicePath)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: cannot parse
	execReturnValue = "Not good output"
	devicePath, err = GetUnderlyingDevicePathForLUKSDevice(context.Background(), "")
	assert.Error(t, err)
	assert.Equal(t, "", devicePath)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: no LUKS device
	execReturnCode = 4
	devicePath, err = GetUnderlyingDevicePathForLUKSDevice(context.Background(), "")
	assert.Error(t, err)
	assert.Equal(t, "", devicePath)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: output has "device:" but nothing else on the line
	execReturnCode = 0
	execReturnValue = `device:`
	devicePath, err = GetUnderlyingDevicePathForLUKSDevice(context.Background(), "")
	assert.Error(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: output has "device:" but extra on the line
	execReturnCode = 0
	execReturnValue = `/dev/mapper/luks-trident_pvc_0c6202cb_be41_46b7_bea9_7f2c5c2c4a41 is active and is in use.
  type:    LUKS2
  cipher:  aes-xts-plain64
  keysize: 512 bits
  key location: keyring
  device:  /dev/mapper/3600a09807770457a795d526950374c76 extra stuff on line
  sector size:  512
  offset:  32768 sectors
  size:    2064384 sectors
  mode:    read/write`
	devicePath, err = GetUnderlyingDevicePathForLUKSDevice(context.Background(), "")
	assert.Error(t, err)
}

func TestCheckPassphrase(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnCode = 0
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()
	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Correct passphrase
	luksDevice, err := NewLUKSDevice("", "test-pvc")
	assert.NoError(t, err)
	correct, err := luksDevice.CheckPassphrase(context.Background(), "passphrase")
	assert.True(t, correct)
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Not correct passphrase
	execReturnCode = 2
	luksDevice, err = NewLUKSDevice("", "test-pvc")
	assert.NoError(t, err)
	correct, err = luksDevice.CheckPassphrase(context.Background(), "passphrase")
	assert.False(t, correct)
	assert.NoError(t, err)

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: error
	execReturnCode = 4
	luksDevice, err = NewLUKSDevice("", "test-pvc")
	assert.NoError(t, err)
	correct, err = luksDevice.CheckPassphrase(context.Background(), "passphrase")
	assert.False(t, correct)
	assert.Error(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestNewLUKSDeviceFromMappingPath(t *testing.T) {
	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case
	execCmd = fakeExecCommand
	execReturnCode = 0
	// Reset exec command after tests
	defer func() {
		execCmd = exec.CommandContext
	}()
	execReturnValue = `/dev/mapper/luks-pvc-test is active and is in use.
  type:    LUKS2
  cipher:  aes-xts-plain64
  keysize: 512 bits
  key location: keyring
  device:  /dev/sdb
  sector size:  512
  offset:  32768 sectors
  size:    2064384 sectors
  mode:    read/write`
	luksDevice, err := NewLUKSDeviceFromMappingPath(context.TODO(), "/dev/mapper/luks-pvc-test", "pvc-test")
	assert.NoError(t, err)

	assert.Equal(t, luksDevice.RawDevicePath(), "/dev/sdb")
	assert.Equal(t, luksDevice.MappedDevicePath(), "/dev/mapper/luks-pvc-test")
	assert.Equal(t, luksDevice.MappedDeviceName(), "luks-pvc-test")

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: error attempting to get device path
	execReturnCode = 2
	// Reset exec command after tests
	luksDevice, err = NewLUKSDeviceFromMappingPath(context.TODO(), "/dev/mapper/luks-pvc-test", "pvc-test")
	assert.Error(t, err)
	assert.Nil(t, luksDevice)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestResize_Positive(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnValue = ""
	execReturnCode = 0

	// Reset exec command
	defer func() {
		execCmd = exec.CommandContext
	}()

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Positive case: Resize no error
	luksDeviceName := "luks-test_pvc"
	luksDevice := &LUKSDevice{"/dev/sdb", luksDeviceName}
	err := luksDevice.Resize(context.Background(), "testpassphrase")
	assert.NoError(t, err)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestResize_Negative(t *testing.T) {
	execCmd = fakeExecCommand
	execReturnValue = ""
	execReturnCode = 0

	// Reset exec command
	defer func() {
		execCmd = exec.CommandContext
	}()

	luksDeviceName := "luks-test_pvc"
	luksDevice := &LUKSDevice{"/dev/sdb", luksDeviceName}

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Bad Passphrase error
	execReturnCode = 2
	err := luksDevice.Resize(context.Background(), "testpassphrase")
	assert.Error(t, err)
	assert.True(t, IsIncorrectLUKSPassphraseError(err))

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Negative case: Misc error
	execReturnCode = 4
	err = luksDevice.Resize(context.Background(), "testpassphrase")
	assert.Error(t, err)
}
