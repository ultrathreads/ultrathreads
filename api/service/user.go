package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/avatar"
	"ultrathreads/util/log"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/uploader"
)

type ScanUserCallback func(users []model.User)

type userRepository interface {
    Get(id int64) *model.User
}

func NewUserService(repo userRepository) *userService {
    return &userService{repo: repo}
}

type userService struct{
	repo userRepository
}

func (s *userService) Get(id int64) *model.User {
	return dao.UserDao.Get(id)
}

func (s *userService) GetBySlug(slug string) *model.User {
    id := hashid.Slug2Id[model.User](slug)

    user := cache.UserCache.Get(id)
    if user == nil {
        user = dao.UserDao.Get(id)
    }

    return user
}

func (s *userService) Take(where ...interface{}) *model.User {
	return dao.UserDao.Take(where...)
}

func (s *userService) Find(cnd *querybuilder.QueryBuilder) []model.User {
	return dao.UserDao.Find(cnd)
}

func (s *userService) FindOne(cnd *querybuilder.QueryBuilder) *model.User {
	return dao.UserDao.FindOne(cnd)
}

func (s *userService) List(cnd *querybuilder.QueryBuilder) (list []model.User, paging *querybuilder.Paging) {
	return dao.UserDao.List(cnd)
}

// Count 统计数量
func (s *userService) Count(cnd *querybuilder.QueryBuilder) int64 { // ✅ int → int64
	return dao.UserDao.Count(cnd)
}

func (s *userService) Update(dto form.UserUpdateForm) error {
	userID := hashid.Slug2Id[model.User](dto.Slug)
	err := dao.UserDao.Updates(userID, map[string]interface{}{
		"nickname":    dto.Nickname,
		"description": dto.Description,
		"level":       dto.Level,
		"update_time": util.NowTimestamp(),
	})
	cache.UserCache.Invalidate(userID)
	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := dao.UserDao.Updates(id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := dao.UserDao.UpdateColumn(id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

// Delete 删除用户
func (s *userService) Delete(id int64) error {
	if err := dao.UserDao.Delete(id); err != nil {
		return err
	}
	cache.UserCache.Invalidate(id)
	return nil
}

// Scan 游标分页扫描全表
func (s *userService) Scan(cb ScanUserCallback) {
	var cursor int64
	for {
		list := dao.UserDao.Find(querybuilder.NewQueryBuilder().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}

// Create 注册用户
func (s *userService) Create(username, email, nickname, password, rePassword string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	if len(username) == 0 {
		return nil, errors.New("用户名不能为空")
	}
	if err := util.IsValidatePassword(password, rePassword); err != nil {
		return nil, err
	}
	if len(email) == 0 {
		return nil, errors.New("请输入邮箱")
	}
	if err := util.IsValidateEmail(email); err != nil {
		return nil, err
	}
	if dao.UserDao.GetByEmail(email) != nil {
		return nil, errors.New("邮箱：" + email + " 已被占用")
	}
	if err := util.IsValidateUsername(username); err != nil {
		return nil, err
	}
	if s.isUsernameExists(username) {
		return nil, errors.New("用户名：" + username + " 已被占用")
	}

	user := &model.User{
		Username:   util.SqlNullString(username),
		Email:      util.SqlNullString(email),
		Nickname:   nickname,
		Password:   util.EncodePassword(password),
		Status:     model.StatusOk,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}

	// ✅ v2 事务 + 修复事务穿透
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		// 🔴 修复：原代码 dao.UserDao.Create 使用全局 db，不受 tx 控制
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		avatarUrl, err := s.HandleAvatar(user.ID, "")
		if err != nil {
			return err
		}

		updateColumns := map[string]interface{}{
			"avatar": avatarUrl,
		}
		if user.ID == 1 {
			updateColumns["level"] = model.UserLevelAdmin
		}

		// 🔴 修复：原代码 dao.UserDao.Updates 使用全局 db
		return tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(updateColumns).Error
	})

	if err != nil {
		return nil, err
	}
	cache.UserCache.Invalidate(user.ID)
	return user, nil
}

// SignInByLoginSource 第三方账号登录/注册
func (s *userService) SignInByLoginSource(loginSource *model.LoginSource) (*model.User, error) {
	user := s.Get(loginSource.UserID.Int64)
	if user != nil {
		if user.Status != model.StatusOk {
			return nil, errors.New("用户已被禁用")
		}
		return user, nil
	}

	var website, description string
	if loginSource.TargetType == model.LoginSourceTypeGithub {
		if blog := gjson.Get(loginSource.ExtraData, "blog"); blog.Exists() && len(blog.String()) > 0 {
			website = blog.String()
		} else if htmlUrl := gjson.Get(loginSource.ExtraData, "html_url"); htmlUrl.Exists() && len(htmlUrl.String()) > 0 {
			website = htmlUrl.String()
		}
		description = gjson.Get(loginSource.ExtraData, "bio").String()
	}

	user = &model.User{
		Username:    sql.NullString{},
		Nickname:    loginSource.Nickname,
		Status:      model.StatusOk,
		Website:     website,
		Description: description,
		CreateTime:  util.NowTimestamp(),
		UpdateTime:  util.NowTimestamp(),
	}

	// ✅ v2 事务 + 修复三处事务穿透
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		// 🔴 修复1: CreateUser
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// 🔴 修复2: UpdateLoginSource
		if err := tx.Model(&model.LoginSource{}).Where("id = ?", loginSource.ID).
			UpdateColumn("user_id", user.ID).Error; err != nil {
			return err
		}
		// 🔴 修复3: UpdateAvatar
		avatarUrl, err := s.HandleAvatar(user.ID, loginSource.Avatar)
		if err != nil {
			return err
		}
		return tx.Model(&model.User{}).Where("id = ?", user.ID).
			UpdateColumn("avatar", avatarUrl).Error
	})

	if err != nil {
		return nil, util.FromError(err)
	}
	cache.UserCache.Invalidate(user.ID)
	return user, nil
}

// HandleAvatar 处理头像
func (s *userService) HandleAvatar(userId int64, avatarUrl string) (string, error) {
	if len(avatarUrl) > 0 {
		return uploader.CopyImage(avatarUrl)
	}
	avatarBytes, err := avatar.Generate(userId)
	if err != nil {
		return "", err
	}
	return uploader.PutImage(avatarBytes)
}

func (s *userService) isEmailExists(email string) bool {
	if len(email) == 0 {
		return false
	}
	return dao.UserDao.GetByEmail(email) != nil
}

func (s *userService) isUsernameExists(username string) bool {
	return dao.UserDao.GetByUsername(username) != nil
}

// UpdateAvatar 更新头像
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// SetUsername 设置用户名（仅允许设置一次）
func (s *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := util.IsValidateUsername(username); err != nil {
		return err
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if len(user.Username.String) > 0 {
		return errors.New("你已设置了用户名，无法重复设置。")
	}
	if s.isUsernameExists(username) {
		return errors.New("用户名：" + username + " 已被占用")
	}
	return s.UpdateColumn(userId, "username", username)
}

// SetEmail 设置邮箱
func (s *userService) SetEmail(userId int64, email string) error {
	email = strings.TrimSpace(email)
	if err := util.IsValidateEmail(email); err != nil {
		return err
	}
	if s.isEmailExists(email) {
		return errors.New("邮箱：" + email + " 已被占用")
	}
	return s.UpdateColumn(userId, "email", email)
}

// SetPassword 首次设置密码
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := util.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	return s.UpdateColumn(userId, "password", util.EncodePassword(password))
}

// UpdatePassword 修改密码
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := util.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}
	if !util.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}
	return s.UpdateColumn(userId, "password", util.EncodePassword(password))
}

