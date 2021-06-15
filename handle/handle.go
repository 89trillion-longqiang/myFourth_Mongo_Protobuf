package handle

import (
	"fmt"
	"giftCode/Protobuf"
	"giftCode/gift"
	"giftCode/mongo"
	"giftCode/userInfo"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
	"time"

	"giftCode/setUpRedis"
)
func HandleAdminCreatGiftcode(des string,GiftNum string,ValidPeriod string,GiftContent string,CreatePer string)  map[string]string{
	retMap := make(map[string]string,2)
	GiftCode := gift.GetRandCode(8)
	CreatTime := time.Unix(time.Now().Unix(),0).Format("2006-01-02 15:04:05")
	err := setUpRedis.HashSet(gift.Gift{
		GiftCode,
		des,
		GiftNum,
		ValidPeriod,
		GiftContent,
		CreatePer,
		CreatTime,
		"0",
		"",
	})
	if err != nil {
		retMap["condition"]="error"
		retMap["GiftCode" ] =  err.Error()
	}else {
		retMap["condition"]="success"
		retMap["GiftCode" ] =  GiftCode
	}

	return retMap
}
func HadnleAdminInquireGiftCode(GiftCode string) (map[string]string,map[string]string){
	retMap := make(map[string]string,2)
	if len(GiftCode) == 0 {
		retMap["condition"]="error"
		retMap["giftCode"] = "GiftCode is empty"
		return retMap , nil
	}

	if setUpRedis.ExistsKey(GiftCode) {
		ret, err := setUpRedis.HashGetAll(GiftCode)
		if err != nil {
			retMap["condition"]="error"
			retMap["giftCode"] =  err.Error()
		} else if len(ret) != 0 {
			retMap["condition"]="success"

			return retMap,ret
		}
	}else {
		retMap["condition"]="error"
		retMap["giftCode"] = "GiftCode is error"
	}
	return retMap , nil
}
func HandleClient(GiftCode string,userName string) map[string]string{
	var errString string
	var flagCondition bool
	retMap := make(map[string]string,3)
	retMap["GiftCode"] = GiftCode
	if GiftCode == ""{
		retMap["condition"]="error"
		retMap["GiftCode" ]="input GiftCode"
		return retMap
	}
	if userName == ""{
		retMap["condition"]="error"
		retMap["GiftCode" ]="input usr"
		return retMap
	}

	ret , err := setUpRedis.HashGetAll(GiftCode)
	ret["GiftCode"] = GiftCode
	if err != nil{
		retMap["condition"]="error"
		retMap["GiftCode" ]= err.Error()
		return retMap
	}

	creatTime,_:=time.Parse("2006-01-02",ret["CreatTime"])
	curTime ,_:=time.Parse("2006-01-02",time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	d:=creatTime.Sub(curTime)
	validFlo, _ := strconv.ParseFloat(ret["ValidPeriod"], 64)

	if validFlo > d.Hours()/24{
		flagCondition =true
	}else {
		errString = "Expired"
	}
	Claim:= ret["ClaimList"]
	ClaimList := strings.Fields(Claim)
	for _, s := range ClaimList {
		if userName == s {

			retMap["condition"]="error"
			retMap["GiftCode"] = "User has received"
			return retMap
		}
	}

	avaNum , _:=strconv.Atoi(ret["AvailableNum"])
	giftNum,_ :=strconv.Atoi(ret["GiftNum"])
	if avaNum+1 <= giftNum {
		flagCondition = true
		outString := "{ " + userName + " "+ time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")+"}"
		avanumS := strconv.Itoa(avaNum+1)
		ret["AvailableNum"] = avanumS
		ret["ClaimList"] = string(ret["ClaimList"]) +" " + outString + ";"

		err := setUpRedis.HashSetMap(ret)
		if err != nil {
			fmt.Printf("%s",err)
		}
	}else {
		errString = ""
		errString = "Insufficient quantity"
		flagCondition = false
	}
	if flagCondition {
		retMap["condition"]="success"
		retMap["GiftContent"] = ret["GiftContent"]
	}else {

		retMap["condition"]   = "error"
		retMap["GiftContent"] = errString
	}
	return retMap
}

