package main

import (
    "fmt"
    // "strings"
    // "strconv"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "path/filepath"
    "github.com/gin-gonic/gin"
)

func get_lang(name string) string {
    lang := "ko-KR"
    data, err := ioutil.ReadFile(filepath.Join("lang", lang + ".json"))
    if err != nil {
		panic(err)
	}
    
    data_map := map[string]string{}
    json.Unmarshal([]byte(data), &data_map)
    
    if _, check := data_map[name]; check {
        return data_map[name]
    } else {
        return name + " (" + lang + ")"
    }
}

func main() {
    gin.SetMode(gin.ReleaseMode)
    gin_router := gin.Default()
    gin_router.LoadHTMLGlob("view/*")
    
    gin_router.GET("/", func(gin_data *gin.Context) {
		gin_data.HTML(http.StatusOK, "index.html", gin.H{
            "title": get_lang("main_page"),
            "content" : "TEST",
        },)
	})
    
    fmt.Println("Run 0.0.0.0:80")
    gin_router.Run("0.0.0.0:80")
}