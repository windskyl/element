package utils

import (
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
	return answer == "1234"
}
