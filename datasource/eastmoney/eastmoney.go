// 东方财富数据源封装

package eastmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// EastMoney 东方财富数据源
type EastMoney struct {
	// http 客户端
	HTTPClient *http.Client
}

// NewEastMoney 创建 EastMoney 实例
func NewEastMoney() EastMoney {
	hc := &http.Client{
		Timeout: 5 * time.Second,
	}
	return EastMoney{
		HTTPClient: hc,
	}
}

// NewMultipartReq 根据参数创建form-data请求
func NewMultipartReq(ctx context.Context, url string, reqData map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range reqData {
		if err := writer.WriteField(k, v); err != nil {
			return nil, errors.Wrap(err, "WriteField error")
		}
	}
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "Writer close error")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequestWithContext error")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

// Post 发送 multipart/form 的 post 请求
func (e EastMoney) Post(ctx context.Context, url string, reqData map[string]string, rspPointer interface{}) error {
	beginTime := time.Now()
	logging.Info(ctx, "EastMoney Post "+url+" begin", zap.Any("reqData", reqData))
	req, err := NewMultipartReq(ctx, url, reqData)
	if err != nil {
		return err
	}
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Post do request error")
	}
	defer resp.Body.Close()

	rspbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, fmt.Sprint("Post read resp body error, resp.Body:", resp.Body))
	}
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Info(ctx, "EastMoney Post "+url+" end", zap.Int64("latency(ms)", latency), zap.String("resp", string(rspbuf)))
	if err := json.Unmarshal(rspbuf, rspPointer); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Post json unmarshal result error. reqData:%#v, resp:%s", reqData, string(rspbuf)))
	}
	return nil
}
