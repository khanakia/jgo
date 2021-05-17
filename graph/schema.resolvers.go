package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/khanakia/jgo/graph/generated"
	"github.com/khanakia/jgo/graph/model"
	"github.com/khanakia/jgo/pkg/auth"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AuthRegister(ctx context.Context, args model.RegisterInput) (*model.User, error) {
	// graphql.AddError(ctx, gqlerror.Errorf("zzzzzt"))
	// return nil, gqlerror.Errorf("BOOM! Headshot")

	user := auth.FindByEmail("khanakia@gmail.com", r.Auth.Dbc.Db)

	// errForm := user.ValidateRegister()
	// fmt.Println(errForm)

	// // var user2 User11
	user2 := model.User{}
	// user2.FullName = "ds"
	copier.Copy(&user2, user)
	// // user2.Email = "test@gmail.com"
	// // user2.ID = "1"
	// fmt.Println(user2)

	// util.OzzoErrToGraphqlErrors(errForm, ctx)

	// i := 0
	// for field, v := range errForm.(validation.Errors) {

	// 	message := strings.Title(field) + " " + v.Error()
	// 	graphql.AddError(ctx, &gqlerror.Error{
	// 		Path:    graphql.GetPath(ctx),
	// 		Message: message,
	// 		Extensions: map[string]interface{}{
	// 			"idx":   i,
	// 			"field": field,
	// 		},
	// 	})
	// 	i++
	// }

	return &user2, nil

	// var user auth.User

	// user.Email = args.Email

	// usern, err := r.Auth.Register(&user)

	// if err != nil {
	// 	return nil, errors.New("Cannot register")
	// }

	// var user2 model.User
	// copier.Copy(&user2, &usern)

	// return &user2, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AuthLogin(ctx context.Context, args model.LoginInput) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
type User11 struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	ID       string `json:"id"`
}
