package mock

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"

	"ultrathreads/model"
	"ultrathreads/dao"
	"ultrathreads/util"
)

var (
	// 头像假数据
	avatars = []string{
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
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
		Username:   util.SqlNullString(username),
		Password:   util.EncodePassword(password),
		Status:     model.StatusOk,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}

	r := RandInt(0, len(avatars)-1)

	return factory.NewFactory(
		u,
	).Attr("Nickname", func(args factory.Args) (interface{}, error) {
		// ✅ admin 用户的昵称也设为 admin，保持一致性
		if i == 0 {
			return "admin", nil
		}
		return fmt.Sprintf("user-%d", i+1), nil
	}).Attr("Avatar", func(args factory.Args) (interface{}, error) {
		return avatars[r], nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		paragraph := randomdata.Paragraph()
		if len(paragraph) >= 70 {
			paragraph = paragraph[:70]
		}
		return paragraph, nil
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
		if err := dao.UserDao.Create(user); err != nil {
			fmt.Printf("mock user error： %v\n", err)
		}
	}
}