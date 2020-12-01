package common

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
)
type JwtTokenInfo struct {
	AppId    int64 `json:"app_id"`    // appId
	UserId   int64 `json:"user_id"`   // 用户id
	//UserName   int64 `json:"user_name"`   // 用户名称
	Passwd   string `json:"passwd"`	//用户密码
	//Expire   int64 `json:"expire"`    // 过期时间
}



func JwtEncry(info JwtTokenInfo) string {

	dataByte,_:= json.Marshal(info)
	var dataStr = string(dataByte)
	data := jwt.StandardClaims{Subject:dataStr,ExpiresAt:time.Now().Unix()-1000}
	tokenInfo := jwt.NewWithClaims(jwt.SigningMethodHS256,data)
	tokenStr,_ := tokenInfo.SignedString([]byte(KeyInfo))
	//Sugar.Info("generate token: ",tokenStr)
	return tokenStr

}
func JwtDecry(tokenStr string)(JwtTokenInfo,bool){
	tokenInfo , _ := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		return KeyInfo,nil
	})

	err := tokenInfo.Claims.Valid()
	if err!=nil{
		println(err.Error())
	}

	fineToken := tokenInfo.Claims.(jwt.MapClaims)
	succ := fineToken.VerifyExpiresAt(time.Now().Unix(),true)
	return fineToken["sub"].(JwtTokenInfo),succ
}