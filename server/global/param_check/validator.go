package param_check

import (
	"dcss/global"
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ch_translations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"
)

// InitValidate 初始化Validator数据校验
func InitValidate() {
	chinese := zh.New()
	uni := ut.New(chinese, chinese)
	global.Trans, _ = uni.GetTranslator("zh")
	global.Validate = validator.New()
	_ = ch_translations.RegisterDefaultTranslations(global.Validate, global.Trans)
	_ = global.Validate.RegisterValidation("checkMobile", checkMobile)

	global.LOG.Infof("初始化validator.v10数据校验器完成")
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `^(0|86|17951)?(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}

func VerifyParam(reqInfo interface{}) error {
	err := global.Validate.Struct(reqInfo)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s", err.Translate(global.Trans))
		}
	}
	return nil
}
