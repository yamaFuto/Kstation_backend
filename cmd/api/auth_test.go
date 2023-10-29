package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_app_getTokenFromHeaderAndVerify(t *testing.T) {
	testUser := jwtUser {
		ID: 1,
		FirstName: "Admin",
		LastName: "user",
	}

	tokens, _ := app.auth.GenerateTokenPair(&testUser)

	var tests = []struct{
		name string
		token string
		errorExpected bool
		setHeader bool
		issuer string
	}{
		{"valid", fmt.Sprintf("Bearer %s", tokens.Token), false, true, app.Domain},
		{"valid expired", fmt.Sprintf("Bearer %s", expiredToken), true, true, app.Domain},
		{"valid expired", "", true, false, app.Domain},
		{"invalid token", fmt.Sprintf("Bearer %s1", tokens.Token), true, true, app.Domain},
		{"no bearer", fmt.Sprintf("Bear %s", tokens.Token), true, true, app.Domain},
		{"three header parts", fmt.Sprintf("Bearer %s 1", tokens.Token), true, true, app.Domain},

		//make sure the next test is the last one to run
		{"wrong issuer", fmt.Sprintf("Bearer %s", tokens.Token), true, true, "anotherdomain.com"},
	}

	for _, e := range tests {
		if e.issuer != app.Domain {
			app.auth.Issuer = e.issuer
			tokens, _ = app.auth.GenerateTokenPair(&testUser)
		}
		req, _ := http.NewRequest("GET", "/", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}

		rr := httptest.NewRecorder()

		_, _, err := app.auth.GetTokenFromHeaderAndVerify(rr, req)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: did not expect error, but got one - %s", e.name, err.Error())
		}

		if err == nil && e.errorExpected {
			t.Errorf("%s: expected error, but did not get one", e.name)
		}
		app.Domain = "example.com"
	}
}