//+build mage

package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"log"
)

const (
	policyName           = "dd-edge-policy"
	accountsTableName    = "accounts"
	thingTypeDescription = "detectordag device"
	thingTypeName        = "detectordag"
	thingGroupName       = "detectordag"
	topicRuleName        = "PowerStatusChanged"
	functionName         = "detectordag-consumer-AFE6GRIVNL4R"
)

var thingTypeProperties = []*string{aws.String("name"), aws.String("account-id")}

var iotClient *iot.IoT
var lambdaClient *lambda.Lambda

func init() {
	// Create an AWS session
	sesh, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Create an IoT client
	iotClient = iot.New(sesh)
	// Create a lambda client
	lambdaClient = lambda.New(sesh)
}

// CreatePolicy creates a policy for the edge devices
func CreatePolicy() error {
	// Read in the policy
	doc, err := ioutil.ReadFile("config/policy.json")
	if err != nil {
		return err
	}
	// Create the policy
	_, err = iotClient.CreatePolicy(&iot.CreatePolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(string(doc)),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == iot.ErrCodeResourceAlreadyExistsException {
			// The policy already exists, happy days
			return nil
		}
	}
	return err
}

// CreateRule creates a rule to fire a lambda function
func CreateRule() error {
	// Get the lambda in question
	lambda, err := lambdaClient.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return err
	}
	// Create a rule for device shadow changes
	_, err = iotClient.CreateTopicRule(&iot.CreateTopicRuleInput{
		RuleName: aws.String(topicRuleName),
		TopicRulePayload: &iot.TopicRulePayload{
			Description:      aws.String("Run a lambda function to handle power status updates"),
			AwsIotSqlVersion: aws.String("2016-03-23"),
			RuleDisabled:     aws.Bool(false),
			Sql:              aws.String("SELECT topic(3) as deviceId, timestamp, current.state.reported as state, current.metadata.reported as updated FROM '$aws/things/+/shadow/update/documents' WHERE current.state.reported.status <> previous.state.reported.status"),
			Actions: []*iot.Action{{
				Lambda: &iot.LambdaAction{
					FunctionArn: lambda.Configuration.FunctionArn,
				},
			}},
		},
	})
	return err
}

// CreateTables creates dynamoDB tables for the application
func CreateTables() error {
	// Create accounts table
	return sh.Run("aws", "dynamodb", "create-table", "--table-name", accountsTableName, "--cli-json-input", "file://db/accounts.json")
}

// CreateThingType creates the 'detectordag' thing type
// We use the thing type to predefine the attributes we want a thing to have
func CreateThingType() error {
	_, err := iotClient.CreateThingType(&iot.CreateThingTypeInput{
		ThingTypeName: aws.String(thingTypeName),
		ThingTypeProperties: &iot.ThingTypeProperties{
			SearchableAttributes: thingTypeProperties,
			ThingTypeDescription: aws.String(thingTypeDescription),
		},
	})
	return err
}

// CreateThingGroup creates the 'detectordag' thing group
// We use the thing group to apply a policy to all dags
func CreateThingGroup() error {
	// Ensure we've created the policy
	mg.Deps(CreatePolicy)
	// Create the thing group
	group, err := iotClient.CreateThingGroup(&iot.CreateThingGroupInput{
		ThingGroupName: aws.String(thingGroupName),
	})
	if err != nil {
		return err
	}
	// Attach the policy
	_, err = iotClient.AttachPolicy(&iot.AttachPolicyInput{
		PolicyName: aws.String(policyName),
		Target:     group.ThingGroupArn,
	})
	return err
}
