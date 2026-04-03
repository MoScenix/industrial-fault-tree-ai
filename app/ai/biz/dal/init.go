package dal

import (
	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
