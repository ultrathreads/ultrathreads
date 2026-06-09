package uploader

import (
	"path/filepath"
	"time"

	"ultrathreads/util"
)

// generateImageKey 根据文件内容 MD5 和日期生成唯一存储路径
// 输出格式: images/2026/06/09/<md5>.jpg
func generateImageKey(data []byte) string {
	md5 := util.MD5Bytes(data)
	dateDir := util.TimeFormat(time.Now(), "2006/01/02/")
	return filepath.Join("images", dateDir, md5+".jpg")
}