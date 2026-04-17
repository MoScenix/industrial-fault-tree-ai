package mysql

import (
	"fmt"
	"os"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/conf"
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

	if err := DB.AutoMigrate(&model.Graph{}, &model.Message{}); err != nil {
		panic(err)
	}
}
