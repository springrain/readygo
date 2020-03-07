package permservice

import (
	"context"
	"errors"
	"readygo/permission/permstruct"
	"strconv"

	"gitee.com/chunanyong/zorm"
)

//FindUserIdByOrgId 根据orgId,查找归属的UserId,不包括子部门,不包括会员
func FindUserIdByOrgId(ctx context.Context, orgId string, page *zorm.Page) ([]string, error) {
	if len(orgId) < 1 {
		return nil, errors.New("orgId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.userId FROM ").Append(permstruct.UserOrgStructTableName)
	finder.Append(" re where  re.orgId=? and re.managerType>0 order by re.managerType desc ", orgId)
	userIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &userIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userIds, nil
}

//FindUserByOrgId 根据orgId,查找归属的User,不包括子部门,不包括会员
func FindUserByOrgId(ctx context.Context, orgId string, page *zorm.Page) ([]permstruct.UserStruct, error) {
	userIds, errUserIds := FindUserIdByOrgId(ctx, orgId, page)
	if errUserIds != nil {
		return nil, errUserIds
	}
	return listUserId2ListUser(ctx, userIds)
}

//FindAllUserIdByOrgId 查询部门下所有的UserId,包括子部门,不包括会员
func FindAllUserIdByOrgId(ctx context.Context, orgId string, page *zorm.Page) ([]string, error) {
	if len(orgId) < 1 {
		return nil, errors.New("orgId不能为空")
	}
	orgStructPtr, errFindById := FindOrgStructById(ctx, orgId)
	if errFindById != nil || orgStructPtr == nil {
		return nil, errFindById
	}
	comcode := orgStructPtr.Comcode
	finder := zorm.NewFinder().Append("SELECT re.userId FROM ").Append(permstruct.UserOrgStructTableName).Append(" re,")
	finder.Append(permstruct.OrgStructTableName)
	finder.Append(" org WHERE org.id=re.orgId and org.comcode like ? and re.managerType>0   order by re.userId asc ", comcode+"%")

	userIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &userIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userIds, nil

}

//FindAllUserByOrgId 查询部门下所有的UserStruct,包括子部门,不包括会员
func FindAllUserByOrgId(ctx context.Context, orgId string, page *zorm.Page) ([]permstruct.UserStruct, error) {

	userIds, errUserIds := FindAllUserIdByOrgId(ctx, orgId, page)
	if errUserIds != nil {
		return nil, errUserIds
	}
	return listUserId2ListUser(ctx, userIds)

}

//FindOrgIdByUserId 根据userId查找所在的部门
func FindOrgIdByUserId(ctx context.Context, userId string, page *zorm.Page) ([]string, error) {
	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.orgId FROM  ").Append(permstruct.UserOrgStructTableName).Append(" re ")
	finder.Append("   WHERE re.userId=?    order by re.managerType desc   ", userId)

	orgIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &orgIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return orgIds, nil

}

//FindOrgByUserId 根据UserId,查找用户的部门对象
func FindOrgByUserId(ctx context.Context, userId string, page *zorm.Page) ([]permstruct.OrgStruct, error) {

	orgIds, errByUserId := FindOrgIdByUserId(ctx, userId, page)

	if errByUserId != nil {
		return nil, errByUserId
	}
	return listOrgId2ListOrg(ctx, orgIds)
}

//FindUserOrgByUserId 根据userId查找部门UserOrg中间表对象
func FindUserOrgByUserId(ctx context.Context, userId string, page *zorm.Page) ([]permstruct.UserOrgStruct, error) {
	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.* FROM  ").Append(permstruct.UserOrgStructTableName).Append(" re ")
	finder.Append("   WHERE re.userId=?    order by re.managerType desc   ", userId)

	userOrgs := make([]permstruct.UserOrgStruct, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &userOrgs, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userOrgs, nil

}

//FindManagerOrgIdByUserId 根据userId查找管理的部门ID,不包括角色扩展的部门
func FindManagerOrgIdByUserId(ctx context.Context, userId string, page *zorm.Page) ([]string, error) {
	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}

	finder := zorm.NewFinder().Append("SELECT re.orgId FROM  ").Append(permstruct.UserOrgStructTableName)
	finder.Append(" re  WHERE re.userId=?  and re.managerType=2  order by re.orgId desc   ", userId)

	orgIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &orgIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return orgIds, nil

}

//FindManagerOrgByUserId 根据userId查找管理的部门对象,不包括角色扩展的部门
func FindManagerOrgByUserId(ctx context.Context, userId string, page *zorm.Page) ([]permstruct.OrgStruct, error) {

	orgIds, errByUserId := FindManagerOrgIdByUserId(ctx, userId, page)

	if errByUserId != nil {
		return nil, errByUserId
	}
	return listOrgId2ListOrg(ctx, orgIds)

}

//FindManagerUserIdByOrgId 查找部门主管的UserId,一个部门只有一个主管,其他管理人员可以通过角色进行扩展分配
func FindManagerUserIdByOrgId(ctx context.Context, orgId string) (string, error) {
	if len(orgId) < 1 {
		return "", errors.New("orgId不能为空")
	}

	finder := zorm.NewFinder().Append("SELECT re.userId FROM ").Append(permstruct.UserOrgStructTableName)
	finder.Append(" re  WHERE   re.orgId=? and  re.managerType=2   order by re.userId desc   ", orgId)

	managerUserId := ""
	errQueryStruct := zorm.QueryStruct(ctx, finder, &managerUserId)
	if errQueryStruct != nil {
		return "", errQueryStruct
	}

	return managerUserId, nil
}

//FindManagerUserByOrgId 查找部门主管的UserStruct对象,一个部门只有一个主管,其他管理人员可以通过角色进行扩展分配.调用  FindManagerUserIdByOrgId 方法
func FindManagerUserByOrgId(ctx context.Context, orgId string) (*permstruct.UserStruct, error) {

	managerUserId, errByOrgId := FindManagerUserIdByOrgId(ctx, orgId)

	if errByOrgId != nil {
		return nil, errByOrgId
	}

	return FindUserStructById(ctx, managerUserId)

}

//FindAllUserCountByOrgId 查询部门下所有的人员数量,包括子部门.不包括会员
func FindAllUserCountByOrgId(ctx context.Context, orgId string) (int, error) {
	if len(orgId) < 1 {
		return -1, errors.New("orgId不能为空")
	}

	org, errByOrgId := FindOrgStructById(ctx, orgId)

	if errByOrgId != nil {
		return -1, errByOrgId
	}

	finder := zorm.NewFinder().Append("SELECT count(re.userId) FROM ").Append(permstruct.UserOrgStructTableName)
	finder.Append(" re,").Append(permstruct.OrgStructTableName)
	finder.Append(" org WHERE org.id=re.orgId and org.comcode like ? and  re.managerType>0 ", org.Comcode+"%")
	//查询总条数
	count := -1
	errCount := zorm.QueryStruct(ctx, finder, &count)
	return count, errCount
}

//FindRoleOrgByRoleId 根据roleId查找roleOrg,角色管理的部门,用于角色自定的部门范围,查询 t_role_org 中间表
func FindRoleOrgByRoleId(ctx context.Context, roleId string, page *zorm.Page) ([]permstruct.RoleOrgStruct, error) {
	if len(roleId) < 1 {
		return nil, errors.New("roleId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.* FROM ").Append(permstruct.RoleOrgStructTableName)
	finder.Append(" re WHERE re.roleId=? order by re.id desc ", roleId)

	roleOrgs := make([]permstruct.RoleOrgStruct, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &roleOrgs, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return roleOrgs, nil
}

//WrapOrgIdFinderByUserId 查询用户有权限管理的所有部门,包括角色关联分配产生的部门权限.分装成Finder的形式用于关联查询的finder实体类
// 1.获取用户所有的 []permstruct.UserRoleStruct,包含主管的部门和角色分配的部门,不包括角色私有部门
// 2.wrapOrgIdFinderByUserRole(list) 生成完整的Finder对象.
func WrapOrgIdFinderByUserId(ctx context.Context, userId string) (*zorm.Finder, error) {
	// 获取用户所有的 []permstruct.UserRoleStruct,包含主管的部门和角色分配的部门,不包括角色私有部门
	roleOrgs, errByUserId := findManagerOrgAndRoleOrgByUserId(ctx, userId)
	if errByUserId != nil {
		return nil, errByUserId
	}
	// 生成 Finder 方法.
	return wrapOrgIdFinderByRoleOrg(ctx, roleOrgs)
}

//WrapOrgIdFinderByPrivateOrgRoleId 查询私有部门角色的部门范围,构造成Finder,用于权限范围限制
//在请求的时候,通过url查询出menuId和对应的roleId,根据RoleId可以知晓是私有还是公共权限,分开调用Finder
func WrapOrgIdFinderByPrivateOrgRoleId(ctx context.Context, roleId string, userId string) (*zorm.Finder, error) {
	if len(roleId) < 1 || len(userId) < 1 {
		return nil, errors.New("roleId或userId不能为空")
	}
	role, errByRoleId := FindRoleStructById(ctx, roleId)
	if errByRoleId != nil {
		return nil, errByRoleId
	}
	if role == nil || role.PrivateOrg == 0 { //只处理 私有部门 类型的角色.
		return nil, errors.New("只处理私有部门类型的角色")
	}

	// 当前用户 如果是主管,所管理的所有部门,如果是私有权限,就严格执行,不处理用户为主管的部门,部门主管可能没有部门权限.
	//  List<RoleOrg> list = wrapManagerRoleOrgByUserId(userId);
	roleOrgs := make([]permstruct.RoleOrgStruct, 0)
	//角色分配的 私有部门
	findRoleOrgIdByRole, errByRole := findRoleOrgIdByRole(ctx, role, userId, nil)
	if errByRole != nil {
		return nil, errByRole
	}
	// 把主管部门
	if len(findRoleOrgIdByRole) > 0 {
		roleOrgs = append(roleOrgs, findRoleOrgIdByRole...)
	}
	// 生成 Finder 方法.
	return wrapOrgIdFinderByRoleOrg(ctx, roleOrgs)
}

// 查询用户有权限管理的所有部门,包括角色关联分配产生的部门权限.分装成Finder的形式用于关联查询的finder实体类
//这个是基础的Finder 封装方法,其他的方式也在转化为List<UserRole> list,调用此方法
// 1.基于 List<UserRole> list 拼装WHERE条件
// 2.完善 前面的查询语句,构造完整的Finder 查询语句
func wrapOrgIdFinderByRoleOrg(ctx context.Context, list []permstruct.RoleOrgStruct) (*zorm.Finder, error) {

	// 基于 list []permstruct.RoleOrgStruct 拼装WHERE条件
	finder, errSQL := wrapOrgIdWheresSQLFinder(ctx, list)
	if errSQL != nil {
		return nil, errSQL
	}
	// 完善 前面的查询语句,构造完整的Finder 查询语句
	return wrapOrgIdFinder(ctx, finder)
}

// 查询用户根据角色派生和自身主管的所有部门
func findManagerOrgAndRoleOrgByUserId(ctx context.Context, userId string) ([]permstruct.RoleOrgStruct, error) {

	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}

	// 如果是主管,其管理的部门
	list, errRoleOrgByUserId := wrapManagerRoleOrgByUserId(ctx, userId)
	if errRoleOrgByUserId != nil {
		return nil, errRoleOrgByUserId
	}

	// 查询用户所有的角色
	listRole, errRoleByUserId := FindRoleByUserId(ctx, userId, nil)
	if errRoleByUserId != nil {
		return nil, errRoleByUserId
	}
	if len(listRole) < 1 {
		return list, nil
	}

	// 查询角色被分配到部门权限,不处理私有部门类型的角色
	for _, role := range listRole {
		if role.PrivateOrg == 1 { // 不处理 私有部门 类型的角色
			continue
		}
		//角色管理的部门
		findRoleOrgIdByRole, errByRole := findRoleOrgIdByRole(ctx, &role, userId, nil)
		if errByRole != nil {
			return nil, errByRole
		}
		//添加到List
		if len(findRoleOrgIdByRole) > 0 {
			list = append(list, findRoleOrgIdByRole...)
		}
	}

	return list, nil
}

// 根据role 对象 查询 Role的关联部门.  roleOrgType 0自己的数据,1所在部门,2所在部门及子部门数据,3.自定义部门数据.部门主管有所管理部门的数据全权限,无论角色是否分配
//  外围需要单独判断是否启用私有角色,不然很容易造成群贤扩大
//  这里只处理角色产生的权限,不考虑用户如果是主管派生的下级部门权限,这种情况有业务自己处理
func findRoleOrgIdByRole(ctx context.Context, role *permstruct.RoleStruct, userId string, page *zorm.Page) ([]permstruct.RoleOrgStruct, error) {

	if role == nil {
		return nil, errors.New("role不能为空")
	}

	// 角色部门类型 roleOrgType 0自己的数据,1所在部门,2所在部门及子部门数据,3.自定义部门数据
	roleOrgType := role.RoleOrgType

	if roleOrgType == 0 { // 用户自己的数据,不包含部门权限
		return nil, nil
	} else if roleOrgType == 1 || roleOrgType == 2 { //用户所在的部门

		if len(userId) < 1 {
			return nil, errors.New("userId不能为空")
		}
		//用户所在的部门,这里只处理角色产生的权限,不考虑用户如果是主管派生的下级部门权限,这种情况有业务自己处理
		orgIdByUserId, errByUserId := FindOrgIdByUserId(ctx, userId, page)
		if errByUserId != nil {
			return nil, errByUserId
		}
		if len(orgIdByUserId) < 1 {
			return nil, nil
		}
		list := make([]permstruct.RoleOrgStruct, 0)
		for _, orgId := range orgIdByUserId {
			re := permstruct.RoleOrgStruct{}
			// 如果是包含子部门权限
			if roleOrgType == 2 {
				re.Children = 1
			} else { // 不包含子部门的权限
				re.Children = 0
			}
			re.OrgId = orgId
			list = append(list, re)
		}
		return list, nil
	} else if roleOrgType == 3 { // 自定义权限
		return findOrgByRoleId(ctx, role.Id, page)
	}
	return nil, nil
}

//findOrgByRoleId 根据roleId,查询role下管理的部门,用于角色自定的部门范围,查询 t_role_org 中间表
func findOrgByRoleId(ctx context.Context, roleId string, page *zorm.Page) ([]permstruct.RoleOrgStruct, error) {
	if len(roleId) < 1 {
		return nil, errors.New("roleId不能为空")
	}

	finder := zorm.NewFinder().Append("SELECT re.* FROM ").Append(permstruct.RoleOrgStructTableName)
	finder.Append(" re WHERE re.roleId=? order by re.id desc ", roleId)

	roleOrgs := make([]permstruct.RoleOrgStruct, 0)
	errQueryList := zorm.QueryStructList(ctx, finder, &roleOrgs, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return roleOrgs, nil
}

//  查询用户作为主管时所有的管理部门,封装成 []permstruct.RoleOrgStruct 格式
func wrapManagerRoleOrgByUserId(ctx context.Context, userId string) ([]permstruct.RoleOrgStruct, error) {
	list := make([]permstruct.RoleOrgStruct, 0)

	// 查询用户直接管理的部门
	managerOrgIdByUserId, errByUserId := FindManagerOrgIdByUserId(ctx, userId, nil)
	if errByUserId != nil {
		return nil, errByUserId
	}
	if len(managerOrgIdByUserId) < 1 {
		return list, nil
	}
	// 构造List<RoleOrg>
	for _, orgId := range managerOrgIdByUserId {
		re := permstruct.RoleOrgStruct{}
		// 部门主管能管理子部门
		re.Children = 1
		re.OrgId = orgId
		list = append(list, re)
	}
	return list, nil
}

//wrapOrgIdFinder 构造完整的finder,基于 wrapOrgIdWheresSQLFinder 产生的 WHERE 条件,完善Finder的查询部分.
func wrapOrgIdFinder(ctx context.Context, whereFinder *zorm.Finder) (*zorm.Finder, error) {
	if whereFinder == nil {
		return nil, errors.New("whereFinder不能为空")
	}

	wheresql, errWhereSQL := whereFinder.GetSQL()
	if errWhereSQL != nil {
		return nil, errWhereSQL
	}
	if len(wheresql) < 1 {
		return nil, errors.New("whereFinder的SQL语句不能为空")
	}

	/*
	   // 查找人员
	   Finder finder = new Finder("SELECT  _system_temp_user_org.userId FROM ");
	   finder.append(Finder.getTableName(UserOrg.class)).append(" _system_temp_user_org,")
	           .append(Finder.getTableName(Org.class)).append(" _system_temp_org ");
	   finder.append(" WHERE _system_temp_user_org.orgId=_system_temp_org.id  ");
	*/

	// 查找部门
	finder := zorm.NewFinder().Append(" SELECT _system_temp_org.id  FROM ").Append(permstruct.OrgStructTableName)
	finder.Append(" _system_temp_org WHERE 1=1 ")

	// 增加 WHERE 条件
	finder.AppendFinder(whereFinder)

	return finder, nil
}

//wrapOrgIdWheresSQLFinder 基于 list []permstruct.RoleOrgStruct 生成 Finder 对象,并不是完整的语句,只是 WHERE 后面的部门条件语句 类似 and ( 1=2 or ....
func wrapOrgIdWheresSQLFinder(ctx context.Context, list []permstruct.RoleOrgStruct) (*zorm.Finder, error) {
	if len(list) < 1 {
		return nil, nil
	}
	//去掉重复的对象,对象是使用 id 作为对比字段的,需要把id设置好进行去重
	roleOrgMap := make(map[string]permstruct.RoleOrgStruct)

	for _, re := range list {
		id := re.OrgId + "_" + strconv.Itoa(re.Children)
		_, ok := roleOrgMap[id]
		if ok { //已经存在
			continue
		}
		r := permstruct.RoleOrgStruct{}

		r.Id = id
		r.OrgId = re.OrgId
		r.Children = re.Children
		roleOrgMap[id] = r
	}

	// 不包含子部门的 部门Id List
	noChildrenList := make([]string, 0)
	// 包含子部门的Finder
	hasChildrenFinder := zorm.NewFinder().Append("  and ( 1=2 ")

	for _, re := range roleOrgMap {
		orgId := re.OrgId
		children := re.Children
		if children == 0 { // 不包含子部门
			noChildrenList = append(noChildrenList, orgId)
		} else if children == 1 { // 包含子部门
			org, errByOrgId := FindOrgStructById(ctx, orgId)
			if errByOrgId != nil {
				return nil, errByOrgId
			}
			hasChildrenFinder.Append(" or _system_temp_org.comcode like ? ", org.Comcode+"%")
		}
	}

	// 处理没有子部门,部门Id or in 到 noChildrenList
	if len(noChildrenList) > 0 {
		// 前面有sql加连接符
		hasChildrenFinder.Append(" or  _system_temp_org.id in (?) ", noChildrenList)
	}

	hasChildrenFinder.Append(") ")

	return hasChildrenFinder, nil
}

//listUserId2ListUser 根据 userIds 查询出 []permstruct.UserStruct
func listUserId2ListUser(ctx context.Context, userIds []string) ([]permstruct.UserStruct, error) {

	if len(userIds) < 1 {
		return nil, nil
	}

	users := make([]permstruct.UserStruct, 0)
	for _, userId := range userIds {
		user, errByUserId := FindUserStructById(ctx, userId)
		if errByUserId != nil {
			return nil, errByUserId
		}
		if user == nil {
			continue
		}
		users = append(users, *user)
	}
	return users, nil

}

//listOrgId2ListOrg  根据 orgIds 查询出 []permstruct.OrgStruct
func listOrgId2ListOrg(ctx context.Context, orgIds []string) ([]permstruct.OrgStruct, error) {
	if len(orgIds) < 1 {
		return nil, nil
	}

	orgs := make([]permstruct.OrgStruct, 0)
	for _, orgId := range orgIds {
		org, errByOrgId := FindOrgStructById(ctx, orgId)
		if errByOrgId != nil {
			return nil, errByOrgId
		}
		if org == nil {
			continue
		}
		orgs = append(orgs, *org)
	}
	return orgs, nil

}

/**


 @Override
   public String updateUserOrg( UserOrg userOrg) throws Exception {
        Integer managerType = userOrg.getManagerType();

        if (userOrg==null||StringUtils.isBlank(userOrg.getOrgId())||managerType==null){
                return  "数据不能为空";
            }

        Finder finder=Finder.getDeleteFinder(UserOrg.class).append(" WHERE userId=:userId and orgId=:orgId ");
        finder.setParam("userId",userOrg.getUserId()).setParam("orgId",userOrg.getOrgId());
        super.update(finder);

        if(managerType<0){// 删除关系
           return null;
        }

        Date now=new Date();
        userOrg.setId(SecUtils.getUUID());
        userOrg.setCreateTime(now);
        userOrg.setUpdateTime(now);
        userOrg.setCreateUserId(SessionUser.getUserId());
        userOrg.setUpdateUserId(SessionUser.getUserId());

        super.save(userOrg);

        return null;
    }



    @Override
    public String updateRoleOrg(RoleOrg roleOrg) throws Exception {

        if(roleOrg==null||StringUtils.isBlank(roleOrg.getOrgId())||StringUtils.isBlank(roleOrg.getRoleId())||roleOrg.getCheck()==null){
            return  "数据不能为空";
        }

        Finder finder=Finder.getDeleteFinder(RoleOrg.class).append(" WHERE roleId=:roleId and orgId=:orgId ");
        finder.setParam("roleId",roleOrg.getRoleId()).setParam("orgId",roleOrg.getOrgId());
        super.update(finder);

        if(roleOrg.getCheck()==false){// 删除关系
            return null;
        }

        Date now=new Date();
        roleOrg.setId(SecUtils.getUUID());
        roleOrg.setCreateTime(now);
        roleOrg.setUpdateTime(now);
        roleOrg.setCreateUserId(SessionUser.getUserId());
        roleOrg.setUpdateUserId(SessionUser.getUserId());
        super.save(roleOrg);
        return null;
    }

    @Override
    public boolean isUserInOrg(String userId, String orgId) throws Exception {

        if (StringUtils.isBlank(userId) || StringUtils.isBlank(orgId)) {
            return false;
        }
        Finder finder = Finder.getSelectFinder(UserOrg.class, " 1 ").append(" WHERE userId=:userId and orgId=:orgId ");
        finder.setParam("userId", userId).setParam("orgId", orgId);
        Integer isUserInOrg = super.queryForObject(finder, Integer.class);

        if (isUserInOrg == null) {
            return false;
        }
        return true;
    }

	**/
