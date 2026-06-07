package form

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func init() {
	// Register custom validate methods
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//     _ = v.RegisterValidation("customValidate", customValidate)
	//     _ = v.RegisterValidation("pwdValidate", pwdValidate)
	// }
}

// Bind : bind request dto and auto verify parameters
func Bind(c *gin.Context, obj interface{}) error {
	_ = c.ShouldBindUri(obj)
	if err := c.ShouldBind(obj); err != nil {
		// ✅ 使用 errors.As 安全断言，避免非 ValidationErrors 导致 panic
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			var tagErrorMsg []string
			for _, e := range valErrs {
				if msgTpl, has := ValidateErrorMessage[e.Tag()]; has {
					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(msgTpl, e.Field(), e.Value()))
				} else {
					tagErrorMsg = append(tagErrorMsg, fmt.Sprintf(ValidateErrorMessage["default"], e.Tag(), e.Field(), e.Value()))
				}
			}
			return errors.New(strings.Join(tagErrorMsg, ","))
		}

		// ✅ 非验证错误（如 JSON 类型不匹配、格式错误等），直接返回原始错误信息
		return errors.Wrap(err, "参数绑定失败")
	}

	return nil
}

// ValidateErrorMessage : customize error messages
var ValidateErrorMessage = map[string]string{
	"default":        "%s - %s is invalid(%s)",
	"customValidate": "%s can not be %s",
	"required":       "%s is required,got empty %#v",
	"pwdValidate":    "%s is not a valid password",
}