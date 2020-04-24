//+build mage

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	policyName        = "dd-edge-policy"
	belanaAppName     = "detectordag-edge"
	accountsTableName = "accounts"
)

type createThingResponse struct {
	ThingName string `json:""`
	ThingArn  string `json:""`
	ThingId   string `json:""`
}

// CreatePolicy creates a policy for the edge devices
func CreatePolicy() error {
	return sh.Run("aws", "iot", "create-policy",
		"--policy-name", policyName,
		"--policy-document", "file://config/policy.json")
}

// CreateRule creates a rule to fire a lambda function
func CreateRule() error {
	return sh.Run("aws", "iot", "create-topic-rule",
		"--rule-name", "power_status_changed",
		"--topic-rule-payload", "file://config/topicRule.json")
}

func CreateTables() error {
	// Create accounts table
	return sh.Run("aws", "dynamodb", "create-table", "--table-name", accountsTableName, "--cli-json-input", "file://db/accounts.json")
}
