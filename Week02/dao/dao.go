package dao

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

var (
	//DB *sql.DB

	//ErrRecordNotFound 定义一个标准错误
	ErrRecordNotFound = errors.New("record not found")
)

//User 用户对象
type User struct {
	ID   int
	Name string
}

//FindUserByID 通过id查询user
func FindUserByID(userID int) (u *User, err error) {
	err = sql.ErrNoRows
	return u, errors.Wrap(err, fmt.Sprintf("find user by id:%v failed", userID))
}

//GetFakeUser 伪造一个user
func GetFakeUser() (u *User) {
	return
}
