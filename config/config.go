package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/15
*@description:
 */
var defaultConfig *Config

func GetConfig() *Config {
	return defaultConfig
}
func SetConfig() {
	defaultConfig = new(Config)
	err := viper.Unmarshal(defaultConfig)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(err)
}

/**
	系统配置信息,字段名称一定要和配置文件的字段忽略大小写一致，
    否则会赋值不成功
*/
type Config struct {
	Server  Server  `json:"server"`
	Redis   Redis   `json:"redis"`
	Client  Client  `json:"client"`
	WeChart WeChart `json:"we_chart"`
}
type WeChart struct {
	AppID                string `json:"app_id"`
	Secret               string `json:"secret"`
	EndpointAccessToken  string `json:"endpoint_access_token"`
	EndpointCode2Session string `json:"endpoint_code2_session"`
	EndpointUserInfo     string `json:"endpoint_user_info"`
}
type Client struct {
	DialContextTimeout    time.Duration `json:"dial_context_timeout"`
	DialContextKeepalive  time.Duration `json:"dial_context_keepalive"`
	DisableKeepalives     bool          `json:"disable_keepalives"`
	DisableCompression    bool          `json:"disable_compression"`
	MaxConnsPrehost       int           `json:"max_conns_prehost"`
	MaxidleConnsPerhost   int           `json:"maxidle_conns_perhost"`
	MaxidleConns          int           `json:"maxidle_conns"`
	IdleConnTimeout       time.Duration `json:"idle_conn_timeout"`
	ResponseHeaderTimeout time.Duration `json:"response_header_timeout"`
	Timeout               time.Duration `json:"timeout"`
}
type Redis struct {
	Url      string   `json:"url"`
	PoolSize int      `json:"pool_size"`
	Consumer Consumer `json:"consumer"`
	Producet Producet `json:"producet"`
}
type Producet struct {
	Inverted Inverted `json:"inverted"`
}
type Consumer struct {
	Inverted Inverted `json:"inverted"`
}
type Inverted struct {
	Topic string `json:"topic"`
}
type Server struct {
	Debug          bool   `json:"debug"`
	Port           int32  `json:"port"`
	IpdbCity       string `json:"ipdb_city"`
	IpdbDistrictV6 string `json:"ipdb_district_v6"`
	IpdbDistrict   string `json:"ipdb_district"`
}
