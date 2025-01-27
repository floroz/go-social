package main

import (
	"context"
	"math/rand/v2"

	"github.com/bxcodec/faker/v3"
	"github.com/floroz/go-social/cmd/api"
	"github.com/floroz/go-social/cmd/database"
	"github.com/floroz/go-social/internal/domain"
	"github.com/floroz/go-social/internal/env"
	"github.com/floroz/go-social/internal/repositories"
	"github.com/floroz/go-social/internal/services"
	"github.com/rs/zerolog/log"
)

func seed(app *api.Application) {
	ctx := context.Background()
	users := seedUsers(ctx, app)
	posts := seedPosts(ctx, app, users)
	seedComments(ctx, app, posts)
}

func seedUsers(ctx context.Context, app *api.Application) []domain.User {
	const maxUsers = 50
	users := make([]domain.User, 0, maxUsers)
	for i := 0; i < maxUsers; i++ {
		user := domain.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Email:     faker.Email(),
			Username:  faker.Username(),
			Password:  "password",
		}

		createdUser, err := app.UserService.Create(ctx, &domain.CreateUserDTO{
			EditableUserField: domain.EditableUserField{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Username:  user.Username,
			},
			Password: user.Password,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to create user")
		}
		users = append(users, *createdUser)
	}
	return users
}

func seedPosts(ctx context.Context, app *api.Application, users []domain.User) []domain.Post {
	posts := make([]domain.Post, 0, len(users))
	for _, user := range users {
		numPosts := rand.IntN(10) + 1
		for i := 0; i < numPosts; i++ {
			dto := domain.CreatePostDTO{
				EditablePostFields: domain.EditablePostFields{
					Content: faker.Sentence(),
				},
			}

			post, err := app.PostService.Create(ctx, user.ID, &dto)

			if err != nil {
				log.Error().Err(err).Msg("failed to create post")
			}
			posts = append(posts, *post)
		}
	}
	return posts

}

func seedComments(ctx context.Context, app *api.Application, posts []domain.Post) []domain.Comment {
	comments := make([]domain.Comment, 0, len(posts))
	for _, post := range posts {
		numComments := rand.IntN(5)
		for i := 0; i < numComments; i++ {
			dto := domain.CreateCommentDTO{
				EditableCommentFields: domain.EditableCommentFields{
					Content: faker.Sentence(),
				},
			}

			comment, err := app.CommentService.Create(ctx, post.UserID, post.ID, &dto)
			if err != nil {
				log.Error().Err(err).Msg("failed to create comment")
			}
			comments = append(comments, *comment)
		}
	}

	return comments
}

func main() {
	env.MustLoadEnv(".env.local")
	db, err := database.ConnectDb()
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}
	defer db.Close()

	config := &api.Config{
		Port: env.GetEnvValue("PORT"),
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)
	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo, commentRepo)
	authService := services.NewAuthService(userRepo)

	app := &api.Application{
		Config:         config,
		UserService:    userService,
		PostService:    postService,
		CommentService: commentService,
		AuthService:    authService,
	}

	seed(app)
}
