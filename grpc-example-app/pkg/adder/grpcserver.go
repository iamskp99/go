package adder

import (
	"context"
	"math/rand"
	"time"

	"github.com/keploy/go-sdk/integrations/kmongo"
	api "github.com/koddr/example-go-grpc-server/pkg/api/api/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GRPCServer struct
type GRPCServer struct {
	api.UnimplementedAdderServer
}

var col *kmongo.Collection

func New(host, db string) (*mongo.Client, error) {
	clientOptions := options.Client()

	clientOptions.ApplyURI("mongodb://" + host + "/" + db + "?retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, clientOptions)
}

// Add method for calculate X + Y
func (s *GRPCServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	dbName, collection := "keploy", "grpc-mongo"
	client, err := New("localhost:27017", dbName)
	db := client.Database(dbName)

	// integrate keploy with mongo
	col = kmongo.NewCollection(db.Collection(collection))
	if err != nil {
		panic(err)
	}
	// mongo, err := mgo.Dial("127.0.0.1:27017")
	// if err != nil {
	// 	log.Fatalf("Could not connect to the MongoDB server: %v", err)
	// }
	// defer mongo.Close()

	// DB = &mong{mongo.DB("mydb").C("mycol")}
	// DB.Operation.Insert(req)
	col.InsertOne(ctx, req)
	return &api.AddResponse{
		Result: rand.Int31n(200),
		Data:   &api.UserData{Name: "Fabio Di Gentanio", Team: &api.TeamData{Name: "Ducati", Championships: "0", Points: "1001"}},
	}, nil
}
