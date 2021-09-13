package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/poniteru/go-coin-watcher/app/dao"
	"strings"
)

func Register(chatId int64, inviteCode string) (token string, err error) {
	// 查用户是否存在并已激活
	exists, err := dao.SelectUserExists(chatId)
	if err != nil {
		return
	}
	if exists == 1 {
		err = errors.New("User is already exists!")
		return
	}
	// 查邀请码是否存在并有效
	available, err := dao.SelectInviteCodeAvailable(inviteCode)
	if err != nil {
		return
	}
	if available == 0 {
		err = errors.New("InviteCode is not available!")
		return
	}
	// 注册成功给用户分配token
	token = strings.ReplaceAll(uuid.New().String(), "-", "")
	err = dao.InsertUser(chatId, token)
	if err != nil {
		return
	}
	// 返回注册成功
	return
}
