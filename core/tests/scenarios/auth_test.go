package features_test

import (
	"github.com/user/gokata/core/tests/support"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Auth", Label(support.Auth), func() {

	support.RunScenarios([]support.S{
		{
			Name:    `Login with valid credentials`,
			Env:     []string{support.SIT, support.Staging},
			Project: []string{"PROJ-123"},
			Steps: func() {
				support.Hit("login").FromService("auth").WithCondition("success")
			},
		},
		{
			Name:    `Login with invalid credentials`,
			Env:     []string{support.SIT, support.Staging},
			Project: []string{"PROJ-123"},
			Steps: func() {
				support.Hit("login").FromService("auth").WithCondition("invalid-credentials")
			},
		},
	})
})
