package service

import (
	"errors"
	"strconv"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

var UserScoreService = newUserScoreService()

func newUserScoreService() *userScoreService {
	return &userScoreService{}
}

type userScoreService struct {
}

func (s *userScoreService) Get(id int64) *model.UserScore {
	return dao.UserScoreDao.Get(id)
}

func (s *userScoreService) Take(where ...interface{}) *model.UserScore {
	return dao.UserScoreDao.Take(where...)
}

func (s *userScoreService) Find(cnd *querybuilder.QueryBuilder) []model.UserScore {
	return dao.UserScoreDao.Find(cnd)
}

func (s *userScoreService) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScore {
	return dao.UserScoreDao.FindOne(cnd)
}

func (s *userScoreService) List(cnd *querybuilder.QueryBuilder) (list []model.UserScore, paging *querybuilder.Paging) {
	return dao.UserScoreDao.List(cnd)
}

func (s *userScoreService) Create(t *model.UserScore) error {
	return dao.UserScoreDao.Create(t)
}

func (s *userScoreService) Update(t *model.UserScore) error {
	return dao.UserScoreDao.Update(t)
}

func (s *userScoreService) Updates(id int64, columns map[string]interface{}) error {
	return dao.UserScoreDao.Updates(id, columns)
}

func (s *userScoreService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.UserScoreDao.UpdateColumn(id, name, value)
}

func (s *userScoreService) Delete(id int64) {
	dao.UserScoreDao.Delete(id)
}

func (s *userScoreService) GetByUserId(userId int64) *model.UserScore {
	return s.FindOne(querybuilder.NewQueryBuilder().Eq("user_id", userId))
}

func (s *userScoreService) CreateOrUpdate(t *model.UserScore) error {
	if t.ID > 0 {
		return s.Update(t)
	} else {
		return s.Create(t)
	}
}

// IncrementCreatePostScore 发帖获积分
func (s *userScoreService) IncrementPostPostScore(post *model.Post) {
	config := SettingService.GetSetting()
	if config.ScoreConfig.PostPostScore <= 0 {
		log.Info("请配置发帖积分")
		return
	}
	err := s.addScore(post.UserId, config.ScoreConfig.PostPostScore, model.EntityTypePost,
		strconv.FormatInt(post.ID, 10), "发表话题")
	if err != nil {
		log.Error(err.Error())
	}
}

// Increment 增加分数
func (s *userScoreService) Increment(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

// Decrement 减少分数
func (s *userScoreService) Decrement(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

// addScore 加分数，也可以加负数
func (s *userScoreService) addScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score == 0 {
		return errors.New("分数不能为0")
	}
	userScore := s.GetByUserId(userId)
	if userScore == nil {
		userScore = &model.UserScore{
			UserId:     userId,
			CreateTime: util.NowTimestamp(),
		}
	}
	userScore.Score = userScore.Score + score
	userScore.UpdateTime = util.NowTimestamp()
	if err := s.CreateOrUpdate(userScore); err != nil {
		return err
	}

	scoreType := model.ScoreTypeIncr
	if score < 0 {
		scoreType = model.ScoreTypeDecr
	}
	err := UserScoreLogService.Create(&model.UserScoreLog{
		UserId:      userId,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreateTime:  util.NowTimestamp(),
	})
	if err == nil {
		cache.UserCache.InvalidateScore(userId)
	}
	return err
}
