package dao

import (
	"github.com/poniteru/go-coin-watcher/common/db"
)

func InsertUser(chatId int64, token string) (err error) {
	sqlStr := "INSERT INTO tg_user_info (chat_id, token, status) VALUES (?,?,1)"
	_, err = db.Instance().Exec(sqlStr, chatId, token)
	return
}

func SelectUserChatId(token string) (chatId int64, err error) {
	sqlStr := "select chat_id from tg_user_info where token=?"
	err = db.Instance().Get(&chatId, sqlStr, token)
	return
}

func SelectUserExists(chatId int64) (exists int8, err error) {
	sqlStr := "SELECT EXISTS (SELECT 1 FROM tg_user_info WHERE chat_id=?)"
	err = db.Instance().Get(&exists, sqlStr, chatId)
	return
}

func SelectInviteCodeAvailable(inviteCode string) (exists int8, err error) {
	sqlStr := "SELECT EXISTS (SELECT 1 FROM tg_invite_code WHERE invite_code=? AND status=1)"
	err = db.Instance().Get(&exists, sqlStr, inviteCode)
	return
}
