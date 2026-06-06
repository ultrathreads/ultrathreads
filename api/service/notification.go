package service

import (
	"sync"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/email"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

var NotificationService = newNotificationService()

func newNotificationService() *notificationService {
	return &notificationService{
		notificationsChan: make(chan *model.Notification),
	}
}

type notificationService struct {
	notificationsChan        chan *model.Notification
	notificationsConsumeOnce sync.Once
}

func (s *notificationService) Get(id int64) *model.Notification {
	return dao.NotificationDao.Get(id)
}

func (s *notificationService) Take(where ...interface{}) *model.Notification {
	return dao.NotificationDao.Take(where...)
}

func (s *notificationService) Find(cnd *querybuilder.QueryBuilder) []model.Notification {
	return dao.NotificationDao.Find(cnd)
}

func (s *notificationService) FindOne(cnd *querybuilder.QueryBuilder) *model.Notification {
	return dao.NotificationDao.FindOne(cnd)
}

func (s *notificationService) List(cnd *querybuilder.QueryBuilder) (list []model.Notification, paging *querybuilder.Paging) {
	return dao.NotificationDao.List(cnd)
}

func (s *notificationService) Create(t *model.Notification) error {
	return dao.NotificationDao.Create(t)
}

func (s *notificationService) Update(t *model.Notification) error {
	return dao.NotificationDao.Update(t)
}

func (s *notificationService) Updates(id int64, columns map[string]interface{}) error {
	return dao.NotificationDao.Updates(id, columns)
}

func (s *notificationService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.NotificationDao.UpdateColumn(id, name, value)
}

func (s *notificationService) Delete(id int64) {
	dao.NotificationDao.Delete(id)
}

// 获取未读消息数量
func (s *notificationService) GetUnReadCount(userId int64) (count int64) {
	return dao.NotificationDao.GetUnReadCount(userId)
}

// 将所有消息标记为已读
func (s *notificationService) MarkRead(userId int64) error {
	return dao.NotificationDao.UpdateStatusBatch(userId)
}

// 用户关注
func (s *notificationService) SendUserWatchNotification(userWatch *model.UserWatch) {
	user := cache.UserCache.Get(userWatch.WatcherID)

	var (
		fromId       = userWatch.WatcherID // 消息发送人
		authorId     int64                 // 被关注人
		content      string                // 消息内容
		quoteContent string                // 引用内容
	)

	authorId = userWatch.UserID
	content = user.Username.String + " 关注了你"
	quoteContent = ""

	if authorId <= 0 {
		return
	}
	// 给被关注者发消息
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeUserWatch, map[string]interface{}{
		"entityType":  model.EntityTypeUser,
		"entityId":    userWatch.WatcherID,
		"userWatchID": userWatch.ID,
	})
}

// 内容被点赞
func (s *notificationService) SendPostLikeNotification(postLike *model.PostLike) {
	user := cache.UserCache.Get(postLike.UserId)

	var (
		fromId       = postLike.UserId // 消息发送人
		authorId     int64              // 点赞者编号
		content      string             // 消息内容
		quoteContent string             // 引用内容
	)
	post := dao.PostDao.Get(postLike.PostId)
	if post != nil {
		authorId = post.UserId
		content = user.Username.String + " 点赞了你的话题：" + post.Title
		quoteContent = ""
	}

	if authorId <= 0 {
		return
	}
	// 给帖子作者发消息
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypePostLike, map[string]interface{}{
		"entityType":  model.EntityTypePost,
		"entityId":    post.ID,
		"postLikeId": postLike.ID,
	})
}

// 生产，将消息数据放入chan
func (s *notificationService) Produce(fromId, toId int64, content, quoteContent string, msgType int, extraDataMap map[string]interface{}) {
	to := cache.UserCache.Get(toId)
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

// 消费，消费chan中的消息
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

// 发送邮件通知
func (s *notificationService) SendEmailNotice(notification *model.Notification) {
	user := cache.UserCache.Get(notification.UserId)
	if user != nil && len(user.Email.String) > 0 {
		siteTitle := cache.SettingCache.GetValue(model.SettingSiteTitle)
		emailTitle := siteTitle + " 新消息提醒"

		email.SendTemplateEmail(user.Email.String, emailTitle, emailTitle, notification.Content,
			notification.QuoteContent, urls.AbsUrl("/user/notifications"))
		log.Info("发送邮件...email=%s", user.Email)
	} else {
		log.Info("邮件未发送，没设置邮箱...")
	}
}
