package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"rag_biddingchunks/internal/domain"
)

// 全局配置变量
var Cfg domain.RagConf

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}
