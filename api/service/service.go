package service

// Services 聚合所有服务实例，作为统一的服务访问入口
type Services struct {
	NodeService *nodeService
	PostService *postService
	UserService  *userService
}

// NewServices 集中初始化所有服务
// 当前阶段仍返回具体类型，后续重构 DI 时只需修改此函数的签名和内部实现
func NewServices() *Services {
	return &Services{
		NodeService: newNodeService(),
		PostService: newPostService(),
		UserService: newUserService(),
	}
}

// Srv 全局服务实例（过渡期使用）
// ⚠️ 注意：这是为了兼容现有代码的临时方案，最终目标是消除这个全局变量
var Srv = NewServices()