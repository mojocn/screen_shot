package main

import (
	"github.com/benbjohnson/phantomjs"
	"os"
	"net/http"
)

//https://github.com/benbjohnson/phantomjs
func main() {
	defer func() {

		err := recover()
		if err != nil {
			println(err)
		}
	}()

	var url = "https://www.163.com"
	// Start the process once.
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		panic(err)
		os.Exit(1)
	}

	defer phantomjs.DefaultProcess.Close()

	page, err := phantomjs.CreateWebPage()
	if err != nil {
		panic(err)
	}
	//set request headers
	requestHeader := http.Header{
		"User-Agent": []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
	}

	if err := page.SetCustomHeaders(requestHeader); err != nil {
		panic(err)
	}

	// Setup the viewport and render the results view.
	if err := page.SetViewportSize(640, 960); err != nil {
		panic(err)
	}
	// Open a URL.
	if err := page.Open(url); err != nil {
		panic(err)
	}
	if err := page.Render("hackernews4.png", "png", 50); err != nil {
		panic(err)
	}
}
