package service

import (
	"errors"

	"github.com/snguovo/Go-000/Week02/tree/main/dao"
)

type Service struct {
}

func NewService() *Service {
	return new(Service)
}
func (s *Service) FindUserByID(userID int) (*dao.User, error) {
	return dao.FindUserByID(userID)
}

//MustFindUserByID 虚拟一个特殊场景 service 来处理找不到的情况
func (s *Service) MustFindUserByID(userID int) (user *dao.User, err error) {
	user, err = dao.FindUserByID(userID)
	if errors.Is(err, dao.ErrRecordNotFound) {
		user = dao.GetFakeUser()
		err = nil
	}
	return
}
