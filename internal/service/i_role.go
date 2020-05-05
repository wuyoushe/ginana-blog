package service

import (
	"context"
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/database"
	"sync"
)

func (s *service) GetEFRoles(c context.Context) (roles []*database.EFRolePolicy, err error) {
	var roleIdList []int64
	s.db.Model(&model.Role{}).Select("id").Pluck("id", &roleIdList)
	var wg sync.WaitGroup
	var ch = make(chan int64, 1)
	wg.Add(len(roleIdList))
	for _, roleId := range roleIdList {
		go func(roleId int64, roles *[]*database.EFRolePolicy, wg *sync.WaitGroup) {
			r, err := s.GetRole(roleId)
			if err != nil {
				return
			}
			ch <- roleId
			for _, policy := range r.Policys {
				role := new(database.EFRolePolicy)
				role.RoleName = r.RoleName
				role.Router = policy.Router
				role.Method = policy.Method
				*roles = append(*roles, role)
			}
			<-ch
			wg.Done()
		}(roleId, &roles, &wg)
	}
	wg.Wait()
	return
}

func (s *service) GetRole(id int64) (role *model.Role, err error) {
	key := s.hm.GetCacheKey(2, id)
	role = new(model.Role)
	err = s.mc.Get(key, role)
	if err != nil {
		role.ID = id
		if err = s.db.Find(role).Related(&role.Policys, "Policys").Error; err != nil {
			err = s.hm.GetMessage(1001, err)
			return
		}
		if err = s.mc.Set(key, role); err != nil {
			err = s.hm.GetMessage(1002, err)
			return
		}
	}
	return
}
