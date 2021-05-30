package webserver

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func TestViperConfig(t *testing.T) {
	InitWithConfigFile("../config.default.toml")
	defer viper.Reset()
	if !goutils.IsInitedViper() {
		t.Error("init viper failed")
	}
}

func TestRun(t *testing.T) {
	InitWithConfigFile("../config.default.toml")
	defer viper.Reset()
	viper.Set("server.mode", "release")
	viper.Set("logging.level", "error")
	app := NewGinEngine(nil)
	app.GET("/666", func(c *gin.Context) {
		c.JSON(200, 666)
	})
	go Run(app)
	time.Sleep(100 * time.Millisecond)

	rsp, err := http.Get("http://localhost" + viper.GetString("server.addr") + "/666")
	if err != nil {
		t.Fatal("request running server error:", err)
	}
	defer rsp.Body.Close()
	if b, err := ioutil.ReadAll(rsp.Body); err != nil {
		t.Error("read running server response body error:", err)
	} else if string(b) != "666" {
		t.Error("running server response invalid:", string(b))
	}
}

func TestUnixRun(t *testing.T) {
	InitWithConfigFile("../config.default.toml")
	defer viper.Reset()
	socketFilename := "/tmp/webserver_test.socket"
	defer os.Remove(socketFilename)
	viper.Set("server.addr", "unix:"+socketFilename)
	viper.Set("server.mode", "release")
	viper.Set("logging.level", "error")
	app := NewGinEngine(nil)
	app.GET("/666", func(c *gin.Context) {
		c.JSON(200, 666)
	})
	go Run(app)
	time.Sleep(100 * time.Millisecond)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, "unix", socketFilename)
			},
		},
	}
	rsp, err := client.Get("http://unix/666")
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	if b, err := ioutil.ReadAll(rsp.Body); err != nil {
		t.Error("read running server response body error:", err)
	} else if string(b) != "666" {
		t.Error("running server response invalid:", string(b))
	}
}
