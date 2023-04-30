package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"context"
	"memox_server/graph/generated"
	"memox_server/src/exception"
	"memox_server/src/service/memory"
	"memox_server/src/service/resource"
	"memox_server/src/service/storage/utils"
	"memox_server/src/service/user"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteResource is the resolver for the deleteResource field.
func (r *mutationResolver) DeleteResource(ctx context.Context, input string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(input)
	if err != nil {
		return false, exception.ErrInvalidID
	}
	res, err := r.resourceSvc.GetResource(ctx, id)
	if err != nil {
		return false, err
	}
	if len(res.Ref) > 0 {
		return false, exception.ErrResourceHasReference
	}
	err = r.resourceSvc.DeleteResource(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetToken is the resolver for the getToken field.
func (r *mutationResolver) GetToken(ctx context.Context, fileName string) (*utils.UploadTokenPayload, error) {
	ok, newFileName := utils.CheckFileName(fileName)
	if !ok {
		return nil, exception.ErrInvalidFileName
	}
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	// 检查容量限制
	u, err := r.userSvc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	used, err := r.User().Used(ctx, &u)
	if err != nil {
		return nil, err
	}
	subscribe, err := r.User().Subscribe(ctx, &u)
	if err != nil {
		return nil, err
	}
	if used > subscribe.Capacity {
		return nil, exception.ErrResourceSizeLimit
	}
	token, err := r.storageSvc.GetToken(ctx, newFileName)
	if err != nil {
		return nil, err
	}
	res, err := r.resourceSvc.NewResource(ctx, newFileName, 0)
	if err != nil {
		return nil, err
	}
	token.ID = res
	return token, nil
}

// LocalUpload is the resolver for the localUpload field.
func (r *mutationResolver) LocalUpload(ctx context.Context, input generated.LocalUploadInput) (string, error) {
	if strings.ToLower(r.conf.Storage.StorageProvider) == "local" {
		return r.resourceSvc.LocalUpload(ctx, input.SessionToken, input.Upload)
	}
	return "", exception.ErrInvalidStorageProvider
}

// AllResources is the resolver for the allResources field.
func (r *queryResolver) AllResources(ctx context.Context, page int64, size int64, byCreate bool, desc bool) ([]*resource.Resource, error) {
	userID, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	resources, err := r.resourceSvc.GetResources(ctx, userID, page, size, byCreate, desc)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// ID is the resolver for the id field.
func (r *resourceResolver) ID(ctx context.Context, obj *resource.Resource) (string, error) {
	return obj.ObjectID.Hex(), nil
}

// User is the resolver for the user field.
func (r *resourceResolver) User(ctx context.Context, obj *resource.Resource) (*user.User, error) {
	u, err := r.userSvc.GetUser(ctx, obj.Uid)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Memories is the resolver for the memories field.
func (r *resourceResolver) Memories(ctx context.Context, obj *resource.Resource) ([]*memory.Memory, error) {
	memoryIDs := obj.Ref
	var memories []*memory.Memory
	for _, id := range memoryIDs {
		m, err := r.memorySvc.GetMemory(ctx, id)
		if err != nil {
			return nil, err
		}
		memories = append(memories, m)
	}
	return memories, nil
}

// Resource returns generated.ResourceResolver implementation.
func (r *Resolver) Resource() generated.ResourceResolver { return &resourceResolver{r} }

type resourceResolver struct{ *Resolver }
