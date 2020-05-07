package service

import (
	"ginana-blog/internal/model"
)

func (s *service) GetLinks() (links []*model.Link, err error) {
	key := s.hm.GetCacheKey(7)
	err = s.mc.Get(key, &links)
	if err != nil {
		if err = s.db.Model(&links).Order("created_at desc").Find(&links).Error; err != nil {
			return nil, s.hm.GetMessage(1001, err)
		}
		if err = s.mc.Set(key, &links); err != nil {
			return nil, s.hm.GetMessage(1002, err)
		}
	}
	return
}
