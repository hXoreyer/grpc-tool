package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gRpcTool/server"
	"io/ioutil"
	"net/http"
	"os"
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

func main() {

	r := gin.Default()

	r.Use(Cors())
	r.Static("/assets", "./assets")

	r.POST("/LinkMethods", LinkMethods)
	r.POST("/MethodParam", GetMethodParam)
	r.POST("/Call", Call)
	r.POST("/set", SetEtc)
	r.POST("/get", GetEtc)

	r.Run(":10580")
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
	fmt.Println(news)

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
	fmt.Println(js)
	rd, e := ioutil.ReadAll(f)
	fmt.Println(string(rd), e)
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
	fmt.Println(js)

	c.JSON(200, js)
}
