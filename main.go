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
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"icode.baidu.com/baidu/go-lib/log"
	"icode.baidu.com/baidu/go-lib/log/log4go"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/config"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core/middlewares"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/middleware/saver"
	"icode.baidu.com/baidu/goodcoder/wangyufeng04/spider"
)

const (
	Version     = "v1.0"
	ServiceName = "mini_spider"
)

var (
	version    = flag.Bool("v", false, "版本")
	configPath = flag.String("c", "conf", "配置文件路径")
	logPath    = flag.String("l", "log", "日志文件路径")
	help       = flag.Bool("h", false, "帮助")
)

// main the function where execution of the program begins
func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	if *version {
		fmt.Println(Version)
		return
	}

	var cfg *config.Config
	var err error
	if cfg, err = config.ParseConfig(*configPath); err != nil {
		panic(errors.Wrap(err, "打开配置文件失败"))
	}

	var (
		logger          log4go.Logger
		f               *os.File
		saverMiddleware *saver.Middleware
	)

	if logger, err = log.Create(
		ServiceName, "DEBUG", *logPath, true, "MIDNIGHT", 0); err != nil {
		panic(errors.Wrap(err, "日志组建初始化失败"))
	}

	defer logger.Close()

	f, err = os.Open(cfg.URLListFile)
	if err != nil {
		panic(errors.Wrap(err, "打开种子文件失败"))
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var urls []string
	if err = decoder.Decode(&urls); err != nil {
		panic(errors.Wrap(err, "配置文件解析失败"))
	}

	limiterMiddleware := middlewares.NewDownLimiterWithConfig(middlewares.DownLimiterConfig{
		ConcurrentRequestsPerDomain: cfg.ThreadCount,
		DownloadDelayPerDomain:      time.Duration(cfg.CrawlInterval * float64(time.Second)),
		Timeout:                     time.Duration(cfg.CrawlTimeout * float64(time.Second)),
		MaxDepth:                    cfg.MaxDepth,
	})

	if saverMiddleware, err = saver.NewMiddleware(&saver.Config{
		Dir:     cfg.OutputDirectory,
		Pattern: cfg.TargetURL,
	}); err != nil {
		panic(errors.Wrap(err, "初始化存储程序失败"))
	}

	engine := core.NewEngine(cfg.ThreadCount)
	engine.SetLogger(logger)
	engine.Use(
		middlewares.NewDupeFilter(),
		limiterMiddleware,
		middlewares.NewStatsCollector(),
		middlewares.NewRecover(),
		middlewares.NewLogging(),
		saverMiddleware,
		middlewares.NewDecoder(),
	)

	s := spider.NewSimpleSpider(urls)
	engine.Run(s)
}
