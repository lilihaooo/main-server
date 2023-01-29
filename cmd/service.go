/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/wallestore/framework"
	"gitlab.com/canyinxinxi/main-server/config"
	"gitlab.com/canyinxinxi/main-server/db"
	"gitlab.com/canyinxinxi/main-server/handler/router"
	"gitlab.com/canyinxinxi/main-server/module/client"
	"gitlab.com/canyinxinxi/main-server/module/consumer"
	"gitlab.com/canyinxinxi/main-server/module/gin_proxy"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"os"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		framework.SetAppName("main-server")
		framework.Init(initFunc)
		framework.Exit(exitFunc)
	},
	Run: func(cmd *cobra.Command, args []string) {
		framework.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop")
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 初始化函数
func initFunc() {

	config.SetConfig()

	if err := redis.Init(1); err != nil {
		log.Error().Err(err).Msg("init_redis_err")
		os.Exit(6)
	}
	//consumer
	if err := consumer.InitConsumerIndex(config.GetConfig().Redis.Consumer.Inverted.Topic); err != nil {
		log.Error().Err(err).Msg("init_consumer_err")
		os.Exit(6)
	}
	//InitProducetIndex
	//if err := producet.InitProducetIndex(config.GetConfig().Redis.Consumer.Inverted.Topic); err != nil {
	//	log.Error().Err(err).Msg("init_consumer_err")
	//	os.Exit(6)
	//}
	//数据全量同步
	if err := db.Reload(); err != nil {
		log.Error().Err(err).Msg("init_db_failed")
		os.Exit(6)
	}
	//http 客户端初始化
	client.Init(2)

	//gin框架初始化
	gin_proxy.Init(3)

	//配置路由器
	router.Router()

}

// 退出函数
func exitFunc() {
}
