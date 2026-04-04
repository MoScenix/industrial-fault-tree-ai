package middleware

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(r *server.Hertz) {
	r.Use(GlobalAuth())
	r.Use(Auth())
}
