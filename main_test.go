package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_login(t *testing.T) {
	tests := []struct {
		name    string
		input   gofight.H
		code    int
		errorm  string
		message string
	}{
		{"invalid params", gofight.H{"username": "", "password": ""}, http.StatusBadRequest, "Parameters can't be empty", ""},
		{"wrong username and password", gofight.H{"username": "test", "password": "test"}, http.StatusUnauthorized, "Authentication failed", ""},
		{"correct username and password", gofight.H{"username": "hello", "password": "itsme"}, http.StatusOK, "", "Successfully authenticated user"},
	}
	g := gofight.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.POST("/login").
				SetForm(tt.input).
				Run(engine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
					assert.Equal(t, tt.code, r.Code)
					data := r.Body.Bytes()
					if tt.errorm != "" {
						e, _ := jsonparser.GetString(data, "error")
						assert.Equal(t, tt.errorm, e)
					}
					if tt.message != "" {
						e, _ := jsonparser.GetString(data, "message")
						assert.Equal(t, tt.message, e)
					}
				})
		})
	}
}

func Test_status(t *testing.T) {
	g := gofight.New()
	e := engine()

	g.GET("/private/me").Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})

	var cookie string
	g.POST("/login").
		SetForm(gofight.H{"username": "hello", "password": "itsme"}).
		Run(e, func(response gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, response.Code)
			r := (*httptest.ResponseRecorder)(response)
			cookie = r.Header().Get("Set-Cookie")
			assert.NotZero(t, cookie)
		})
	g.GET("/private/me").SetHeader(gofight.H{"Cookie": cookie}).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func Test_logout(t *testing.T) {
	g := gofight.New()
	e := engine()

	g.GET("/logout").Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

	var cookie string
	g.POST("/login").
		SetForm(gofight.H{"username": "hello", "password": "itsme"}).
		Run(e, func(response gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, response.Code)
			r := (*httptest.ResponseRecorder)(response)
			cookie = r.Header().Get("Set-Cookie")
			assert.NotZero(t, cookie)
		})
	g.GET("/logout").SetHeader(gofight.H{"Cookie": cookie}).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusOK, r.Code)
	})
}
