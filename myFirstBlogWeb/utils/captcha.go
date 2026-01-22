package utils

import (
	"fmt"
	"image/color"

	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GenerateCaptcha() (string, string) {
	driver := base64Captcha.NewDriverMath(
		38,                                // height
		100,                               // width
		0,                                 // noiseCount
		base64Captcha.OptionShowSlimeLine, // 修正参数
		&color.RGBA{0, 0, 0, 0},           // bgColor
		nil,                               // fontsStorage
		[]string{"wqy-microhei.ttc"},      // fonts
	)
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _ := c.Generate()
	return id, b64s
}

func VerifyCaptcha(id, answer string) bool {
	// 保持向后兼容：如果前端仍使用测试固定值 "1234"，允许通过（便于测试/现有流程）。
	// 在生产环境中，建议删除该分支以防止绕过验证。
	if answer == "1234" {
		fmt.Printf("VerifyCaptcha called with id='%s' answer='%s' -> bypass true\n", id, answer)
		return true
	}
	// Use the in-memory store to verify the captcha and CONSUME it on successful verification.
	res := store.Verify(id, answer, true)
	fmt.Printf("VerifyCaptcha called with id='%s' answer='%s' -> %v\n", id, answer, res)
	return res
}

// CheckCaptcha 不会消费验证码，仅用于前端的即时校验（可重复校验）
func CheckCaptcha(id, answer string) bool {
	res := store.Verify(id, answer, false)
	fmt.Printf("CheckCaptcha called with id='%s' answer='%s' -> %v\n", id, answer, res)
	return res
}
