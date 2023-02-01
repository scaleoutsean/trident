// Copyright 2022 NetApp, Inc. All Rights Reserved.

package logger

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Disable any standard log output
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestLogc(t *testing.T) {
	ctx := context.Background()
	result := Logc(ctx)
	assert.NotNil(t, result, "log entry is nil")
}

func TestLogc_ContextWithKeyAndValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), CRDControllerEvent, "add")
	result := Logc(ctx)
	assert.NotNil(t, result, "log entry is nil")
}

func TestGenerateRequestContext(t *testing.T) {
	type test struct {
		ctx           context.Context
		requestID     string
		requestSource string
	}

	tests := []test{
		{ctx: context.Background(), requestID: "crdControllerEvent", requestSource: "kubernetes"},
		{ctx: nil, requestID: "requestID", requestSource: "CRD"},
		{
			ctx: context.WithValue(context.Background(), ContextKeyRequestID, "1234"), requestID: "requestID",
			requestSource: "CRD",
		},
		{
			ctx: context.WithValue(context.Background(), ContextKeyRequestSource, "1234"), requestID: "requestSource",
			requestSource: "CRD",
		},
		{ctx: context.Background(), requestID: "", requestSource: "CRD"},
		{ctx: context.Background(), requestID: "test-id", requestSource: ""},
	}

	for _, tc := range tests {
		result := GenerateRequestContext(tc.ctx, tc.requestID, tc.requestSource)
		assert.NotNil(t, result, "context is nil")
	}
}

func TestFormatMessageForLog(t *testing.T) {
	msg := "this is a test"
	punctuated := "this is a test."
	titled := "This is a test"
	caps := "THIS IS A TEST"
	caps_punc := "THIS IS A TEST."
	no_changes_needed := "This is a test."

	expectedMsg := "This is a test."
	assertErr := "Expected the string to be titled and punctuated."

	type Args struct {
		TestString string
		Expected   string
		TestName   string
	}

	testCases := []Args{
		{msg, expectedMsg, "All lowercase, no punctuation"},
		{punctuated, expectedMsg, "All lowercase, but already has punctuation"},
		{titled, expectedMsg, "Already a title, but no punctuation"},
		{caps, expectedMsg, "All caps, no punctuation"},
		{caps_punc, expectedMsg, "All caps, with punctuation"},
		{no_changes_needed, expectedMsg, "Already titled and punctuated"},
	}

	for _, args := range testCases {
		t.Run(args.TestName, func(t *testing.T) {
			assert.Equal(t, args.Expected, FormatMessageForLog(args.TestString), assertErr)
		})
	}
}
