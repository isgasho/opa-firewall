package firewall

import (
	"context"
	"time"

	"github.com/cainelli/opa-firewall/pkg/iptree"
	"github.com/open-policy-agent/opa/rego"
	"github.com/sirupsen/logrus"
)

// Firewall defines the data structure used by firewall handler
type Firewall struct {
	Configuration    *Configuration
	Logger           *logrus.Logger
	PreparedEval     rego.PreparedEvalQuery
	IPTrees          IPTrees
	Policies         map[string]PolicyEvent
	PoliciesBacklog  int
	CompileInterval  time.Duration
	context          context.Context
	warmedUp         chan bool
	startedConsuming bool
}

const (
	// EventTypePatch ...
	EventTypePatch = "PATCH"
	// EventTypeFull ...
	EventTypeFull = "FULL"
)

//IPTrees is a map where first key is the policyName, the second the bucket name its value is the iptree for the given bucket
type IPTrees map[string]map[string]*iptree.IPTree

// PolicyEvent ...
type PolicyEvent struct {
	// Type can be FULL or INCR. FULL events must contain the rego policy which will be overridden during compilation.
	// INCR types can skip the rego and send JSON patches into the Data field.
	Type string `json:"type" yaml:"type"`
	// Name of the rule. This must be unique across the running packages and during
	// initialization we do checks to avoid conflicts.
	Name string `json:"name" yaml:"name"`
	// Rego contains the declarative rego policy.
	// https://www.openpolicyagent.org/docs/latest/#rego
	Rego string `json:"rego" yaml:"rego"`
	// Data contains the data used by the rego policy. It will be placed in a shared data structure between all the packages.
	// The key to access the data is the name of the package. Ex.:
	// name: "partner" ; data: {"dev_key": "a1b2c3d4e5f6d7"}:
	// allow {
	// 		data.partner.dev_key = input.header["x-dev-access-key"]
	// }
	Data interface{} `json:"data,omitempty" yaml:"data"`
	// IPBuckets are the origin data structure that we build ip binary tree. Ex.
	// {"blacklist":{"40.127.145.4":"2020-03-11T12:05:57.137118+01:00"}}
	// The blacklist is the bucket name which can be used on the rego policy. The ip as key and its value is
	// the expiration time of the IP in the binary tree. A example use case in a rego policy would be:
	// deny {
	//   ip_in_tree(input.ip, blacklist)
	// }
	IPBuckets IPBuckets `json:"ipbuckets,omitempty" yaml:"ipbuckets"`
}

// IPBuckets key is bucketName ...
type IPBuckets map[string]IPBucket

// IPBucket key is ip, value is expiration time ...
type IPBucket map[string]time.Time

// Configuration defines the configuration section for firewall handler
type Configuration struct {
	IsEnabled bool
	DryRun    bool
}

const (
	// PolicyTopicName ...
	PolicyTopicName = "firewall-policies"
	// EventsTopicName ...
	EventsTopicName = "firewall-events"
)
