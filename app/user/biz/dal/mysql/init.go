package mysql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/conf"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	otelgorm "gorm.io/plugin/opentelemetry/tracing"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if err := DB.Use(otelgorm.NewPlugin(otelgorm.WithDBSystem("mysql"), otelgorm.WithoutMetrics())); err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	if err := ensureDefaultAdmin(); err != nil {
		panic(err)
	}
}

func ensureDefaultAdmin() error {
	ctx := context.Background()
	q := model.NewUserQuery(ctx, DB)
	_, err := q.GetUserByAccount("root")
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("rootroot"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = q.CreateUser(model.User{
		UserAccount:  "root",
		Name:         "root",
		UserRole:     "admin",
		PasswordHash: string(hashed),
	})
	return err
}