// IncrTopicCount post_count + 1
// ✅ 修复：原代码先查后写存在并发竞态，改为原子 SQL 自增
func (s *userService) IncrTopicCount(userId int64) int64 {
	if err := dao.UserDao.UpdateColumn(userId, "post_count", gorm.Expr("post_count + ?", 1)); err != nil {
		log.Error("IncrTopicCount failed: %v", err)
		return 0
	}
	cache.UserCache.Invalidate(userId)
	// 回查最新值返回
	user := dao.UserDao.Get(userId)
	if user == nil {
		return 0
	}
	return user.TopicCount
}

// IncrCommentCount comment_count + 1
// ✅ 修复：同上，改为原子 SQL 自增
func (s *userService) IncrCommentCount(userId int64) int64 {
	if err := dao.UserDao.UpdateColumn(userId, "comment_count", gorm.Expr("comment_count + ?", 1)); err != nil {
		log.Error("IncrCommentCount failed: %v", err)
		return 0
	}
	cache.UserCache.Invalidate(userId)
	user := dao.UserDao.Get(userId)
	if user == nil {
		return 0
	}
	return user.CommentCount
}

// SyncUserCount 同步用户计数
func (s *userService) SyncUserCount() {
	s.Scan(func(users []model.User) {
		for _, user := range users {
			topicCount := dao.PostDao.Count(querybuilder.NewQueryBuilder().Eq("user_id", user.ID).Eq("status", model.StatusOk))
			if err := dao.UserDao.UpdateColumn(user.ID, "post_count", topicCount); err != nil {
				log.Error("SyncUserCount update post_count failed for user %d: %v", user.ID, err)
			}
			cache.UserCache.Invalidate(user.ID)
		}
	})
}

var (
	errInvalidAccount = errors.New("账号或密码错误")
	errInvalidCode    = errors.New("请输入正确验证码")
	errAccountLocked  = errors.New("账号已被锁定,请联系管理员")
)

// VerifyAndReturnUserInfo 登录验证并返回用户信息
func (s *userService) VerifyAndReturnUserInfo(username, password string) (bool, error, model.User) {
	var userModel *model.User
	if err := util.IsValidateEmail(username); err == nil {
		userModel = dao.UserDao.GetByEmail(username)
	} else {
		userModel = dao.UserDao.GetByUsername(username)
	}

	if userModel == nil || userModel.ID < 1 {
		return false, errInvalidAccount, model.User{}
	}
	if !util.ValidatePassword(userModel.Password, password) {
		log.Error("password wrong: username=%s", userModel.Nickname)
		return false, errInvalidAccount, model.User{}
	}
	return true, nil, *userModel
}