package products

import (
	"github.com/user/gokata/core/tests/support"
	. "github.com/onsi/gomega"
)

func RequestSearchProducts(query string) *support.HTTPResponse {
	req := support.BuildRequest("GET", "/products/search?q="+query, nil)
	return support.ExecuteHTTPRequest(req)
}

func (s Steps) ExecuteSearchProducts(condition string) {
	ctx := support.Ctx()

	switch condition {
	case "success":
		resp := RequestSearchProducts(ctx.TDString("search_query"))
		support.ValidateStatusCode(resp, 200)

		support.ValidateJMESPathExists(resp.Body, "products")
		total, _ := support.JMESPathQueryInt(resp.Body, "total")
		Expect(total).To(BeNumerically(">", 0), "search should return results")
		ctx.Set("search_total", total)

	case "no-results":
		resp := RequestSearchProducts(ctx.TDString("search_query"))
		support.ValidateStatusCode(resp, 200)
		support.ValidateJMESPathValue(resp.Body, "total", float64(0))
	}
}
