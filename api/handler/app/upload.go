package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"

	"ultrathreads/handler/base"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/uploader"
)

type UploadHandler struct {
	base.BaseHandler
}

// Upload 单文件上传
func (h *UploadHandler) Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	defer file.Close()

	fileBytes, err := readAndValidateFile(file, header.Size)
	if err != nil {
		h.Fail(ctx, util.NewErrorMsg(err.Error()))
		return
	}

	log.Info("上传文件：%s, size: %d", header.Filename, len(fileBytes))

	url, err := uploader.PutImage(fileBytes)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, gin.H{"url": url})
}

// UploadFromEditor 编辑器多文件上传
func (h *UploadHandler) UploadFromEditor(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}

	mForm, err := ctx.MultipartForm()
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	var (
		errFiles []string
		succMap  = make(map[string]string)
	)

	files := mForm.File["file[]"]
	for _, fileHeader := range files {
		f, err := fileHeader.Open()
		if err != nil {
			log.Error("打开文件失败 %s: %v", fileHeader.Filename, err)
			errFiles = append(errFiles, fileHeader.Filename)
			continue
		}

		// ✅ 修复：循环内必须立即 Close，不能用 defer（defer 在函数返回时才执行）
		fileBytes, err := readAndValidateFile(f, fileHeader.Size)
		f.Close() // 显式关闭，确保每次迭代都释放句柄

		if err != nil {
			log.Error("读取文件失败 %s: %v", fileHeader.Filename, err)
			errFiles = append(errFiles, fileHeader.Filename)
			continue
		}

		url, err := uploader.PutImage(fileBytes)
		if err != nil {
			log.Error("上传文件失败 %s: %v", fileHeader.Filename, err)
			errFiles = append(errFiles, fileHeader.Filename)
			continue
		}

		succMap[fileHeader.Filename] = url
	}

	h.Success(ctx, gin.H{
		"errFiles": errFiles,
		"succMap":  succMap,
	})
}

// UploadFromURL 通过 URL 转存图片
func (h *UploadHandler) UploadFromURL(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}

	data := make(map[string]string)
	if err := ctx.ShouldBindJSON(&data); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	rawURL := strings.TrimSpace(data["url"])
	if rawURL == "" {
		h.Fail(ctx, util.NewErrorMsg("url不能为空"))
		return
	}

	// SSRF 防护、大小限制等已在 uploader.CopyImage -> safeDownload 中统一处理
	output, err := uploader.CopyImage(rawURL)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, gin.H{
		"originalURL": rawURL,
		"url":         output,
	})
}

// readAndValidateFile 统一文件读取与大小校验
// 双重防护：header.Size 预检 + LimitReader 实际读取兜底，防止 Content-Length 伪造导致 OOM
func readAndValidateFile(reader io.Reader, declaredSize int64) ([]byte, error) {
	maxBytes := uploader.GetMaxBytes()

	if declaredSize > maxBytes {
		return nil, fmt.Errorf("文件大小超过限制（最大 %d MB）", maxBytes/1024/1024)
	}

	limitedReader := io.LimitReader(reader, maxBytes+1)
	fileBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	if int64(len(fileBytes)) > maxBytes {
		return nil, fmt.Errorf("文件实际大小超过限制（最大 %d MB）", maxBytes/1024/1024)
	}

	return fileBytes, nil
}
