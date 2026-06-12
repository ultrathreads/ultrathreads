package converter

import (
	"strconv"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/avatar"
	"ultrathreads/util/hashid"
)

func ToUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.ID = id
		user.Username = util.SqlNullString(strconv.FormatInt(id, 10))
		user.Avatar = avatar.DefaultAvatar
		user.CreateTime = util.NowTimestamp()
	}
	return ToUser(user)
}

func ToUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	return ToUser(user)
}

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
	slug := hashid.Id2Slug[model.User](user.ID)
	ret := &model.UserInfo{
		Slug: 		  slug,
		Username:     user.Username.String,
		Nickname:     user.Nickname,
		Avatar:       a,
		Level:        user.Level,
		LevelName:    levelName,
		Website:      user.Website,
		Description:  user.Description,
		Score:        0, // 占位，下方按状态赋值
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		PasswordSet:  len(user.Password) > 0,
		Status:       user.Status,
		CreateTime:   user.CreateTime,
	}
	if user.Status == model.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Avatar = avatar.DefaultAvatar
		ret.Website = ""
		ret.Description = ""
	} else {
		ret.Score = cache.UserCache.GetScore(user.ID)
	}
	return ret
}

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