package render

import (
	"ultrathreads/model"
	"ultrathreads/util/avatar"
	"ultrathreads/util/hashid"
)

// ToUser 纯函数：仅负责将 model.User 转换为响应结构体
// 无 I/O、无额外依赖
func ToUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}

	a := user.Avatar
	if len(a) == 0 {
		a = avatar.DefaultAvatar
	}

	levelName := "普通用户"
	if user.Level == model.UserLevelAdmin {
		levelName = "管理员"
	}

	ret := &model.UserInfo{
		Slug:         hashid.Id2Slug[model.User](user.ID),
		Username:     user.Username.String,
		Nickname:     user.Nickname,
		Avatar:       a,
		Level:        user.Level,
		LevelName:    levelName,
		Website:      user.Website,
		Description:  user.Description,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		PasswordSet:  len(user.Password) > 0,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt,
	}

	// 黑名单用户脱敏处理
	if user.Status == model.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Avatar = avatar.DefaultAvatar
		ret.Website = ""
		ret.Description = ""
	}

	return ret
}

// ToDefaultUser 替代原 ToUserDefaultIfNull 的兜底逻辑
func ToDefaultUser(id int64) *model.UserInfo {
	return &model.UserInfo{
		Slug:      hashid.Id2Slug[model.User](id),
		Nickname:  "系统通知",
		Avatar:    avatar.DefaultAvatar,
		LevelName: "普通用户",
	}
}

// ToUsers 批量转换
func ToUsers(users []model.User) []model.UserInfo {
	if len(users) == 0 {
		return []model.UserInfo{}
	}
	responses := make([]model.UserInfo, 0, len(users))
	for i := range users {
		if item := ToUser(&users[i]); item != nil {
			responses = append(responses, *item)
		}
	}
	return responses
}