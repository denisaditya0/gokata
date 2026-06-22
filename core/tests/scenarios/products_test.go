package features_test

import (
	"github.com/user/gokata/core/tests/support"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Products", Label(support.Products), func() {

	support.RunScenarios([]support.S{
		{
			Name:    `Get product by ID`,
			Env:     []string{support.SIT, support.Staging},
			Project: []string{"PROJ-123"},
			Steps: func() {
				support.Hit("get-product").FromService("products").WithCondition("success")
			},
		},
		{
			Name: `Get product not found`,
			Env:  []string{support.SIT, support.Staging},
			Steps: func() {
				support.Hit("get-product").FromService("products").WithCondition("not-found")
			},
		},
		{
			Name:    `Search products`,
			Env:     []string{support.SIT, support.Staging, support.Sanity},
			Project: []string{"PROJ-456"},
			Steps: func() {
				support.Hit("search-products").FromService("products").WithCondition("success")
			},
		},
		{
			Name:    `Search products no results`,
			Env:     []string{support.SIT},
			Project: []string{"PROJ-456"},
			Steps: func() {
				support.Hit("search-products").FromService("products").WithCondition("no-results")
			},
		},
		{
			Name:    `E2E: Login then get product`,
			Env:     []string{support.SIT, support.Staging},
			Project: []string{"PROJ-123"},
			Steps: func() {
				support.Hit("login").FromService("auth").WithCondition("success")
				support.Hit("get-product").FromService("products").WithCondition("success")
			},
		},
	})
})
