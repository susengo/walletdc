package server

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/susengo/commontools/gintool"
	jwttool "github.com/susengo/commontools/jwt"
	"github.com/susengo/commontools/log"
	"github.com/susengo/commontools/password"
	"github.com/susengo/swing/model"
)

var logger = log.GetLogger("server", log.ERROR)

const TokenKey = "baas user secret"

type UserService struct {
	DbEngine *xorm.Engine
}

func (l *UserService) Add(user *model.User) (bool, string) {
	if user.Password != "" {
		user.Password = password.Encode(user.Password, 12, "default")
	}
	i, err := l.DbEngine.Insert(user)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *UserService) Update(user *model.User) (bool, string) {
	if user.Password != "" {
		user.Password = password.Encode(user.Password, 12, "default")
	}
	i, err := l.DbEngine.Where("id = ?", user.Id).Update(user)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *UserService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&model.User{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

func (l *UserService) GetByUser(user *model.User) (bool, *model.User) {
	has, err := l.DbEngine.Get(user)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, user
}

func (l *UserService) GetList(user *model.User, page, size int) (bool, []model.UserDetail, int64) {

	pager := gintool.CreatePager(page, size)

	users := make([]*model.User, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if user.Account != "" {
		where += " and account = ? "
		values = append(values, user.Account)
	}
	if user.Name != "" {
		where += " and name like ? "
		values = append(values, "%"+user.Name+"%")
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&users)
	if err != nil {
		logger.Error(err.Error())
	}

	total, err := l.DbEngine.Where(where, values...).Count(new(model.User))
	if err != nil {
		logger.Error(err.Error())
	}

	userIds := make([]int, len(users))
	userDatas := make([]model.UserDetail, len(users))
	for i, u := range users {
		userIds[i] = u.Id
		userDatas[i].Id = u.Id
		userDatas[i].Account = u.Account
		userDatas[i].Password = u.Password
		userDatas[i].Avatar = u.Avatar
		userDatas[i].Name = u.Name
		userDatas[i].Created = u.Created
	}

	roles := make([]model.UserRole, 0)
	err = l.DbEngine.In("user_id", userIds).Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}

	for i, d := range userDatas {
		keys := make([]string, 0)
		for _, r := range roles {
			if r.UserId == d.Id {
				keys = append(keys, r.RoleKey)
			}
		}
		d.Roles = keys
		userDatas[i] = d
	}

	return true, userDatas, total
}

func (l *UserService) GetToken(user *model.User) *model.JwtToken {

	info := make(map[string]interface{})
	now := time.Now()
	info["userId"] = user.Id
	info["exp"] = now.Add(time.Hour * 1).Unix() // 1 小时过期
	info["iat"] = now.Unix()
	tokenString := jwttool.CreateToken(TokenKey, info)

	return &model.JwtToken{
		Token: tokenString,
	}
}

func (l *UserService) CheckToken(token string, user *model.User) (*model.UserInfo, error) {

	info, ok := jwttool.ParseToken(token, TokenKey)
	infoMap := info.(jwt.MapClaims)
	if ok {
		expTime := infoMap["exp"].(float64)
		if float64(time.Now().Unix()) >= expTime {
			return nil, fmt.Errorf("%s", "token已过期")
		} else {
			l.DbEngine.Get(user)
			ur := make([]model.UserRole, 0)
			err := l.DbEngine.Where("user_id = ?", user.Id).Find(&ur)
			if err != nil {
				logger.Error(err.Error())
			}
			roles := make([]string, len(ur))
			for i, m := range ur {
				roles[i] = m.RoleKey
			}
			info := &model.UserInfo{
				Avatar:  user.Avatar,
				Roles:   roles,
				Name:    user.Name,
				Account: user.Account,
			}
			return info, nil
		}
	} else {
		return nil, fmt.Errorf("%s", "token无效")
	}
}

func (l *UserService) AddAuth(ur *model.UserRole) (bool, string) {

	i, err := l.DbEngine.Insert(ur)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *UserService) DelAuth(ur *model.UserRole) (bool, string) {

	i, err := l.DbEngine.Delete(ur)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "del success"
	}
	return false, "del fail"
}

func (l *UserService) HasAdminRole(account string) bool {
	user := &model.User{Account: account}
	_, user = l.GetByUser(user)

	ur := make([]model.UserRole, 0)
	err := l.DbEngine.Where("user_id = ?", user.Id).Find(&ur)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, m := range ur {
		if m.RoleKey == "admin" {
			return true
		}
	}
	return false
}
func NewUserService(engine *xorm.Engine) *UserService {
	return &UserService{
		DbEngine: engine,
	}
}
