package resolver

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"time_speak_server/src/config"
	"time_speak_server/src/service/hashtag"
	"time_speak_server/src/service/history"
	"time_speak_server/src/service/mail"
	"time_speak_server/src/service/memory"
	"time_speak_server/src/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	conf       config.Config
	db         *mongo.Database
	r          *redis.Client
	userSvc    *user.Svc
	mailSvc    *mail.Svc
	memorySvc  *memory.Svc
	hashtagSvc *hashtag.Svc
	historySvc *history.Svc
}

func NewResolver(conf config.Config, db *mongo.Database, r *redis.Client) *Resolver {
	return &Resolver{
		conf:       conf,
		db:         db,
		r:          r,
		userSvc:    user.NewUserSvc(conf.User, db, r),
		mailSvc:    mail.NewMailSvc(conf.Mail, r),
		memorySvc:  memory.NewMemorySvc(conf.Memory, db, r),
		hashtagSvc: hashtag.NewHashtagSvc(conf.Hashtag, db, r),
		historySvc: history.NewHistorySvc(conf.History, db, r),
	}
}
