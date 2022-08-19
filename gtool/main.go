package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gRpcTool/server"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	PackageName string    `json:"packageName"`
	Icon        string    `json:"icon"`
	Services    []Service `json:"services"`
}

type Service struct {
	Name    string   `json:"name"`
	Icon    string   `json:"icon"`
	Methods []Method `json:"methods"`
}

type Method struct {
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	InputType  Param  `json:"inputType"`
	OutputType Param  `json:"outputType"`
}

type Param struct {
	Names []string `json:"names"`
}

/*
func portInUse(portNumber int) int {
	res := -1
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %d", portNumber)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(r) > 0 {
		pid, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			res = -1
		} else {
			res = pid
		}
	}
	return res
}

func findFreePort() int {
	var ports = make([]int, 20)
	for i := 0; i < 20; i++ {
		ports[i] = 10580 + i
	}

	for _, v := range ports {
		in := portInUse(v)
		if in == -1 {
			return v
		}
	}
	return -1
}
*/

func main() {
	/*
		port := findFreePort()
		if port == -1 {
			fmt.Println("0")
		}
		fmt.Println(port)
	*/
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()

	r.Use(Cors())
	r.Static("/assets", "./assets")

	r.POST("/LinkMethods", LinkMethods)
	r.POST("/MethodParam", GetMethodParam)
	r.POST("/Call", Call)
	r.POST("/set", SetEtc)
	r.POST("/get", GetEtc)

	if err := r.Run(":10580"); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			fmt.Println("端口占用")
		} else {
			fmt.Println("未知错误")
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func LinkMethods(c *gin.Context) {
	param := make(map[string]interface{})
	c.BindJSON(&param)
	url := param["url"].(string)
	ctx := context.Background()
	refClient, err := server.NewClient(ctx, url)
	if err != nil {
		c.String(500, err.Error())
	}
	defer refClient.Close()
	err = refClient.ListService()
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, refClient.Service)
}

func GetMethodParam(c *gin.Context) {
	param := make(map[string]interface{})
	c.BindJSON(&param)
	url := param["url"].(string)
	serviceName := param["service"].(string)
	methodName := param["method"].(string)
	ctx := context.Background()

	refClient, err := server.NewClient(ctx, url)
	if err != nil {
		c.String(500, err.Error())
	}
	defer refClient.Close()

	ret, _ := refClient.GetParams(serviceName, methodName)
	c.JSON(200, ret)

}

func Call(c *gin.Context) {
	param := make(map[string]interface{})
	c.BindJSON(&param)
	url := param["url"].(string)
	serviceName := param["service"].(string)
	methodName := param["method"].(string)
	data := param["data"].(string)
	ctx := context.Background()
	news := strings.Replace(data, "'", "\"", -1)

	refClient, err := server.NewClient(ctx, url)
	if err != nil {
		c.String(500, err.Error())
	}
	defer refClient.Close()

	ret, err := refClient.Call(serviceName, methodName, news)
	if err != nil {
		c.String(500, err.Error())
	}
	c.JSON(200, ret)
}

func SetEtc(c *gin.Context) {
	f, _ := os.OpenFile("./info.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	defer f.Close()
	js := make(map[string]interface{})
	c.BindJSON(&js)
	json.NewEncoder(f).Encode(js)

	c.String(200, "保存成功")
}

func GetEtc(c *gin.Context) {
	f, _ := os.OpenFile("./info.json", os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	js := make(map[string]interface{})
	dat := &bytes.Buffer{}
	b, _ := ioutil.ReadAll(f)
	dat.Write(b)
	json.NewDecoder(dat).Decode(&js)

	c.JSON(200, js)
}
