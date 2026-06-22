package products

import (
	"github.com/user/gokata/core/tests/support"
	. "github.com/onsi/gomega"
)

func RequestGetProduct(id string, token string) *support.HTTPResponse {
	req := support.BuildRequest("GET", "/products/"+id, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return support.ExecuteHTTPRequest(req)
}

func (s Steps) ExecuteGetProduct(condition string) {
	ctx := support.Ctx()

	switch condition {
	case "success":
		resp := RequestGetProduct(ctx.TDString("product_id"), ctx.GetString("token"))
		support.ValidateStatusCode(resp, 200)

		title, _ := support.JMESPathQueryString(resp.Body, "title")
		Expect(title).NotTo(BeEmpty(), "product title should not be empty")
		ctx.Set("product_title", title)

	case "not-found":
		resp := RequestGetProduct(ctx.TDString("product_id"), ctx.GetString("token"))
		support.ValidateStatusCode(resp, 404)
	}
}
