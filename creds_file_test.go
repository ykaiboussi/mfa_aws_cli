package main

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEditCredFile(t *testing.T) {
	// insert fake data
	id := uuid.New().String()
	key := uuid.New().String()
	token := uuid.New().String()

	l, err := EditCredFile("test_data/credentials", "default", "default-mfa", id, key, token)
	if err != nil {
		t.Fatalf("Failed to edit file error: %s", err)
	}

	for k, value := range l {
		if strings.Contains(k, "[default-mfa]") {
			if strings.Contains(value.IDEnv, "aws_secret_access_key="+id) {
				assert.Equal(t, value.IDEnv, "aws_secret_access_key="+id)
			}
			if strings.Contains(value.KeyEnv, "aws_session_token="+key) {
				assert.Equal(t, value.KeyEnv, "aws_session_token="+key)
			}
			if strings.Contains(value.Token, "aws_session_token="+token) {
				assert.Equal(t, value.Token, "aws_session_token="+token)
			}
		}
		if strings.Contains(k, "[default]") {
			if strings.Contains(value.IDEnv, "aws_access_key_id=some_id_default") {
				assert.Equal(t, value.IDEnv, "aws_access_key_id=some_id_default")
			}
			if strings.Contains(value.KeyEnv, "aws_secret_access_key=some_key") {
				assert.Equal(t, value.KeyEnv, "aws_secret_access_key=some_key")
			}
			if strings.Contains(value.AWSRegion, "region=us-west-1") {
				assert.Equal(t, value.AWSRegion, "region=us-west-1")
			}
			if strings.Contains(value.Output, "output=json") {
				assert.Equal(t, value.Output, "output=json")
			}
		}
	}
}
