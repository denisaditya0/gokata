package products

import "github.com/user/gokata/core/tests/support"

type Steps struct{ support.Steps }

func init() { support.RegisterServiceSteps("products", Steps{}) }
