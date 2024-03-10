package middlewares

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	constLib "go-gateway/lib/const"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)


func TranslationMiddleware() gin.HandlerFunc {
	en := en.New()
	zh := zh.New()

	uni := ut.New(zh, zh, en)
	vtor := validator.New()

	//自定义验证方法
	CustomRegisterValidation(vtor)
	return func(ctx *gin.Context) {
		locale := ctx.DefaultQuery("locale", "zh")
		isZh := locale == "zh"
		trans, _ := uni.GetTranslator(locale)
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(vtor,trans)
			vtor.RegisterTagNameFunc(func(field reflect.StructField) string {
				return field.Tag.Get("en_comment") 
			})
		default:
			zh_translations.RegisterDefaultTranslations(vtor,trans)
			vtor.RegisterTagNameFunc(func(field reflect.StructField) string {
				return field.Tag.Get("comment") 
			})
		}

		//针对于自定义的 validator func 设置翻译器
		CustomRegisterTranslation(vtor, trans, isZh)
		// 设置全局 Context translator and validator
		ctx.Set(constLib.TranslatorKey, trans)
		ctx.Set(constLib.ValidatorKey, vtor)
		ctx.Next()
	}
}


func CustomRegisterTranslation(val *validator.Validate, trans ut.Translator, isZh bool) {
	val.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
		if isZh {
			return ut.Add("valid_username", "{0} 填写不正确", true)
		} else {
			return ut.Add("valid_username", "{0} is not correct", true)
		}
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_username", fe.Field())
		return t
	})

	val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
		return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_service_name", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_rule", trans, func(ut ut.Translator) error {
		return ut.Add("valid_rule", "{0} 必须是非空字符", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_rule", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
		return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_url_rewrite", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_header_transfor", trans, func(ut ut.Translator) error {
		return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_header_transfor", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_ipportlist", trans, func(ut ut.Translator) error {
		return ut.Add("valid_ipportlist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_ipportlist", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_iplist", trans, func(ut ut.Translator) error {
		return ut.Add("valid_iplist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_iplist", fe.Field())
		return t
	})
	val.RegisterTranslation("valid_weightlist", trans, func(ut ut.Translator) error {
		return ut.Add("valid_weightlist", "{0} 不符合输入格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_weightlist", fe.Field())
		return t
	})
}	

func CustomRegisterValidation(val *validator.Validate) {
	val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "admin"
	})

	val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
		return matched
	})

	val.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
		return matched
	})

	val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if len(strings.Split(ms, " ")) != 2 {
				return false
			}
		}
		return true
	})

	val.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if len(strings.Split(ms, " ")) != 3 {
				return false
			}
		}
		return true
	})

	val.RegisterValidation("valid_ipportlist", func(fl validator.FieldLevel) bool {
		reg, err := regexp.Compile(`^\S+\:\d+$`)

		if err != nil {
			return false
		}
		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if matched := reg.Match([]byte(ms)); !matched {
				return false
			}
		}
		return true
	})

	val.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		req, err := regexp.Compile(`\S+`)

		if err != nil {
			return false
		}

		for _, item := range strings.Split(fl.Field().String(), ",") {
			matched := req.Match([]byte(item)) //ip_addr
			if !matched {
				return false
			}
		}
		return true
	})

	val.RegisterValidation("valid_weightlist", func(fl validator.FieldLevel) bool {
		fmt.Println(fl.Field().String())
		req, err := regexp.Compile(`^\d+$`)

		if err != nil {
			return false
		}

		for _, ms := range strings.Split(fl.Field().String(), ",") {
			if matched := req.Match([]byte(ms)); !matched {
				return false
			}
		}
		return true
	})


}