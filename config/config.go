package config

import (
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/gcfg.v1"
)

const configFile = "spider.conf"

// SpiderConfig 爬虫配置
type SpiderConfig struct {
	// URLListFile 种子文件路径
	URLListFile string `gcfg:"urlListFile" validate:"required"`
	// OutputDirectory 抓取结果存储目录
	OutputDirectory string `gcfg:"outputDirectory" validate:"required"`
	// MaxDepth 最大抓取深度(种子为0级)
	MaxDepth int64 `gcfg:"maxDepth" validate:"gte=0,lte=128"`
	// CrawlInterval 抓取间隔. 单位: 秒
	CrawlInterval float64 `gcfg:"crawlInterval" validate:"gte=0"`
	// CrawlTimeout 抓取超时. 单位: 秒
	CrawlTimeout float64 `gcfg:"crawlTimeout" validate:"gte=0"`
	// TargetURL 需要存储的目标网页URL pattern(正则表达式)
	TargetURL string `gcfg:"targetUrl"`
	// ThreadCount 抓取routine数
	ThreadCount int64 `gcfg:"threadCount" validate:"gte=0"`
}

type Config struct {
	SpiderConfig `gcfg:"spider"`
}

func ParseConfig(filePath string) (cfg *Config, err error) {
	path := filepath.Join(filePath, configFile)
	cfg = new(Config)
	if err = gcfg.ReadFileInto(cfg, path); err != nil {
		return
	}
	v := validator.New()
	err = v.Struct(cfg)
	return
}
