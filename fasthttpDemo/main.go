package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	typeOfRequest := string(ctx.Method())
	if typeOfRequest == "GET" {
		fmt.Fprintln(ctx, "Welcome Its a GET Request!")
	}

	if typeOfRequest == "POST" {
		fmt.Fprintln(ctx, "Welcome Its a POST Request!")
	}
}

func main() {
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}
