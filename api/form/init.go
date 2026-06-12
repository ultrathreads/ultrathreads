package form

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkgerrors "github.com/pkg/errors"
)

// ValidateErrorMessage : customize error messages
// ✅ 移到这里，确保 Bind 函数一定能访问到
var ValidateErrorMessage = map[string]string{
	"default":        "%s - %s is invalid(%s)",
	"customValidate": "%s can not be %s",
	"required":       "%s is required,got empty %#v",
	"pwdValidate":    "%s is not a valid password",
}

func init() {
	// Register custom validate methods
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//     _ = v.RegisterValidation("customValidate", customValidate)
	//     _ = v.RegisterValidation("pwdValidate", pwdValidate)
	// }
}

// Bind : bind request dto and auto verify parameters
func Bind(c *gin.Context, obj interface{}) error {
	// 1. 始终先从 URI 绑定（:slug 等路径参数）
	_ = c.ShouldBindUri(obj)

	// 2. 再尝试 JSON/Query/Form 绑定
	err := c.ShouldBind(obj)

	// 3. ✅ Body 为空(EOF)时，静默忽略 JSON 绑定错误
	if err != nil && errors.Is(err, io.EOF) {
		return nil
	}

	// 4. 其他绑定错误正常处理
	if err != nil {
		var valErrs validator.ValidationErrors
		if pkgerrors.As(err, &valErrs) {
			var tagErrorMsg []string
			for _, e := range valErrs {
				if msgTpl, has := ValidateErrorMessage[e.Tag()]; has {
					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(msgTpl, e.Field(), e.Value()))
				} else {
					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(ValidateErrorMessage["default"], e.Tag(), e.Field(), e.Value()))
				}
			}
			return pkgerrors.New(strings.Join(tagErrorMsg, ","))
		}
		return pkgerrors.Wrap(err, "参数绑定失败")
	}

	return nil
}