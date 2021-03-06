package main

import (
	"log"
	"time"

	"github.com/cainelli/opa-firewall/pkg/policies"
	nouseragent "github.com/cainelli/opa-firewall/pkg/policies/no-user-agent"
	"github.com/sirupsen/logrus"
)

func main() {

	log.Print("initializing server")

	logger := logrus.New()
	policyController := policies.New([]policies.PolicyInterface{
		nouseragent.New(logger),
	}, logger)

	for {
		select {
		case <-time.After(5 * time.Second):

			policyController.Run()
		}
	}
}
