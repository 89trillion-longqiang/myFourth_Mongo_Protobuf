package gift

import (
	"bytes"
	"math/rand"
	"time"
)

type Gift struct {
	GiftCode string      ///礼包码
	Description string	 ///礼包描述
	GiftNum string		 ///可领取数
	ValidPeriod string	 ///有效期
	GiftContent string   ///礼包内容
	CreatePer string     ///创建用户
	CreatTime string     ///创建时间

	AvailableNum string  ///已领取次数
	ClaimList string     ///领取列表
}

func (g Gift)IsEmpty() bool{
	if len(g.GiftCode) == 0 {return true}
	if len(g.Description) == 0 {return true}
	if len(g.GiftNum) == 0 {return true}
	if len(g.ValidPeriod) == 0 {return true}
	if len(g.GiftContent) == 0 {return true}
	if len(g.CreatePer) == 0 {return true}
	if len(g.CreatTime) == 0 {return true}
	return false
}

func GetRandCode(codeLen int) string {
	// 1. 定义原始字符串
	rawStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取
	rand.Seed(time.Now().UnixNano())
	for rawStrLen := len(rawStr);codeLen > 0; codeLen-- {
		randNum := rand.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	return b.String()
}