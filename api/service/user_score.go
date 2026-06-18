package service

import (
	"errors"
	"time"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util/querybuilder"
)

// UserScoreService 用户积分业务契约
type UserScoreService interface {
	Get(id int64) *model.UserScore
	Take(where ...interface{}) *model.UserScore
	Find(cnd *querybuilder.QueryBuilder) []model.UserScore
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserScore
	List(cnd *querybuilder.QueryBuilder) ([]model.UserScore, *querybuilder.Paging)
	Create(t *model.UserScore) error
	Update(t *model.UserScore) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetByUserId(userId int64) *model.UserScore
	CreateOrUpdate(t *model.UserScore) error
	Increment(userId int64, score int, sourceType, sourceId, description string) error
	Decrement(userId int64, score int, sourceType, sourceId, description string) error
}

func NewUserScoreService(repo repository.UserScoreRepository, scoreLogSvc UserScoreLogService, userCache cache.UserCacheInterface) UserScoreService {
	return &userScoreService{repo: repo, scoreLogSvc: scoreLogSvc, userCache: userCache}
}

type userScoreService struct {
	repo        repository.UserScoreRepository
	scoreLogSvc UserScoreLogService
	userCache   cache.UserCacheInterface
}

func (s *userScoreService) Get(id int64) *model.UserScore {
	return s.repo.Get(id)
}

func (s *userScoreService) Take(where ...interface{}) *model.UserScore {
	return s.repo.Take(where...)
}

func (s *userScoreService) Find(cnd *querybuilder.QueryBuilder) []model.UserScore {
	return s.repo.Find(cnd)
}

func (s *userScoreService) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScore {
	return s.repo.FindOne(cnd)
}

func (s *userScoreService) List(cnd *querybuilder.QueryBuilder) ([]model.UserScore, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *userScoreService) Create(t *model.UserScore) error {
	return s.repo.Create(t)
}

func (s *userScoreService) Update(t *model.UserScore) error {
	return s.repo.Update(t)
}

func (s *userScoreService) Updates(id int64, columns map[string]interface{}) error {
	return s.repo.Updates(id, columns)
}

func (s *userScoreService) UpdateColumn(id int64, name string, value interface{}) error {
	return s.repo.UpdateColumn(id, name, value)
}

func (s *userScoreService) Delete(id int64) {
	s.repo.Delete(id)
}

func (s *userScoreService) GetByUserId(userId int64) *model.UserScore {
	return s.FindOne(querybuilder.NewQueryBuilder().Eq("user_id", userId))
}

func (s *userScoreService) CreateOrUpdate(t *model.UserScore) error {
	if t.ID > 0 {
		return s.Update(t)
	}
	return s.Create(t)
}

func (s *userScoreService) Increment(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

func (s *userScoreService) Decrement(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

func (s *userScoreService) addScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score == 0 {
		return errors.New("分数不能为0")
	}
	userScore := s.GetByUserId(userId)
	if userScore == nil {
		userScore = &model.UserScore{
			UserId:    userId,
			CreatedAt: time.Now(),
		}
	}
	userScore.Score = userScore.Score + score
	userScore.UpdatedAt = time.Now()
	if err := s.CreateOrUpdate(userScore); err != nil {
		return err
	}

	scoreType := model.ScoreTypeIncr
	if score < 0 {
		scoreType = model.ScoreTypeDecr
	}
	err := s.scoreLogSvc.Create(&model.UserScoreLog{
		UserId:      userId,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreatedAt:   time.Now(),
	})
	if err == nil {
		s.userCache.InvalidateScore(userId)
	}
	return err
}
