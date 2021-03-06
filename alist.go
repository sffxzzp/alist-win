package main

import (
	"os"
	"fmt"
	"github.com/Xhofe/alist/bootstrap"
	"github.com/Xhofe/alist/conf"
	_ "github.com/Xhofe/alist/drivers"
	"github.com/Xhofe/alist/model"
	"github.com/Xhofe/alist/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/sffxzzp/go-webview2"
)

func Init() bool {
	bootstrap.InitConf()
	bootstrap.InitCron()
	bootstrap.InitModel()
	server.InitIndex()
	bootstrap.InitSettings()
	bootstrap.InitAccounts()
	bootstrap.InitCache()
	pass, _ := model.GetSettingByKey("password")
	_ = os.WriteFile("password.txt", []byte(pass.Value), 0777)
	return true
}

func main() {
	if conf.Version {
		fmt.Printf("Built At: %s\nGo Version: %s\nAuthor: %s\nCommit ID: %s\nVersion: %s\nWebVersion: %s\n",
			conf.BuiltAt, conf.GoVersion, conf.GitAuthor, conf.GitCommit, conf.GitTag, conf.WebTag)
		return
	}
	if !Init() {
		return
	}
	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	server.InitApiRouter(r)
	base := fmt.Sprintf("%s:%d", conf.Conf.Address, conf.Conf.Port)
	log.Infof("start server @ %s", base)
	var prefix string
	if conf.Conf.Scheme.Https {
		prefix = "https"
		go r.RunTLS(base, conf.Conf.Scheme.CertFile, conf.Conf.Scheme.KeyFile)
	} else {
		prefix = "http"
		go r.Run(base)
	}
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:    false,
		DataPath: "./data",
		WindowOptions: webview2.WindowOptions{
			Title:  "alist",
			Width:  1440,
			Height: 900,
			IconId: 2,
			Center: true,
		},
	})
	defer w.Destroy()
	w.Navigate(fmt.Sprintf("%s://127.0.0.1:%d", prefix, conf.Conf.Port))
	w.Run()
}
