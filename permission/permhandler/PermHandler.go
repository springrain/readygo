/*
 * @Author: your name
 * @Date: 2020-02-26 20:15:09
 * @LastEditTime: 2021-03-05 17:23:21
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permhandler\PermHandler.go
 */
package permhandler

import (
	"context"
	"fmt"
	"net/http"
	"readygo/apistruct"
	"readygo/permission/permservice"
	"readygo/permission/permstruct"
	"readygo/permission/permutil"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

//JWTTokenName jwt的token名称
var JWTTokenName = "READYGOTOKEN"

//PermHandler 权限过滤器
func PermHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		//处理跨域
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		method := string(c.Request.Method())
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		//装逼一点,禁止所有的GET方法
		// if method == "GET" {
		// 	c.AbortWithStatus(http.StatusMethodNotAllowed)
		// }

		responseBody := apistruct.ResponseBodyModel{}
		//请求的uri
		uri := string(c.Request.URI().Path())
		hlog.Info(uri)

		//如果是不拦截的URL  TODO 此处因为权限拦截不支持正则 先放开swagger
		if isExcludePath(uri) || strings.Contains(uri, "swagger") {
			c.Next(ctx)
			return
		}

		// 权限的URL不能重复,通过访问的url,根据/api/user/menu返回的数据再遍历一次,获取menuId 和 roleId
		// 前后端分离之后,后台的菜单实际只管理了数据,并不管理前端的菜单层次结构.

		user := permstruct.UserVOStruct{}
		token := string(c.GetHeader(JWTTokenName))
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

		//如果用户默认有的权限
		if isUserDefaultPath(uri) {
			c.Next(ctx)
			return
		}

		// 角色关联部门实现数据权限,角色指定 roleOrgType <0自己的数据,1所在部门,2所在部门及子部门数据,3.自定义部门数据>.
		// 0就是类似员工,1和2 是根据使用者自身数据进行数据授权,适合公用全局属性,例如专员只能查看他所在的部门的数据,虽然他不是部门主管.

		// 角色指定 privateOrg 角色的部门是否私有,0否,1是,默认0.当角色私有时,菜单只使用此角色的部门权限,不再扩散到全局角色权限,用于设置特殊的菜单权限.
		// 公共权限时部门主管有所管理部门的数据全权限,无论角色是否分配.私有部门权限时,严格按照配置的数据执行,部门主管可能没有部门权限.
		// privateOrg私有权限和公共权限分别处理,不能交叉.处理公共权限时会跳过私有权限
		//  privateOrg 和 roleOrgType 交叉情况,比较复杂,场景也很少,暂时未细测,如果是私有的所在部门权限,应该只能查看所在部门的数据,也不会扩散权限.

		// 角色关联人员,部门,菜单,作为整个权限设计的中心枢纽.
		// 角色都有归属部门,其部门主管或上级主管才可以修改角色属性,其他使用人员只能往角色里添加人员,
		// 不能选择部门或则其他操作,只能添加人员,不然存在提权风险,例如 员工角色下有1000人, 如果给 角色 设置了部门,那这1000人都起效了.
		// 角色 shareRole 设置共享的角色可以被下级部门直接查看到,并添加人员.同样 也是只能添加人员.

		// 1.根据访问的url,通过permservice.FindMenuByUserId(ctx,userId)查询对应的menuId和roleId.需要确保url地址唯一,多个菜单url相同可以使用软跳转,暂时不处理.
		// 2.根据userId查询缓存的List<Role>,验证是否包含这个roleId
		// 3.根据roleId查询缓存的List<Menu>,验证是否包含这个menuId.
		// 4.查看roleId如果是私有权限,UserVo 就设置 privateOrgRoleId,业务调用SessionUser.getPrivateOrgRoleId获取私有的roleId,
		// 通过permservice.WrapOrgIdFinderByPrivateOrgRoleId(ctx, roleId, userId) 获取权限的 Finder
		// 5.如果是公共权限,这里不做处理,业务方法调用 通过permservice.WrapOrgIdFinderByUserId(ctx, userId) 获取权限的Finder

		// 注意:在返回前端菜单权限时,要包含menuId和roleId,私有privateOrg的roleId优先,如果同一个menuId存在多个定制roleId冲突,按照role的排序,同一个menuId只保留一个roleId.

		// 注意:缓存的清理,使用缓存,代码组装用户权限的树形结构.

		//TODO 这里需要添加权限判断逻辑
		// 不知道 u_10001什么意思
		// if userID == "u_10001" {
		// 	c.Next(ctx)
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
		//用户是否有uri的权限.循环遍历用户有权限的菜单URL,因为pageUrl是唯一的,可以取出menuId和roleId
		//	roleID := ""
		// for _, item := range permMenuList {
		// 	if item.Pageurl == "" {
		// 		continue
		// 	}
		// 	if strings.ToLower(item.Pageurl) == strings.ToLower(uri) {
		// 		roleID = item.RoleId
		// 		break
		// 	}
		// }

		// if roleID == "" {
		// 	responseBody.Status = http.StatusUnauthorized
		// 	responseBody.Message = "没有当前操作权限"
		// 	c.AbortWithStatusJSON(responseBody.Status, responseBody)
		// 	return
		// }

		// //根据roleId 查询 role
		// role, err := permservice.FindRoleStructById(ctx, roleID)
		// if err != nil {
		// 	responseBody.Status = http.StatusUnauthorized
		// 	responseBody.Message = fmt.Sprintf("%s%s", "FindRoleStructById失败", err.Error())
		// 	c.AbortWithStatusJSON(responseBody.Status, responseBody)
		// 	return
		// }

		userVO, err := permservice.FindUserVOStructByUserId(ctx, userID)
		if err != nil {
			responseBody.Status = http.StatusUnauthorized
			responseBody.Message = fmt.Sprintf("%s%s", "FindUserVOStructByUserId失败", err.Error())
			c.AbortWithStatusJSON(responseBody.Status, responseBody)
			return
		}

		// 如果是私有的部门权限,setPrivateOrgRoleId,业务调用SessionUser.getPrivateOrgRoleId,如果不是NULL,就调用IUserRoleOrgService.wrapOrgIdFinderByPrivateOrgRoleId(String roleId,String userId) 获取权限的 Finder
		// if role.PrivateOrg == 1 {
		// 	userVO.PrivateOrgRoleId = role.Id
		// }

		// 设置当前登录用户到上下文
		ctx, _ = permstruct.BindContextCurrentUser(ctx, userVO)
		//重新覆盖ctx
		c.Next(ctx)
	}
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
