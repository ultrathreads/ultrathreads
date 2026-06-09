package uploader

import (
	"fmt"
	"io"
	"net"
	"net/url"

	"ultrathreads/util/log"
)

// privateRanges 私有/保留 IP 段列表
var privateRanges = func() []*net.IPNet {
	cidrs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"0.0.0.0/8",
		"100.64.0.0/10", // CGNAT
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil {
			nets = append(nets, network)
		}
	}
	return nets
}()

// safeDownload 安全下载远程资源，带 SSRF 防护和大小限制
func safeDownload(rawURL string) ([]byte, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	// 仅允许 http/https 协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported url scheme: %s", parsedURL.Scheme)
	}

	// DNS 解析并校验目标 IP，阻断内网访问
	host := parsedURL.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("dns lookup failed for %s: %w", host, err)
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return nil, fmt.Errorf("access to private/reserved IP %s is forbidden (SSRF protection)", ip.String())
		}
	}

	maxSize := GetMaxBytes()

	// 流式读取，不自动将响应体载入内存
	resp, err := httpClient.R().
		SetDoNotParseResponse(true).
		Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	// ✅ 关键修复：SetDoNotParseResponse(true) 要求调用方必须在所有路径关闭 RawBody
	defer resp.RawBody().Close()

	// ✅ 修复：通过底层 *http.Response 获取 ContentLength
	contentLength := resp.RawResponse.ContentLength
	if contentLength > 0 && contentLength > maxSize {
		return nil, fmt.Errorf("remote file too large: %d bytes, max allowed: %d", contentLength, maxSize)
	}

	// LimitReader 兜底，防止实际传输超过限制或 Content-Length 缺失/伪造
	limitedReader := &io.LimitedReader{R: resp.RawBody(), N: maxSize + 1}
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	if int64(len(data)) > maxSize {
		return nil, fmt.Errorf("remote file exceeds max size limit (%d bytes)", maxSize)
	}

	log.Info("Safe download success: %s, size: %d", rawURL, len(data))
	return data, nil
}

// isPrivateIP 判断 IP 是否属于私有/保留地址段
func isPrivateIP(ip net.IP) bool {
	for _, network := range privateRanges {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}