package main

import (
	"fmt"
	"testing"

	"giftCode/gift"
)

func Test_GetRandCode(t *testing.T){

	retString := gift.GetRandCode(8)
	if len(retString) != 8{
		fmt.Println("Test_GetRandCode error")
	}else {
		fmt.Println("Test_GetRandCode pass")
	}
}

