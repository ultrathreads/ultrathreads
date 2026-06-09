package uploader

// uploaderInterface 统一上传接口（包内使用，不对外暴露）
type uploaderInterface interface {
	PutImage(data []byte) (string, error)
	PutObject(key string, data []byte) (string, error)
	CopyImage(originUrl string) (string, error)
}

// PutImage 上传图片，自动生成存储 key
func PutImage(data []byte) (string, error) {
	return getUploader().PutImage(data)
}

// PutObject 按指定 key 上传二进制数据
func PutObject(key string, data []byte) (string, error) {
	return getUploader().PutObject(key, data)
}

// CopyImage 从远程 URL 下载并转存图片
func CopyImage(originUrl string) (string, error) {
	return getUploader().CopyImage(originUrl)
}