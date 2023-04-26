package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"context"
	"time"
	"time_speak_server/graph/generated"
	"time_speak_server/src/exception"
	"time_speak_server/src/service/mail"
	"time_speak_server/src/service/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input generated.LoginInput) (*generated.LoginPayload, error) {
	match, err := r.userSvc.CheckPasswordByUsername(ctx, input.Username, input.Password)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if !match {
		return nil, exception.ErrUsernameOrPasswordWrong
	}

	jwt, token, err := r.userSvc.GetTokenByUsername(ctx, input.Username)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	return &generated.LoginPayload{
		ID:         jwt.ID,
		Token:      token,
		Expire:     jwt.ExpiresAt,
		Permission: jwt.Permission,
	}, nil
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input generated.RegisterInput) (string, error) {
	if r.mailSvc.CheckEmailVerifyCode(ctx, input.Email, input.EmailVerifyCode) {
		_, err := r.userSvc.GetUserByMail(ctx, input.Email)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// 未注册
				id, err := r.userSvc.AddUser(ctx, input.Username, input.Password, input.Email)
				if err != nil {
					return "", err
				}
				return id.Hex(), nil
			} else {
				return "", err
			}
		}
		return "", exception.ErrMailOccupied
	}
	return "", exception.ErrVerifyCodeWrong
}

// Forget is the resolver for the forget field.
func (r *mutationResolver) Forget(ctx context.Context, input generated.ForgetInput) (bool, error) {
	if r.mailSvc.CheckEmailVerifyCode(ctx, input.Email, input.EmailVerifyCode) {
		_, err := r.userSvc.GetUserByMail(ctx, input.Email)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// 未注册
				return false, exception.ErrUserNotFound
			} else {
				return false, err
			}
		}
		u, err := r.userSvc.GetUserByMail(ctx, input.Email)
		if err != nil {
			return false, err
		}
		p, err := user.EncryptPassword(input.Password)
		if err != nil {
			return false, err
		}

		err = r.userSvc.UpdateUser(ctx, u.ObjectID, func(m bson.M) bson.M {
			m["password"] = p
			m["password_change_time"] = time.Now().Unix()
			return m
		})
		if err != nil {
			return false, err
		}

		_ = r.r.Del(ctx, mail.CodeKey(input.Email)).Err()
		return true, nil
	}
	return false, exception.ErrVerifyCodeWrong
}

// SendEmailCode is the resolver for the sendEmailCode field.
func (r *mutationResolver) SendEmailCode(ctx context.Context, input generated.SendEmailCodeInput) (bool, error) {
	act := "找回"
	if input.Register {
		act = "注册"
	}
	err := r.mailSvc.NewEmailVerifyCode(ctx, input.Mail, act)
	if err != nil {
		if err == mail.ErrVerifyCodeCoolDown {
			return false, exception.ErrTooManyRequest
		} else {
			return false, exception.InternalError(err)
		}
	}
	return true, nil
}
