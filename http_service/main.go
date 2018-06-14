// example of HTTP server that uses the captcha package.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/benbjohnson/phantomjs"
	"log"
	"net/http"
	"os"
)

const OutputDir = "screen_shot"

//ConfigJsonBody json request body.
type jsonBody struct {
	Url            string `json:"url"`
	ViewportWidth  int    `json:"viewport_width"`
	ViewportHeight int    `json:"viewport_height"`
	OutputFileName string `json:"output_file_name"`
	OutputFileExt  string `json:"output_file_ext"`
	Quility        int    `json:"quility"`
	OutputUri      string `json:"output_uri"`
}

// base64Captcha create http handler
func screen_shot(w http.ResponseWriter, r *http.Request) {
	//parse request parameters
	//接收客户端发送来的请求参数json
	decoder := json.NewDecoder(r.Body)
	var postParameters jsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	//截图网页
	phantomjs_screen_shot(&postParameters)
	//set json response
	//设置json响应
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body := map[string]interface{}{"code": 1, "data": postParameters, "msg": "success"}
	json.NewEncoder(w).Encode(body)
}

func phantomjs_screen_shot(config *jsonBody) {
	uri := fmt.Sprintf("%s/%s", OutputDir, config.OutputFileName)
	config.OutputFileName = uri
	config.OutputUri = uri
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
	if err := page.SetViewportSize(config.ViewportWidth, config.ViewportHeight); err != nil {
		panic(err)
	}
	// Open a URL.
	if err := page.Open(config.Url); err != nil {
		panic(err)
	}
	if err := page.Render(config.OutputFileName, config.OutputFileExt, config.Quility); err != nil {
		panic(err)
	}
}

//start a net/http server
//启动golang net/http 服务器
func main() {

	pathPrefix := "/" + OutputDir + "/"
	staticDir := "./" + OutputDir
	http.Handle(pathPrefix, http.StripPrefix(pathPrefix, http.FileServer(http.Dir(staticDir))))

	//api for create captcha
	http.HandleFunc("/api/shot", screen_shot)

	fmt.Println("Server is at localhost:1122")
	if err := http.ListenAndServe("localhost:1122", nil); err != nil {
		log.Fatal(err)
	}
}
