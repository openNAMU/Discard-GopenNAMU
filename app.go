package main

import (
    "fmt"
    "os"
    // "time"
    // "reflect"
    "strings"
    "strconv"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "path/filepath"
    "database/sql"
    "html/template"
    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
)

var global_set = map[string]string{
    "lang" : "ko-KR",
    "db_name" : "data",
}

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

func get_set_in_render() map[string]string {
    set := map[string]string{}
    set["web_name"] = "Test"
    set["web_head"] = ""
    
    return set
}

func do_init_db() *sql.DB {
    db, err := sql.Open("sqlite3", global_set["db_name"])
    if err != nil {
        panic(err)
    }
    
    db.Exec("PRAGMA read_uncommitted = true")
    
    return db
}

func do_init_set() {
    do_version_check := func() string {
        version_file, err := ioutil.ReadFile("version.json")
        if err != nil {
            panic(err)
        }
        
        version_data := map[string]string{}
        json.Unmarshal([]byte(version_file), &version_data)

        version_data_now := map[string]string{}
        if _, err := os.Stat("version_now.json"); os.IsNotExist(err) {            
            version_data_now["version_update"] = ""
        } else {
            version_file_now, err := ioutil.ReadFile("version_now.json")
            if err != nil {
                panic(err)
            }
            
            json.Unmarshal([]byte(version_file_now), &version_data_now)
        }
        
        version_file_now, err := os.Create("version_now.json")
        if err != nil {
            panic(err)
        }
        defer version_file_now.Close()
        
        json_byte, _ := json.Marshal(version_data)
        json_data := string(json_byte)
        fmt.Fprintf(version_file_now, json_data)
        
        fmt.Println("openPAN version : " + version_data["version"])
        if version_data["version_update"] != version_data_now["version_update"] {
            return version_data_now["version_update"]
        } else {
            return ""
        }
    }
    
    do_set_check := func() map[string]string {
        set_data := map[string]string{}
        if _, err := os.Stat("set.json"); os.IsNotExist(err) {
            set_list := map[string]string{
                "db_name" : "",
            }
            
            set_file, err := os.Create("set.json")
            if err != nil {
                panic(err)
            }
            defer set_file.Close()
            
            for for_name_a, _ := range set_list {
                var data_var string
                
                fmt.Print("DB name (data) : ")
                fmt.Scanln(&data_var)
                
                if data_var == "" {
                    set_list[for_name_a] = "data"
                } else {
                    set_list[for_name_a] = data_var
                }
            }
            
            json_byte, _ := json.Marshal(set_list)
            json_data := string(json_byte)
            fmt.Fprintf(set_file, json_data)
        }
        
        set_file, err := ioutil.ReadFile("set.json")
        if err != nil {
            panic(err)
        }

        json.Unmarshal([]byte(set_file), &set_data)
        
        global_set["db_name"] = set_data["db_name"]
        
        return set_data
    }
    
    do_update := func(version_data string) {
        db := do_init_db()
        defer db.Close()
        
        version_data_int, _ := strconv.Atoi(version_data)
        if 2 > version_data_int {
            
        }
    }
    
    version_data := do_version_check()
    do_set_check()
    if version_data != "" { do_update(version_data) }
}

func main() {
    gin.SetMode(gin.ReleaseMode)
    gin_router := gin.Default()
    gin_router.SetFuncMap(template.FuncMap{
        "get_lang": get_lang,
    })
    gin_router.LoadHTMLGlob(filepath.Join("view", "beer", "index.html"))
    
    do_init_set()
    
    gin_router.GET("/", func(gin_data *gin.Context) {
        db := do_init_db()
        defer db.Close()
        
		gin_data.HTML(http.StatusOK, "index.html", gin.H{
            "set" : get_set_in_render(),
            "title": get_lang("main_page"),
            "content" : template.HTML("<span>TEST</span>"),
        },)
	})
    
    gin_router.GET("/view/*path", func(gin_data *gin.Context) {
        param_path := gin_data.Param("path")
        
        if _, err := os.Stat(filepath.Join("view", param_path)); os.IsNotExist(err) {
            gin_data.Data(404, "text/html; charset=utf-8", []byte(""))
        } else {
            file, err := ioutil.ReadFile(filepath.Join("view", param_path))
            if err != nil {
                panic(err)
            }
            
            param_path_split := strings.Split(param_path, "/")
            param_path_split_last := param_path_split[len(param_path_split) - 1]
            
            param_extension := strings.Split(param_path_split_last, ".")
            param_extension_last := param_extension[len(param_extension) - 1]
            
            if len(param_extension) < 2 {
                gin_data.Data(http.StatusOK, "text/html; charset=utf-8", file)
            } else if param_extension_last == "css" {
                gin_data.Data(http.StatusOK, "text/css; charset=utf-8", file)
            } else if param_extension_last == "js" {
                gin_data.Data(http.StatusOK, "text/javascript; charset=utf-8", file)
            } else {
                gin_data.Data(http.StatusOK, "text/html; charset=utf-8", file)
            }
        }
    })
    
    gin_router.NoRoute(func(gin_data *gin.Context) {
        gin_data.HTML(http.StatusOK, "index.html", gin.H{
            "set" : get_set_in_render(),
            "title": get_lang("main_page"),
            "content" : template.HTML("<span>TEST</span>"),
        },)
    })
    
    fmt.Println("Run 0.0.0.0:80")
    gin_router.Run("0.0.0.0:80")
}