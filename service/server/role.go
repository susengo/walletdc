package server

import (
	"github.com/go-xorm/xorm"
	"github.com/susengo/commontools/gintool"
	"github.com/susengo/swing/model"
)

type RoleService struct {
	DbEngine *xorm.Engine
}

func (l *RoleService) Add(role *model.Role) (bool, string) {

	i, err := l.DbEngine.Insert(role)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *RoleService) Update(role *model.Role) (bool, string) {
	i, err := l.DbEngine.Where("rkey = ?", role.Rkey).Update(role)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *RoleService) Delete(key string) (bool, string) {
	i, err := l.DbEngine.Where("rkey = ?", key).Delete(&model.Role{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

func (l *RoleService) GetByRole(role *model.Role) (bool, *model.Role) {
	has, err := l.DbEngine.Get(role)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, role
}

func (l *RoleService) GetList(role *model.Role, page, size int) (bool, []*model.Role, int64) {

	pager := gintool.CreatePager(page, size)

	roles := make([]*model.Role, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if role.Name != "" {
		where += " and name like ? "
		values = append(values, "%"+role.Name+"%")
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}
	total, err := l.DbEngine.Where(where, values...).Count(new(model.Role))
	if err != nil {
		logger.Error(err.Error())
	}

	return true, roles, total
}

func (l *RoleService) GetAll() (bool, []*model.Role) {

	roles := make([]*model.Role, 0)

	err := l.DbEngine.Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}
	return true, roles
}

func NewRoleService(engine *xorm.Engine) *RoleService {
	return &RoleService{
		DbEngine: engine,
	}
}
