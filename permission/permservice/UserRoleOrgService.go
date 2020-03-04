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

//listUserId2ListUser 根据 ListUserId 查询封装List<User> 对象
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
		users = append(users, user)
	}
	return users, nil

}
