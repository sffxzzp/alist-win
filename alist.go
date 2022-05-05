package main

import (
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
	if conf.Password {
		pass, err := model.GetSettingByKey("password")
		if err != nil {
			log.Errorf(err.Error())
			return false
		}
		fmt.Printf("your password: %s\n", pass.Value)
		return false
	}
	server.InitIndex()
	bootstrap.InitSettings()
	bootstrap.InitAccounts()
	bootstrap.InitCache()
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
	var err error
	if conf.Conf.Scheme.Https {
		go r.RunTLS(base, conf.Conf.Scheme.CertFile, conf.Conf.Scheme.KeyFile)
	} else {
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
		},
	})
	defer w.Destroy()
	w.Navigate(base)
	w.Run()
}
