package permapi

import (
	"readygo/ginext"
	"readygo/permission/permservice"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/zorm"
	"github.com/gin-gonic/gin"
)

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	Page       *zorm.Page
	UserName   string   `form:"userName"`   // 用户名
	QueryValue string   `form:"queryValue"` // 模糊查询
	Status     int      `form:"status"`     // 用户状态(1:启用 2:停用)
	RoleIDs    []uint64 `form:"-"`          // 角色ID列表
}

// Query 查询数据
func QueryMenu(c *gin.Context) {
	ctx := c.Request.Context()
	var params MenuQueryParam
	if err := ginext.ParseQuery(c, &params); err != nil {
		ginext.ResError(c, err)
		return
	}
	finder := zorm.NewFinder()
	finder.Append("SELECT r.* from ").Append(permstruct.MenuStructTableName).Append(" m,")

	result, err := permservice.FindMenuStructList(ctx, finder, params.Page)
	if err != nil {
		ginext.ResError(c, err)
		return
	}
	ginext.ResPage(c, result, *params.Page)
}
