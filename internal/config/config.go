package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// 默认配置文件路径
const DefaultConfigFile = ".acr/config.yaml"

// Config 结构体，保存所有配置信息
type Config struct {
	Token  string
	Prompt string
	Model  string
	Url    string
}

// InitConfigFile 初始化配置文件（若已存在则返回提示，若不存在则创建并写入默认内容）
func InitConfigFile(configFile string) error {
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	configFile = getConfigPath(configFile)

	dir := filepath.Dir(configFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %v", err)
		}
	}
	if _, err := os.Stat(configFile); err == nil {
		return fmt.Errorf("配置文件已存在: %s", configFile)
	}
	v := viper.New()
	v.Set("token", "")
	v.Set("prompt", "请帮我审查以下代码变更，指出潜在问题并给出建议")
	v.Set("model", "")
	v.Set("url", "")
	if err := v.SafeWriteConfigAs(configFile); err != nil {
		return fmt.Errorf("初始化配置文件失败: %v", err)
	}
	return nil
}

// UpdateConfigFile 批量更新配置项，若文件不存在则新建
func UpdateConfigFile(configFile string, updates Config) error {
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	configFile = getConfigPath(configFile)
	dir := filepath.Dir(configFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %v", err)
		}
	}
	v := viper.New()
	v.SetConfigFile(configFile)
	_ = v.ReadInConfig() // 不存在也不报错
	if updates.Token != "" {
		v.Set("token", updates.Token)
	}
	if updates.Prompt != "" {
		v.Set("prompt", updates.Prompt)
	}
	if updates.Model != "" {
		v.Set("model", updates.Model)
	}
	if updates.Url != "" {
		v.Set("url", updates.Url)
	}

	if err := v.WriteConfigAs(configFile); err != nil {
		// 文件不存在则创建
		if os.IsNotExist(err) {
			if err := v.SafeWriteConfigAs(configFile); err != nil {
				return fmt.Errorf("写入配置失败: %v", err)
			}
		} else {
			return fmt.Errorf("写入配置失败: %v", err)
		}
	}
	return nil
}

// getConfigPath 处理 ~ 路径，返回绝对路径
func getConfigPath(filePath string) string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return ""
	}
	return filepath.Join(homeDir, filePath)
}

// LoadConfig 加载配置，优先级：命令行参数 > 环境变量 > 配置文件
func LoadConfig(configFile string) (*Config, error) {
	v := viper.New()

	// 配置文件路径
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	configFile = getConfigPath(configFile)
	v.SetConfigFile(configFile)

	// 环境变量前缀
	v.SetEnvPrefix("AI_CODE_REVIEWER")
	v.AutomaticEnv()

	// 默认值
	v.SetDefault("model", "gpt-3.5-turbo")
	v.SetDefault("prompt", "请帮我审查以下代码变更，指出潜在问题并给出建议")

	// 读取配置文件（可选）
	if _, err := os.Stat(configFile); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	cfg := &Config{
		Token:  v.GetString("token"),
		Prompt: v.GetString("prompt"),
		Model:  v.GetString("model"),
		Url:    v.GetString("url"),
	}

	// if cfg.Token == "" {
	// 	return nil, fmt.Errorf("API token 未配置，请在配置文件、环境变量或命令行参数中设置 token")
	// }

	return cfg, nil
}
