# go+phantomJS图片截取微服务
##背景
- 前端程序员不愿意(技术不行)使用canvs截取图片分享到微信朋友圈
##准备工作
- *unix系统安装[phantomJS](http://phantomjs.org/download.html)可执行文件,`phantomjs`添加到系统环境变量
- 检验phantomjs安装是否成功,`在终端中运行$:` `phantomjs`不报错,则安装成功
- 安装go package [`github.com/benbjohnson/phantomjs`](https://github.com/benbjohnson/phantomjs),主要功能方便go调用phantomJS二进制文件命令

##go调用phantomJS代码解析
- defer函数捕捉panic
    ```go
    func main() {
    	defer func(){
    
    		err:= recover()
    		if err != nil {
    			println(err)
    		}
    	}()
    	var url = "https://www.163.com"
    ```
- 创建一个在golang里面phantomJS创建一个phanomJS进程
    ```go
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		panic(err)
		os.Exit(1)
	}
	defer phantomjs.DefaultProcess.Close()
    ```
- 创建phantomJS page设置请求headers和view port
    ```go
    page, err := phantomjs.CreateWebPage()
	if err != nil {
		panic(err)
	}
	//set request headers 设置UserAgent为iPhone
	requestHeader := http.Header{
		"User-Agent" :[]string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
	}

	if err := page.SetCustomHeaders(requestHeader);err != nil {
		panic(err)
	}

	// Setup the viewport and render the results view.
	if err := page.SetViewportSize(640, 960); err != nil {
		panic(err)
	}
    ```
- 截图并输出png文件
    ```go
    // Open a URL.
	if err := page.Open(url); err != nil {
		panic(err)
	}
	if err := page.Render("hackernews4.png", "png", 50); err != nil {
		panic(err)
	}}
    ```
- 完整代码
  `main.go`
  ```go
    package main

    import (
    	"github.com/benbjohnson/phantomjs"
    	"os"
    	"net/http"
    )
    
    //https://github.com/benbjohnson/phantomjs
    func main() {
    	defer func(){
    
    		err:= recover()
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
    		"User-Agent" :[]string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
    	}
    
    	if err := page.SetCustomHeaders(requestHeader);err != nil {
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
    	if err := page.Render("163_news_screen_shot.png", "png", 50); err != nil {
    		panic(err)
    	}}
    ```
##使用go标准库创建截图微服务
- main.go完成代码
```go
    // example of HTTP server that uses the captcha package.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/benbjohnson/phantomjs"
)

const OutputDir = "screen_shot"

//ConfigJsonBody json request body.
type jsonBody struct {
	Url              string `json:"url"`
	ViewportWidth     int `json:"viewport_width"`
	ViewportHeight     int `json:"viewport_height"`
	OutputFileName     string `json:"output_file_name"`
	OutputFileExt     string  `json:"output_file_ext"`
	Quility int `json:"quility"`
	OutputUri	string `json:"output_uri"`
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

func phantomjs_screen_shot(config *jsonBody)  {
	uri := fmt.Sprintf("%s/%s", OutputDir,config.OutputFileName)
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
		"User-Agent" :[]string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
	}

	if err := page.SetCustomHeaders(requestHeader);err != nil {
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
	
	pathPrefix := "/"+OutputDir +"/"
	staticDir := "./"+OutputDir
	http.Handle(pathPrefix,http.StripPrefix(pathPrefix, http.FileServer(http.Dir(staticDir))))

	//api for create captcha
	http.HandleFunc("/api/shot", screen_shot)


	fmt.Println("Server is at localhost:1122")
	if err := http.ListenAndServe("localhost:1122", nil); err != nil {
		log.Fatal(err)
	}
}

```

- postmanAPI接口
```
POST /api/shot HTTP/1.1
Host: localhost:1122
Content-Type: application/json
Cache-Control: no-cache

{
	"url":"http://www.163.com",
	"viewport_width":480,
	"viewport_height":960,
	"output_file_name":"awesome.jpg",
	"output_file_ext":"jpg",
	"quility":90
}
- 返回参数
```JavaScript
{
    "code": 1,
    "data": {
        "url": "http://www.163.com",
        "viewport_width": 480,
        "viewport_height": 960,
        "output_file_name": "screen_shot/awesome.jpg",
        "output_file_ext": "jpg",
        "quility": 90,
        "output_uri": "screen_shot/awesome.jpg"
    },
    "msg": "success"
}
```

截图图片地址`localhost:1122/` + `output_uri`
```


###
- todo::图片上传到阿里云oss
- todo::相同url不重复截图



## [GitHub源码地址](https://github.com/mojocn/screen_shot)













