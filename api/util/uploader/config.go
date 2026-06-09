package uploader

import "github.com/spf13/viper"

const defaultMaxSizeMB int64 = 3

// GetMaxBytes 获取配置的最大上传字节数，带默认值兜底
func GetMaxBytes() int64 {
	maxMB := viper.GetInt64("uploader.max_size")
	if maxMB <= 0 {
		maxMB = defaultMaxSizeMB
	}
	return maxMB * 1024 * 1024
}