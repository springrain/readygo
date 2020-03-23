/*
 * @Author: your name
 * @Date: 2020-02-26 20:15:09
 * @LastEditTime: 2020-03-12 19:24:31
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permhandler\PermHandler.go
 */
package permhandler

import (
	"fmt"
	"net/http"
	"readygo/apistruct"
	"readygo/permission/permservice"
	"readygo/permission/permstruct"
	"readygo/permission/permutil"
	"strings"

	"github.com/gin-gonic/gin"

	"gitee.com/chunanyong/logger"
)

//权限过滤器
func PermHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		//处理跨域
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		method := c.Request.Method
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		//装逼一点,禁止所有的GET方法
		// if method == "GET" {
		// 	c.AbortWithStatus(http.StatusMethodNotAllowed)
		// }

		responseBody := apistruct.ResponseBodyModel{}
		ctx := c.Request.Context()
		//请求的uri
		uri := GetPatternURI(c)
		logger.Info(uri)

		if IsExcludePath(uri) {
			c.Next()
			return
		}

		user := permstruct.UserVOStruct{}
		token := c.GetHeader("READYGOTOKEN")
		if token == "" {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = "缺少Token"
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		// 获取Token并检测有效期
		userID, err := permutil.JWEGetInfoFromToken(token, &user)
		if err != nil {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = fmt.Sprintf("%s%s", "解析Token失败", err.Error())
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		//TODO 这里需要添加权限判断逻辑
		// 不知道 u_10001什么意思
		// if userID == "u_10001" {
		// 	c.Next()
		// 	return
		// }
		// 获取权限拥有的菜单信息
		permMenuList, err := permservice.FindMenuByUserId(ctx, userID)
		if err != nil {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = fmt.Sprintf("%s%s", "获取用户菜单失败", err.Error())
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		if len(permMenuList) == 0 {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = "没有任何权限"
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		roleID := ""
		for _, item := range permMenuList {
			if item.Pageurl == "" {
				continue
			}
			if strings.ToLower(item.Pageurl) == strings.ToLower(uri) {
				roleID = item.RoleId
				break
			}
		}

		if roleID == "" {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = "没有当前操作权限"
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		role, err := permservice.FindRoleStructById(ctx, roleID)
		if err != nil {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = fmt.Sprintf("%s%s", "FindRoleStructById失败", err.Error())
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		userVO, err := permservice.FindUserVOStructByUserId(ctx, userID)
		if err != nil {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = fmt.Sprintf("%s%s", "FindUserVOStructByUserId失败", err.Error())
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		userVO.PrivateOrgRoleId = role.Id

		// 设置当前登录用户到上下文
		ctx, _ = SetCurrentUser(c.Request.Context(), userVO)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetPatternURI 获取格式化路径
func GetPatternURI(c *gin.Context) string {
	return c.Request.RequestURI
}

/*
func jwe() {
	// Generate a public/private key pair to use for this example.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	publicKey := &privateKey.PublicKey
	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, nil)
	if err != nil {
		panic(err)
	}

	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	var plaintext = []byte("Lorem ipsum dolor sit amet")
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()

	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err = jose.ParseEncrypted(serialized)
	if err != nil {
		panic(err)
	}

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate the the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	decrypted, err := object.Decrypt(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(decrypted))
}
*/
