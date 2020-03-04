package permservice

import (
	"errors"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//FindUserIdByOrgId 根据orgId,查找归属的UserId,不包括子部门
func FindUserIdByOrgId(dbConnection *zorm.DBConnection, orgId string, page *zorm.Page) ([]string, error) {
	if len(orgId) < 1 {
		return nil, errors.New("orgId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.userId FROM ").Append(permstruct.UserOrgStructTableName)
	finder.Append(" re where  re.orgId=? and re.managerType>0 order by re.managerType desc ", orgId)
	userIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, finder, &userIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userIds, nil
}

//FindUserByOrgId 根据orgId,查找归属的User,不包括子部门
func FindUserByOrgId(dbConnection *zorm.DBConnection, orgId string, page *zorm.Page) ([]permstruct.UserStruct, error) {
	userIds, errUserIds := FindUserIdByOrgId(dbConnection, orgId, page)
	if errUserIds != nil {
		return nil, errUserIds
	}
	return listUserId2ListUser(dbConnection, userIds)
}

//FindAllUserIdByOrgId 查询部门下所有的UserId,包括子部门
func FindAllUserIdByOrgId(dbConnection *zorm.DBConnection, orgId string, page *zorm.Page) ([]string, error) {
	if len(orgId) < 1 {
		return nil, errors.New("orgId不能为空")
	}
	orgStructPtr, errFindById := FindOrgStructById(dbConnection, orgId)
	if errFindById != nil || orgStructPtr == nil {
		return nil, errFindById
	}
	comcode := orgStructPtr.Comcode
	finder := zorm.NewFinder().Append("SELECT re.userId FROM ").Append(permstruct.UserOrgStructTableName).Append(" re,")
	finder.Append(permstruct.OrgStructTableName)
	finder.Append(" org WHERE org.id=re.orgId and org.comcode like ? and re.managerType>0   order by re.userId asc ", comcode+"%")

	userIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, finder, &userIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userIds, nil

}

//FindAllUserByOrgId 查询部门下所有的UserStruct,包括子部门
func FindAllUserByOrgId(dbConnection *zorm.DBConnection, orgId string, page *zorm.Page) ([]permstruct.UserStruct, error) {

	userIds, errUserIds := FindAllUserIdByOrgId(dbConnection, orgId, page)
	if errUserIds != nil {
		return nil, errUserIds
	}
	return listUserId2ListUser(dbConnection, userIds)

}

//FindOrgIdByUserId 根据userId查找所在的部门
func FindOrgIdByUserId(dbConnection *zorm.DBConnection, userId string, page *zorm.Page) ([]string, error) {
	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.orgId FROM  ").Append(permstruct.UserOrgStructTableName).Append(" re ")
	finder.Append("   WHERE re.userId=?    order by re.managerType desc   ", userId)

	orgIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, finder, &orgIds, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return orgIds, nil

}

//FindOrgByUserId 根据UserId,查找用户的部门对象
func FindOrgByUserId(dbConnection *zorm.DBConnection, userId string, page *zorm.Page) ([]permstruct.OrgStruct, error) {

	orgIds, errByUserId := FindOrgIdByUserId(dbConnection, userId, page)

	if errByUserId != nil {
		return nil, errByUserId
	}
	return listOrgId2ListOrg(dbConnection, orgIds)
}

func FindUserOrgByUserId(dbConnection *zorm.DBConnection, userId string, page *zorm.Page) ([]permstruct.UserOrgStruct, error) {
	if len(userId) < 1 {
		return nil, errors.New("userId不能为空")
	}
	finder := zorm.NewFinder().Append("SELECT re.* FROM  ").Append(permstruct.UserOrgStructTableName).Append(" re ")
	finder.Append("   WHERE re.userId=?    order by re.managerType desc   ", userId)

	userOrgs := make([]permstruct.UserOrgStruct, 0)
	errQueryList := zorm.QueryStructList(dbConnection, finder, &userOrgs, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	return userOrgs, nil

}

//listUserId2ListUser 根据 userIds 查询出 []permstruct.UserStruct
func listUserId2ListUser(dbConnection *zorm.DBConnection, userIds []string) ([]permstruct.UserStruct, error) {

	if len(userIds) < 1 {
		return nil, nil
	}

	users := make([]permstruct.UserStruct, 0)
	for _, userId := range userIds {
		user, errByUserId := FindUserStructById(dbConnection, userId)
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
func listOrgId2ListOrg(dbConnection *zorm.DBConnection, orgIds []string) ([]permstruct.OrgStruct, error) {
	if len(orgIds) < 1 {
		return nil, nil
	}

	orgs := make([]permstruct.OrgStruct, 0)
	for _, orgId := range orgIds {
		org, errByOrgId := FindOrgStructById(dbConnection, orgId)
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
