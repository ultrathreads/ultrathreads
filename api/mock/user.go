package mock

import (
	"fmt"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"

	"ultrathreads/model"
	"ultrathreads/util"
)

var (
	// 头像假数据
	avatars = []string{
		"https://www.dismall.com/uc_server/data/avatar/000/05/04/86_avatar_middle.jpg",
		"https://www.dismall.com/uc_server/data/avatar/000/04/32/97_avatar_middle.jpg",
		"https://www.dismall.com/uc_server/data/avatar/000/00/04/57_avatar_middle.jpg",
		"https://www.dismall.com/uc_server/data/avatar/000/03/48/31_avatar_middle.jpg",
		"https://www.dismall.com/uc_server/data/avatar/000/05/33/40_avatar_middle.jpg",
		"https://www.dismall.com/uc_server/data/avatar/000/04/40/78_avatar_middle.jpg",
	}

	websites = []string{
		"https://www.ultrathreads.com",
		"https://www.baidu.com",
		"https://www.toutiao.com",
		"https://www.qq.com",
		"https://www.163.com",
		"https://www.taobao.com",
	}
)

func userFactory(i int) *factory.Factory {

	password := "123456"
	username := ""

	// ✅ 第一个用户特殊处理：用户名 admin，密码 admin12345
	if i == 0 {
		username = "admin"
		password = "admin12345"
	} else if i == 1 {
		username = "ultrathreads"
		password = "ultrathreads"
	}

	u := &model.User{
		Username:  util.SqlNullString(username),
		Password:  util.EncodePassword(password),
		Status:    model.StatusOk,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r := RandInt(0, len(avatars)-1)
	w := RandInt(0, len(websites)-1)

	return factory.NewFactory(
		u,
	).Attr("Nickname", func(args factory.Args) (interface{}, error) {
		// ✅ admin 用户的昵称也设为 admin，保持一致性
		if i == 0 {
			return "admin", nil
		}
		return fmt.Sprintf("user-%d", i+1), nil
	}).Attr("Avatar", func(args factory.Args) (interface{}, error) {
		if i == 0 {
			return "", nil
		}
		return avatars[r], nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		paragraph := randomdata.Paragraph()
		if len(paragraph) >= 70 {
			paragraph = paragraph[:70]
		}
		return paragraph, nil
	}).Attr("Website", func(args factory.Args) (interface{}, error) {
		return websites[w], nil
	}).Attr("Email", func(args factory.Args) (interface{}, error) {
		if i == 0 {
			return util.SqlNullString("admin@test.com"), nil // ✅ admin 专用邮箱
		}
		if i == 1 {
			return util.SqlNullString("2@test.com"), nil
		}
		return util.SqlNullString(randomdata.Email()), nil
	}).Attr("Level", func(args factory.Args) (interface{}, error) {
		if i == 0 || i == 1 {
			return 10, nil
		}
		return 0, nil
	})
}

// UserTableSeeder -
func UserTableSeeder(needCleanTable bool, totalUsers int) {
	if needCleanTable {
		dropAndCreateTable(&model.User{})
	}

	for i := 0; i < totalUsers; i++ {
		user := userFactory(i).MustCreate().(*model.User)
		fmt.Println("Email:", user.Email)
		if err := userDao.Create(user); err != nil {
			fmt.Printf("mock user error： %v\n", err)
		}
	}
}
