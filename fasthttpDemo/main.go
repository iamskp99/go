package main

import (
	"fmt"
	"github.com/keploy/go-sdk/integrations/kfasthttp"
	"github.com/keploy/go-sdk/keploy"
	"github.com/valyala/fasthttp"
)

func main() {
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "shashankFasthttpDEMO",
			Port: "8080",
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:8081/api",
		},
	})
	mw := kfasthttp.FastHttpMiddlware(k)
	m := func(ctx *fasthttp.RequestCtx) {
		typeOfRequest := string(ctx.Method())
		if typeOfRequest == "GET" {
			fmt.Fprintln(ctx, "Welcome, Its a GET Request!")
		}

		if typeOfRequest == "POST" {
			fmt.Fprintln(ctx, "Welcome, Its a POST Request!")
		}
	}
	fasthttp.ListenAndServe(":8080", mw(m))
}
