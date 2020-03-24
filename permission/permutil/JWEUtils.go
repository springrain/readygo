package permutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

//rsa 私钥
var privateKey *rsa.PrivateKey

//rsa公钥
var publicKey *rsa.PublicKey

//jwe签名
var signer jose.Signer

//jwe加密
var encrypter jose.Encrypter

//expireDuration 过期时间
var expireDuration time.Duration

//var jwtSecretByte []byte

var _jwtSecret string

//NewJWEConfig 根据配置新建JWE对象
//jweRSAPrivatePemFilePath 证书路径
//jwtSecretToken jwt的加密token
//jweExpireSecond 超时的秒数 默认 24小时=60*60*24
func NewJWEConfig(jweRSAPrivatePemFilePath string, jwtSecret string, jweExpireSecond int) error {
	if jweExpireSecond == 0 {
		jweExpireSecond = 60 * 60 * 24
	}
	expireDuration = time.Second * time.Duration(jweExpireSecond)
	_jwtSecret = jwtSecret
	//jwtSecretByte = []byte(_jwtSecret)
	//加载加密私钥文件
	privateKeyPEMByte, err := ioutil.ReadFile(jweRSAPrivatePemFilePath)
	if err != nil {
		return err
	}
	keyPEMBlock, rest := pem.Decode(privateKeyPEMByte)
	if len(rest) > 0 {
		return errors.New("Decode key failed!jweRSAPrivatePemFilePath:" + jweRSAPrivatePemFilePath)
	}
	//获取私钥
	privateKey, err = x509.ParsePKCS1PrivateKey(keyPEMBlock.Bytes)
	if err != nil {
		return err
	}
	//获取公钥
	publicKey = &privateKey.PublicKey

	//创建jose.Signer
	//修改成每个用户独立的签名秘钥，这里不再统一生成
	// signer, err = jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: jwtSecretByte}, (&jose.SignerOptions{}).WithType("JWT"))
	// if err != nil {
	// 	return err
	// }
	//创建 jose.Encrypter
	encrypter, err = jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"))
	return err
}

//JWECreateToken 创建Token字符串 id：账户唯一ID extInfo:需要在token中保存的扩展信息
func JWECreateToken(id string, extInfo interface{}) (raw string, err error) {

	//设置过期时间
	cl := jwt.Claims{
		ID:     id,
		Expiry: jwt.NewNumericDate(time.Now().Add(expireDuration)),
	}
	jwtSecretByte := []byte(_jwtSecret + id)

	signerOption := jose.SignerOptions{}
	signerOption.WithType("JWT")
	signerOption.WithHeader("R", id)
	signer, err = jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: jwtSecretByte}, &signerOption) //(&jose.SignerOptions{}).WithType("JWT")
	if err != nil {
		return "", err
	}

	nestedBuilder := jwt.SignedAndEncrypted(signer, encrypter).Claims(cl)
	if extInfo != nil {
		nestedBuilder = nestedBuilder.Claims(extInfo)
	}
	raw, err = nestedBuilder.CompactSerialize()

	return raw, err

}

//根据token获取用户id 和 扩展信息  扩展信息extInfo传结构体的指针
func JWEGetInfoFromToken(token string, extInfo interface{}) (id string, err error) {
	tok, err := jwt.ParseSignedAndEncrypted(token)
	if err != nil {
		return "", err
	}

	//解密
	nested, err := tok.Decrypt(privateKey)
	if err != nil {
		return "", err
	}

	if len(nested.Headers) < 1 {
		return "", errors.New("缺少Header")
	}
	if _, ok := nested.Headers[0].ExtraHeaders["R"]; !ok {
		return "", errors.New("缺少Header")
	}
	id = nested.Headers[0].ExtraHeaders["R"].(string)

	//验签并返回Claim对象和扩展对象
	claim := jwt.Claims{}

	//先不验证签名获取token信息
	// if err := nested.UnsafeClaimsWithoutVerification(&claim, extInfo); err != nil {
	// 	return "", err
	// }
	//签名秘钥拼接用户id
	jwtSecretByte := []byte(_jwtSecret + id)
	//验证签名
	if err := nested.Claims(jwtSecretByte, &claim, extInfo); err != nil {
		return "", err
	}

	//验证token有效性 这里暂时只验证了有效期
	if err := claim.Validate(
		jwt.Expected{
			Time: time.Now(),
		}); err != nil {
		return "", err
	}

	return claim.ID, err
}
