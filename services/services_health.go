// 服务健康度检查

package services

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/spf13/viper"
)

// CheckMySQL 检查 mysql 服务状态
func CheckMySQL(ctx context.Context) map[string]string {
	// 检查 mysql
	env := viper.GetString("env")
	envMySQLStatus := "ok"
	if envMySQL, err := goutils.GormMySQL(env); err != nil {
		envMySQLStatus = err.Error()
	} else if sqlDB, err := envMySQL.DB(); err != nil {
		envMySQLStatus = err.Error()
	} else if err := sqlDB.Ping(); err != nil {
		envMySQLStatus = err.Error()
	}
	return map[string]string{
		env: envMySQLStatus,
	}
}

// CheckRedis 检查 redis 服务状态
func CheckRedis(ctx context.Context) map[string]string {
	env := viper.GetString("env")
	envRedisStatus := "ok"
	if envRedis, err := goutils.RedisClient(env); err != nil {
		envRedisStatus = err.Error()
	} else if _, err := envRedis.Ping(context.TODO()).Result(); err != nil {
		envRedisStatus = err.Error()
	}
	return map[string]string{
		env: envRedisStatus,
	}

}

// CheckAtomicLevelServer 检查 logging 的 AtomicLevel server 是否正常
func CheckAtomicLevelServer(ctx context.Context) string {
	client := &http.Client{}
	url := "http://localhost" + viper.GetString("logging.atomic_level_server.addr") + viper.GetString("logging.atomic_level_server.path")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err.Error()
	}
	req.SetBasicAuth(viper.GetString("basic_auth.username"), viper.GetString("basic_auth.password"))
	req.Header.Set(string(logging.TraceIDKeyname), logging.CtxTraceID(ctx))
	rsp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	lvl, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err.Error()
	}
	type levelJSON struct {
		Level string `json:"level"`
	}
	level := levelJSON{}
	if err := json.Unmarshal(lvl, &level); err != nil {
		return err.Error()
	}
	if level.Level == "" {
		return "atomiclevel server response invalid level json."
	}
	return "ok"
}
