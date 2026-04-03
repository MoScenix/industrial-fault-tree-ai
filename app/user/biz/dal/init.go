package dal

import (
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
