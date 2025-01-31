package query

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/VinGitonga/gin-auth/modules/config"
	"github.com/VinGitonga/gin-auth/modules/database"
	"github.com/VinGitonga/gin-auth/modules/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoAppDB struct {
	App *config.GoAppTools
	DB  *mongo.Client
}

// InsertUser implements database.DBRepo.
func (g *GoAppDB) InsertUser(user *model.User) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	regMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	ok := regMail.MatchString(user.Email)

	if !ok {
		g.App.ErrorLogger.Println("invalid registred details")
		return false, 0, errors.New("invalid registered details")
	}

	filters := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M

	err := User(*g.DB, "users").FindOne(ctx, filters).Decode(&res)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			_, insertErr := User(*g.DB, "users").InsertOne(ctx, user)

			if insertErr != nil {
				g.App.ErrorLogger.Fatalf("cannot add user to the database: %v", insertErr)
			}

			return true, 1, nil
		}
		g.App.ErrorLogger.Fatal(err)
	}

	return true, 2, nil
}

// UpdateInfo implements database.DBRepo.
func (g *GoAppDB) UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	filter := bson.D{{Key: "_id", Value: userID}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: tk["t1"]}, {Key: "new_token", Value: tk["t2"]}}}}

	_, err := User(*g.DB, "users").UpdateOne(ctx, filter, update)

	if err != nil {
		return false, err
	}

	return true, nil
}

// VerifyUser implements database.DBRepo.
func (g *GoAppDB) VerifyUser(email string) (primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var res bson.M

	filter := bson.D{{Key: "email", Value: email}}

	err := User(*g.DB, "users").FindOne(ctx, filter).Decode(&res)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			g.App.ErrorLogger.Println("no documents found for this query")
			return nil, err
		}

		g.App.ErrorLogger.Fatalf("cannot execute the database query perfectly: %v", err)
	}

	return res, nil
}

func NewGoAppDB(app *config.GoAppTools, db *mongo.Client) database.DBRepo {
	return &GoAppDB{
		App: app,
		DB:  db,
	}
}
