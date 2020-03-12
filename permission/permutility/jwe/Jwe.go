/*
 * @Author: your name
 * @Date: 2020-02-27 14:22:57
 * @LastEditTime: 2020-03-12 19:06:43
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\utility\Jwe.go
 */
package jwe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"gitee.com/chunanyong/logger"

	"github.com/joho/godotenv"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var jweENCRYPTKEY string
var jweSIGNEDKEY string
var jweEXPIRETIME string
var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey
var signedKey []byte
var sig jose.Signer
var enc jose.Encrypter
var isJweInited = false

func init() {
	godotenv.Load()
	jweENCRYPTKEY = os.Getenv("JWEENCRYPTKEY")
	jweSIGNEDKEY = os.Getenv("JWESIGNEDKEY")
	jweEXPIRETIME = os.Getenv("JWEEXPIRETIME")

	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err.(error), logger.String("JWE", "初始化JWE失败"))
		}
	}()

	signedKey = []byte(jweSIGNEDKEY)

	//加载加密私钥文件
	privateKeyPEMByte, err := ioutil.ReadFile(jweENCRYPTKEY)
	if err != nil {
		panic(err)
	}
	keyPEMBlock, rest := pem.Decode(privateKeyPEMByte)
	if len(rest) > 0 {
		panic("Decode key failed!")
	}
	//获取私钥
	privateKey, err = x509.ParsePKCS1PrivateKey(keyPEMBlock.Bytes)
	if err != nil {
		panic(err)
	}
	//获取公钥
	publicKey = &privateKey.PublicKey
	//创建jose.Signer
	sig, err = jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: signedKey},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	enc, err = jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"))
	if err != nil {
		panic(err)
	}

	isJweInited = true
}

func errorRecover(err *error) {
	if e := recover(); e != nil {
		if er, ok := e.(error); ok {
			*err = er
		} else {
			*err = errors.New(e.(string))
		}
		logger.Error(*err)
	}
}

// 创建Token字符串 id：账户唯一ID extInfo:需要在token中保存的扩展信息
func CreateToken(id string, extInfo interface{}) (raw string, err error) {

	raw = ""
	defer errorRecover(&err)

	if !isJweInited {
		panic("JWE未初始化成功")
	}

	//设置过期时间
	dutation, _ := time.ParseDuration(jweEXPIRETIME)
	cl := jwt.Claims{
		Subject: "readygo",
		Issuer:  "readygo",
		ID:      id,
		Expiry:  jwt.NewNumericDate(time.Now().Add(dutation)),
	}
	nestedBuilder := jwt.SignedAndEncrypted(sig, enc).Claims(cl)
	if extInfo != nil {
		nestedBuilder = nestedBuilder.Claims(extInfo)
	}
	raw, err = nestedBuilder.CompactSerialize()
	if err != nil {
		panic(err)
	}
	return raw, nil

}

//根据token获取用户id 和 扩展信息  扩展信息extInfo传结构体的指针
func GetInfoFromToken(token string, extInfo interface{}) (id string, err error) {

	defer errorRecover(&err)

	tok, err := jwt.ParseSignedAndEncrypted(token)
	if err != nil {
		panic(err)
	}

	//解密
	nested, err := tok.Decrypt(privateKey)
	if err != nil {
		panic(err)
	}

	//验签并返回Claim对象和扩展对象
	claim := jwt.Claims{}
	if err := nested.Claims(signedKey, &claim, extInfo); err != nil {
		panic(err)
	}

	//验证token有效性 这里暂时只验证了有效期
	if err := claim.Validate(
		jwt.Expected{
			Time: time.Now(),
		}); err != nil {
		panic(err)
	}

	return claim.ID, nil
}
