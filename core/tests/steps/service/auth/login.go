package auth

import (
	"github.com/user/gokata/core/tests/support"
	. "github.com/onsi/gomega"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RequestLogin(username, password string) *support.HTTPResponse {
	payload := LoginRequest{Username: username, Password: password}
	req := support.BuildRequest("POST", "/auth/login", payload)
	return support.ExecuteHTTPRequest(req)
}

func (s Steps) ExecuteLogin(condition string) {
	ctx := support.Ctx()

	switch condition {
	case "success":
		resp := RequestLogin(ctx.TDString("username"), ctx.TDString("password"))
		support.ValidateStatusCode(resp, 200)

		token, _ := support.JMESPathQueryString(resp.Body, "accessToken")
		Expect(token).NotTo(BeEmpty(), "accessToken should not be empty")
		ctx.Set("token", token)

	case "invalid-credentials":
		resp := RequestLogin(ctx.TDString("username"), ctx.TDString("password"))
		support.ValidateStatusCode(resp, 400)
	}
}
