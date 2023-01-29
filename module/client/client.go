package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"gitlab.com/canyinxinxi/main-server/config"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/19
*@description:
 */

var Http *pool

type pool struct {
	Client *http.Client
}

func Init(sequence int) *pool {
	Http = &pool{
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   time.Second * config.GetConfig().Client.DialContextTimeout,
					KeepAlive: time.Second * config.GetConfig().Client.DialContextKeepalive,
				}).DialContext,
				DisableKeepAlives:     config.GetConfig().Client.DisableKeepalives,
				DisableCompression:    config.GetConfig().Client.DisableCompression,
				MaxIdleConns:          config.GetConfig().Client.MaxidleConns,
				MaxIdleConnsPerHost:   config.GetConfig().Client.MaxidleConnsPerhost,
				MaxConnsPerHost:       config.GetConfig().Client.MaxConnsPrehost,
				IdleConnTimeout:       time.Second * config.GetConfig().Client.IdleConnTimeout,
				ResponseHeaderTimeout: time.Second * config.GetConfig().Client.ResponseHeaderTimeout,
			},
			Timeout: time.Microsecond * config.GetConfig().Client.Timeout,
		},
	}
	return Http
}

func Close(sequence int) {
	Http.Client.CloseIdleConnections()
}

func (this *pool) Request(req *http.Request, deadline time.Duration) (res *http.Response, byteBody []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deadline))
	defer cancel()
	res, err = this.Client.Do(req.WithContext(ctx))
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return res, nil, err
	}

	if res.StatusCode != 200 {
		return res, nil, errors.New("res code:" + strconv.Itoa(res.StatusCode))
	}
	buffer := bytes.NewBuffer(nil)
	defer buffer.Reset()
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			return res, nil, err
		}
		defer reader.Reset(reader)
		io.Copy(buffer, reader)
	default:
		buffer = bytes.NewBuffer(nil)
		io.Copy(buffer, res.Body)
	}
	byteBody = buffer.Bytes()
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			return
		}
	}
	return
}
