package main

import (
	"kstation_backend/internal/repository/dbrepo"
	"os"
	"testing"
	"time"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2OTUxNjAxOTIsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.SR_1LAvyFImubpwW23wktEJX0881RO-Xy5_cF3sxFVc"

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "secretString"
	app.auth = Auth{
		Issuer: app.JWTIssuer,
		Audience: app.JWTAudience,
		Secret: app.JWTSecret,
		TokenExpiry: time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath: "/",
		CookieName: "refresh_token",
		CookieDomain: app.CookieDomain,
	}
	os.Exit(m.Run())
}