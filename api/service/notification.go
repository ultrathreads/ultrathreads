package service

import (
	"sync"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/email"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

// NotificationService 通知业务契约
type NotificationService interface {
	Get(id int64) *model.Notification
	Take(where ...interface{}) *model.Notification
	Find(cnd *querybuilder.QueryBuilder) []model.Notification
	FindOne(cnd *querybuilder.QueryBuilder) *model.Notification
	List(cnd *querybuilder.QueryBuilder) ([]model.Notification, *querybuilder.Paging)
	Create(t *model.Notification) error
	Update(t *model.Notification) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetUnReadCount(userId int64) int64
	MarkRead(userId int64) error
	SendUserWatchNotification(userWatch *model.UserWatch)
	SendPostLikeNotification(postLike *model.PostLike)
	Produce(fromId, toId int64, content, quoteContent string, msgType int, extraDataMap map[string]interface{})
	Consume()
	SendEmailNotice(notification *model.Notification)
}

func NewNotificationService(repo repository.NotificationRepository, postRepo repository.PostRepository, userCache cache.UserCacheInterface, settingCache cache.SettingCacheInterface) NotificationService {
	return &notificationService{
		repo:              repo,
		postRepo:          postRepo,
		userCache:         userCache,
		settingCache:      settingCache,
		notificationsChan: make(chan *model.Notification),
	}
}

type notificationService struct {
	repo                     repository.NotificationRepository
	postRepo                 repository.PostRepository
	userCache                cache.UserCacheInterface
	settingCache             cache.SettingCacheInterface
	notificationsChan        chan *model.Notification
	notificationsConsumeOnce sync.Once
}

func (s *notificationService) Get(id int64) *model.Notification {
	return s.repo.Get(id)
}

func (s *notificationService) Take(where ...interface{}) *model.Notification {
	return s.repo.Take(where...)
}

func (s *notificationService) Find(cnd *querybuilder.QueryBuilder) []model.Notification {
	return s.repo.Find(cnd)
}

func (s *notificationService) FindOne(cnd *querybuilder.QueryBuilder) *model.Notification {
	return s.repo.FindOne(cnd)
}

func (s *notificationService) List(cnd *querybuilder.QueryBuilder) ([]model.Notification, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *notificationService) Create(t *model.Notification) error {
	return s.repo.Create(t)
}

func (s *notificationService) Update(t *model.Notification) error {
	return s.repo.Update(t)
}

func (s *notificationService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *notificationService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *notificationService) Delete(id int64) {
	s.repo.Delete(id)
}

func (s *notificationService) GetUnReadCount(userId int64) (count int64) {
	return s.repo.GetUnReadCount(userId)
}

func (s *notificationService) MarkRead(userId int64) error {
	return s.repo.UpdateStatusBatch(userId)
}

func (s *notificationService) SendUserWatchNotification(userWatch *model.UserWatch) {
	user := s.userCache.Get(userWatch.WatcherID)

	var (
		fromId       = userWatch.WatcherID
		authorId     int64
		content      string
		quoteContent string
	)

	authorId = userWatch.UserID
	content = user.Username.String + " 关注了你"
	quoteContent = ""

	if authorId <= 0 {
		return
	}
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeUserWatch, map[string]interface{}{
		"entityType":  model.EntityTypeUser,
		"entityId":    userWatch.WatcherID,
		"userWatchID": userWatch.ID,
	})
}

func (s *notificationService) SendPostLikeNotification(postLike *model.PostLike) {
	user := s.userCache.Get(postLike.UserId)

	var (
		fromId       = postLike.UserId
		authorId     int64
		content      string
		quoteContent string
	)
	post := s.postRepo.Get(postLike.PostId)
	if post != nil {
		authorId = post.UserId
		content = user.Username.String + " 点赞了你的话题：" + post.Title
		quoteContent = ""
	}

	if authorId <= 0 {
		return
	}
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypePostLike, map[string]interface{}{
		"entityType": model.EntityTypePost,
		"entityId":   post.ID,
		"postLikeId": postLike.ID,
	})
}

func (s *notificationService) Produce(fromId, toId int64, content, quoteContent string, msgType int, extraDataMap map[string]interface{}) {
	to := s.userCache.Get(toId)
	if to == nil {
		return
	}

	s.Consume()

	var (
		extraData string
		err       error
	)
	if extraData, err = util.FormatJson(extraDataMap); err != nil {
		log.Error("格式化extraData错误")
	}
	s.notificationsChan <- &model.Notification{
		FromId:       fromId,
		UserId:       toId,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
		ExtraData:    extraData,
		Status:       model.NotificationStatusUnread,
		CreateTime:   util.NowTimestamp(),
	}
}

func (s *notificationService) Consume() {
	s.notificationsConsumeOnce.Do(func() {
		go func() {
			log.Info("开始消费系统消息...")
			for {
				msg := <-s.notificationsChan
				log.Info("处理消息：from=%s to=%s", msg.FromId, msg.UserId)

				if err := s.Create(msg); err != nil {
					log.Info("创建消息发生异常...")
				} else {
					s.SendEmailNotice(msg)
				}
			}
		}()
	})
}

func (s *notificationService) SendEmailNotice(notification *model.Notification) {
	user := s.userCache.Get(notification.UserId)
	if user != nil && len(user.Email.String) > 0 {
		siteTitle := s.settingCache.GetValue(model.SettingSiteTitle)
		emailTitle := siteTitle + " 新消息提醒"
		email.SendTemplateEmail(user.Email.String, emailTitle, emailTitle, notification.Content,
			notification.QuoteContent, urls.AbsUrl("/user/notifications"))
		log.Info("发送邮件...email=%s", user.Email)
	} else {
		log.Info("邮件未发送，没设置邮箱...")
	}
}
