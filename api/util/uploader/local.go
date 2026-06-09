package uploader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"ultrathreads/util/urls"
)

type localUploader struct{}

func NewLocal() *localUploader {
	return &localUploader{}
}

func (l *localUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return l.PutObject(key, data)
}

func (l *localUploader) PutObject(key string, data []byte) (string, error) {
	basePath := viper.GetString("uploader.local.path")
	if basePath == "" {
		return "", fmt.Errorf("uploader.local.path is not configured")
	}

	// 安全校验：清洗路径并阻止目录遍历
	cleanKey := filepath.Clean(key)
	if strings.Contains(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return "", fmt.Errorf("invalid file key detected: %s", key)
	}

	filename := filepath.Join(basePath, cleanKey)

	// 二次校验：确保最终绝对路径仍在 basePath 下
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("resolve base path failed: %w", err)
	}
	absFile, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("resolve file path failed: %w", err)
	}
	if !strings.HasPrefix(absFile, absBase+string(os.PathSeparator)) && absFile != absBase {
		return "", fmt.Errorf("path traversal attack blocked: %s", key)
	}

	// 创建目录，权限 0755
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return "", fmt.Errorf("create directory failed: %w", err)
	}

	// 写入文件，权限 0644（非 ModePerm 0777）
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("write file failed: %w", err)
	}

	return urls.UrlJoin(viper.GetString("uploader.local.host"), cleanKey), nil
}

func (l *localUploader) CopyImage(originUrl string) (string, error) {
	data, err := safeDownload(originUrl)
	if err != nil {
		return "", fmt.Errorf("copy image download failed: %w", err)
	}
	return l.PutImage(data)
}