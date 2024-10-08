package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"database/sql"

	"github.com/biboyqg/social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe",
}

var titles = []string{
	"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
	"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
	"Home Office Setup", "Digital Detox", "Gardening Basics",
	"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
	"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
	"Fitness at Home", "Personal Finance Tips", "Creative Writing",
	"Mental Health Awareness", "Learning New Skills",
}

var contents = []string{
	"In this post, we'll explore how to develop good habits that stick and transform your life.",
	"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind.",
	"Learn practical tips for eating healthy on a budget without sacrificing flavor.",
	"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget.",
	"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started.",
	"Increase your productivity with these simple and effective strategies.",
	"Set up the perfect home office to boost your work-from-home efficiency and comfort.",
	"A digital detox can help you reconnect with the real world and improve your mental health.",
	"Start your gardening journey with these basic tips for beginners.",
	"Transform your home with these fun and easy DIY projects.",
	"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try.",
	"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices.",
	"Master time management with these tips and get more done in less time.",
	"Nature has so much to offer. Discover the benefits of spending time outdoors.",
	"Whip up delicious meals with these simple and quick cooking recipes.",
	"Stay fit without leaving home with these effective at-home workout routines.",
	"Take control of your finances with these practical personal finance tips.",
	"Unleash your creativity with these inspiring writing prompts and exercises.",
	"Mental health is just as important as physical health. Learn how to take care of your mind.",
	"Learning new skills can be fun and rewarding. Here are some ideas to get you started.",
}

var tags = []string{
	"Self Improvement", "Minimalism", "Health", "Travel", "Mindfulness",
	"Productivity", "Home Office", "Digital Detox", "Gardening", "DIY",
	"Yoga", "Sustainability", "Time Management", "Nature", "Cooking",
	"Fitness", "Personal Finance", "Writing", "Mental Health", "Learning",
}

var comments = []string{
	"Great post! Thanks for sharing.",
	"I completely agree with your thoughts.",
	"Thanks for the tips, very helpful.",
	"Interesting perspective, I hadn't considered that.",
	"Thanks for sharing your experience.",
	"Well written, I enjoyed reading this.",
	"This is very insightful, thanks for posting.",
	"Great advice, I'll definitely try that.",
	"I love this, very inspirational.",
	"Thanks for the information, very useful.",
}


func Seed(store store.Storage, db *sql.DB) error {
	ctx := context.Background()

	users := generateUsers(100)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("error creating user", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("error creating post", err)
			return err
		}
	}

	comments := generateComments(500, posts, users)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("error creating comment", err)
			return err
		}
	}

	log.Println("successfully seeded", len(users), "users", len(posts), "posts", len(comments), "comments")

	return nil
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)

	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			RoleID: 1,
		}
	}

	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)

	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID: user.ID,
			Title: titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(n int, posts []*store.Post, users []*store.User) []*store.Comment {
	cms := make([]*store.Comment, n)

	for i := 0; i < n; i++ {
		cms[i] = &store.Comment{
			PostID: posts[rand.Intn(len(posts))].ID,
			UserID: users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
