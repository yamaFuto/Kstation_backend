package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	// "net/url"
	"strings"
	"testing"
	"time"
)

func Test_app_authenticate(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com", "password":"secret"}`, http.StatusAccepted},
		{"not json", `I'm not JSON`, http.StatusBadRequest},
		{"empty json", `{}`, http.StatusBadRequest},
		{"empty email", `{"email":""}`, http.StatusBadRequest},
		{"empty password", `{"email":"admin@example.com"}`, http.StatusBadRequest},
		{"invalid user", `{"email":"admin@someotherdomain.com", "password":"secret"}`, http.StatusBadRequest},
	}

	for _, e := range theTests {
		var reader io.Reader
		// string型にreaderを付与した状態のbuffer(byte[]型で記憶)にする
		// struct型の場合はbytesでbyte[]型にしてからnewbufferに格納するか、url.valuesの中に格納してencode()でurlタイプのstringにする
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/auth", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_app_refreshToken(t *testing.T) {
	var tests = []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{"valid", "", http.StatusOK, true},
		{"valid but not yet ready to expire", "", http.StatusTooEarly, false},
		{"expired token", expiredToken, http.StatusUnauthorized, false},
	}

	testUser := jwtUser {
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
	}

	oldRefreshTime := app.auth.RefreshExpiry

	for _, e := range tests {
		var tkn string
		if e.token == "" {
			if e.resetRefreshTime {
				app.auth.RefreshExpiry = time.Second * 1
			}
			tokens, _ := app.auth.GenerateTokenPair(&testUser)
			tkn = tokens.RefreshToken
		} else {
			tkn = e.token
		}

		req, _ := http.NewRequest("GET", "/refresh-token", strings.NewReader(""))

		req.AddCookie(&http.Cookie{
			Name: "refresh_token",
			Value: tkn,
		})

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.refreshToken)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status of %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		app.auth.RefreshExpiry = oldRefreshTime
	}
}