package config

import (
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/gcfg.v1"
)

const configFile = "spider.conf"

// SpiderConfig spider config
type SpiderConfig struct {
	URLListFile     string  `gcfg:"urlListFile" validate:"required"`
	OutputDirectory string  `gcfg:"outputDirectory" validate:"required"`
	MaxDepth        int64   `gcfg:"maxDepth" validate:"gte=0,lte=128"`
	CrawlInterval   float64 `gcfg:"crawlInterval" validate:"gte=0"`
	CrawlTimeout    float64 `gcfg:"crawlTimeout" validate:"gte=0"`
	TargetURL       string  `gcfg:"targetUrl"`
	ThreadCount     int64   `gcfg:"threadCount" validate:"gte=0"`
}

type Config struct {
	SpiderConfig `gcfg:"spider"`
}

// ParseConfig 配置解析
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
