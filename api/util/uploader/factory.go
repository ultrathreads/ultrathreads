package uploader

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"ultrathreads/util/log"
)

var (
	activeUploader uploaderInterface
	initOnce       sync.Once

	// httpClient 全局复用连接池，禁止自动重定向防止 SSRF 绕过
	httpClient = resty.New().
			SetTimeout(30 * time.Second).
			SetRedirectPolicy(resty.NoRedirectPolicy()).
			SetTransport(&http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			})
)

// getUploader 懒加载并缓存当前激活的上传器实例
func getUploader() uploaderInterface {
	initOnce.Do(func() {
		switch viper.GetString("uploader.enable") {
		case "aliyun", "oss":
			activeUploader = NewAliyun()
			log.Info("Uploader initialized: aliyun oss")
		default:
			activeUploader = NewLocal()
			log.Info("Uploader initialized: local filesystem")
		}
	})
	return activeUploader
}