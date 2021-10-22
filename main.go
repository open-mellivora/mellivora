// Copyright 2021 Baidu Inc. All rights reserved.
// Use of this source code is governed by a xxx
// license that can be found in the LICENSE file.

/*
modification history
--------------------
2021/10/18 19:23:42, by wangyufeng04@baidu.com, create
*/

// Package main is special.  It defines a
// standalone executable program, not a library.
// Within package main the function main is also
// special—it’s where execution of the program begins.
// Whatever main does is what the program does.
package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/pkg/errors"
	"icode.baidu.com/baidu/go-lib/log"
	"icode.baidu.com/baidu/go-lib/log/log4go"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/config"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares/coding"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares/downlimiter"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares/dupefilter"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares/logging"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares/recovery"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares/statscollector"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/middleware/saver"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/spider"
)

const serviceName = "mini_spider"

var (
	_          = flag.String("v", "", "版本")
	configPath = flag.String("c", "conf", "配置文件路径")
	logPath    = flag.String("l", "log", "日志文件路径")
	_          = flag.String("h", "", "帮助")
)

// main the function where execution of the program begins
func main() {
	flag.Parse()
	var cfg *config.Config
	var err error
	if cfg, err = config.ParseConfig(*configPath); err != nil {
		err = errors.Wrap(err, "打开配置文件失败")
		panic(err)
	}
	var (
		logger          log4go.Logger
		f               *os.File
		saverMiddleware *saver.Middleware
	)
	if logger, err = log.Create(serviceName,
		"DEBUG", *logPath, true,
		"MIDNIGHT", 0); err != nil {
		err = errors.Wrap(err, "日志组建初始化失败")
		panic(err)
	}

	defer logger.Close()
	f, err = os.Open(cfg.URLListFile)
	if err != nil {
		err = errors.Wrap(err, "打开种子文件失败")
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var urls []string
	if err = decoder.Decode(&urls); err != nil {
		err = errors.Wrap(err, "配置文件解析失败")
		panic(err)
	}

	limiterMiddleware := downlimiter.NewMiddleware(&downlimiter.Config{
		ConcurrentRequests:          cfg.ThreadCount,
		ConcurrentRequestsPerDomain: cfg.ThreadCount,
		DownloadDelayPerDomain:      time.Duration(cfg.CrawlInterval * float64(time.Second)),
		Timeout:                     time.Duration(cfg.CrawlTimeout * float64(time.Second)),
		MaxDepth:                    cfg.MaxDepth,
	})

	saverMiddleware, err = saver.NewMiddleware(&saver.Config{
		Dir:     cfg.OutputDirectory,
		Pattern: cfg.TargetURL,
	})
	if err != nil {
		err = errors.Wrap(err, "初始化存储程序失败")
		panic(err)
	}

	engine := core.NewEngine()
	engine.SetLogger(logger)
	engine.Use(
		dupefilter.NewMiddleware(nil),  // 去重
		limiterMiddleware,              // 请求限制
		statscollector.NewMiddleware(), // 状态收集
		recovery.NewMiddleware(nil),    // panic捕获
		logging.NewMiddleware(),        // 日志打印
		saverMiddleware,                // 网页存储
		coding.NewDecoder(),            // 统一转为utf-8编码
	)

	s := spider.NewSimpleSpider(urls)
	engine.Run(s)
}
