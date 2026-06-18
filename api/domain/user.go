package domain

// User 用户领域模型
type User struct {
	ID           int64
	Username     string
	Email        string
	Nickname     string
	Avatar       string
	Password     string
	Website      string
	Description  string
	Status       int
	TopicCount   int64
	CommentCount int64
	Level        int
	CreateTime   int64
	UpdateTime   int64
}

// UserScore 用户积分
type UserScore struct {
	ID         int64
	UserId     int64
	Score      int
	CreateTime int64
	UpdateTime int64
}

// UserWatch 用户关注
type UserWatch struct {
	ID         int64
	UserID     int64
	WatcherID  int64
	CreateTime int64
}

// UserScoreLog 用户积分流水
type UserScoreLog struct {
	ID          int64
	UserId      int64
	SourceType  string
	SourceId    string
	Description string
	Type        int
	Score       int
	CreateTime  int64
}
