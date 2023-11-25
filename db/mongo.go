package db

import (
	"context"
	"log"

	"github.com/Hunnnn77/hello/service"
	"github.com/Hunnnn77/hello/util"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collections service.Colls

func generate(colUser, colActress *mongo.Collection) service.Colls {
	return service.Colls{
		UserCollection:    colUser,
		ActressCollection: colActress,
	}
}

func Initialize() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := util.ByField("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database(util.ByField("DB"))
	Collections = generate(db.Collection(util.ByField("COL_USER")), db.Collection(util.ByField("COL_ACTRESS")))
}
