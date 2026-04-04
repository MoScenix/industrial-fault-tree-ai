package dal

import (
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/dal/redis"
)

func Init() {
	redis.Init()
	//mysql.Init()
}