func HandleLogin(id string)  (string ,map[string]string) {
	retMap := make(map[string]string,3)
	var (
		data string
		flag bool
	)
	ret , err :=mongo.FindMongo("uid",id,"info")
	if ret.Uid == ""{
		fmt.Println("UID  is  empty")
		data = "no_exist"


		newUid :=userInfo.GetRandCode(8)
		for flag=true ; flag ;{
			if mongo.ExistId(newUid,"info") {
				newUid =userInfo.GetRandCode(8)
			}else {
				flag = false
			}
		}

		if mongo.InsertMongo(userInfo.UserInfo{Uid: newUid, Gold: "0", Diamond: "0"},"info") {
			retMap["Uid"] = newUid
			retMap["Gold"] = "0"
			retMap["Diamond"] = "0"
		}else {
			data = "error"
		}
	}else if ret.Uid != "" {
		data = "success"
		retMap["Uid"] = ret.Uid
		retMap["Gold"] = ret.Gold
		retMap["Diamond"] = ret.Diamond
	}else if err != nil {
		data = err.Error()
	}

	return data,retMap
}

func HandleVerGiftCode(GiftCode string,Uid string) ([]byte ,string){
	if GiftCode == ""{
		return nil,"error"
	}
	if Uid == ""{
		return nil,"error"
	}
	ret , err := setUpRedis.HashGetAll(GiftCode)
	ret["GiftCode"] = GiftCode
	if err != nil{
		return nil,"error"
	}
	creatTime,_:=time.Parse("2006-01-02",ret["CreatTime"])
	curTime ,_:=time.Parse("2006-01-02",time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	d:=creatTime.Sub(curTime)
	validFlo, _ := strconv.ParseFloat(ret["ValidPeriod"], 64)

	if validFlo > d.Hours()/24{

	}else {
		return nil, "Expired"
	}

	Claim:= ret["ClaimList"]
	ClaimList := strings.Fields(Claim)
	for _, s := range ClaimList {
		if Uid == s {
			return nil,"User has received"
		}
	}

	avaNum , _:=strconv.Atoi(ret["AvailableNum"])
	giftNum,_ :=strconv.Atoi(ret["GiftNum"])
	if avaNum+1 <= giftNum {
		outString := "{ " + Uid + " "+ time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")+"}"
		avanumS := strconv.Itoa(avaNum+1)
		ret["AvailableNum"] = avanumS
		ret["ClaimList"] = string(ret["ClaimList"]) +" " + outString + ";"

		err := setUpRedis.HashSetMap(ret)
		if err != nil {
			fmt.Printf("%s",err)
			return nil, "error"
		}
	}else {
		return  nil,"Insufficient quantity"
	}
	change := make(map[uint32]uint64,2)
	balance := make(map[uint32]uint64,2)
	counter := make(map[uint32]uint64,2)
	GiftContent := ret["GiftContent"]
	tempS :=strings.Split(GiftContent,":")
	beforeInfo,err:= mongo.FindMongo("uid",Uid,"info")

	if len(tempS) == 4{
		GoldNum_B ,_:=strconv.Atoi(beforeInfo.Gold)
		DiaNUm_B ,_ := strconv.Atoi(beforeInfo.Diamond)
		balance[0001] = uint64(GoldNum_B)
		balance[0002] = uint64(DiaNUm_B)
		GoldInc,_ := strconv.Atoi(tempS[1])
		DiaInc,_ := strconv.Atoi(tempS[3])
		change[0001] = uint64(GoldInc)
		change[0002] = uint64(DiaInc)
		GloNew := GoldNum_B+GoldInc
		DiaNew := DiaNUm_B +DiaInc

		mongo.UpdataMongo("uid",Uid, "gold", strconv.Itoa(GloNew))
		mongo.UpdataMongo("uid",Uid, "diamond",strconv.Itoa(DiaNew))
		counter[0001] = uint64(GloNew)
		counter[0002] = uint64(DiaNew)
	}

	if err != nil{
		return nil,"error"
	}

	var proInfo = Protobuf.GeneralReward{

		Changes:change,
		Balance:balance,
		Counter:counter,
	}
	retByte, _ := proto.Marshal(&proInfo)
	return retByte,"pass"
}



