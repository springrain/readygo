/*
 * @Author: your name
 * @Date: 2020-03-12 12:26:40
 * @LastEditTime: 2020-03-12 12:29:13
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\permission\permstruct\UserVOStruct.go
 */
package permstruct

type UserVOStruct struct {
	UserId string

	Account string

	Email string

	UserName string

	Password   string
	CaptchaKey string

	Imgcaptcha string
	UserType   int
	Active     int

	// 私有的部门权限,用于处理单独url的特殊权限,调用 IUserRoleOrgService.wrapOrgIdFinderByPrivateOrgRoleId(String userId) 获取权限的finder;
	PrivateOrgRoleId string
	Captcha          string
}
