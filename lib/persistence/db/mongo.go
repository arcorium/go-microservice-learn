package db

import (
	"Microservices/lib/util"
	"Microservices/service/event/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoDBService struct {
	client *mongo.Client
}

func newMongoDBService(connection_ string) (*MongoDBService, error) {
	if client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connection_)); err == nil {
		return &MongoDBService{client: client}, nil
	} else {
		return &MongoDBService{client: nil}, err
	}
}

func (m *MongoDBService) GetDatabase(name_ string, options_ ...*options.DatabaseOptions) *mongo.Database {
	return m.client.Database(name_, options_...)
}

func (m *MongoDBService) GetCollection(dbName_ string, dbOptions_ *options.DatabaseOptions, colName_ string, colOptions_ ...*options.CollectionOptions) *mongo.Collection {
	return m.GetDatabase(dbName_, dbOptions_).Collection(colName_, colOptions_...)
}

func (m *MongoDBService) AddLocation(location *model.Location) (any, error) {
	return "mizhan", nil
}

func (m *MongoDBService) AddEvent(event *model.Event) (any, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Get database collection
	coll := m.GetCollection(DB_NAME, options.Database(), COLL_EVENT)
	if coll == nil {
		return primitive.ObjectID{}, errors.New("there are no such collections")
	}

	// Insert object
	event.Location.Id = primitive.NewObjectID()
	if res, err := coll.InsertOne(ctx, *event); err != nil {
		log.Println(err)
	} else {
		return res.InsertedID.(primitive.ObjectID), nil
	}

	return primitive.ObjectID{}, errors.New("cannot add event into database")
}

func (m *MongoDBService) FindEventById(id_ any) (*model.Event, error) {
	id, err := primitive.ObjectIDFromHex(id_.(string))
	if err != nil {
		return nil, err
	} else {
		return m.FindOneEvent(bson.M{"_id": id})
	}
}

func (m *MongoDBService) FindEventByName(name_ string) (*model.Event, error) {
	return m.FindOneEvent(bson.M{"name": name_})
}

func (m *MongoDBService) FindOneEvent(filter_ any) (*model.Event, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Creating session
	session := util.PackReturn(m.client.StartSession())
	defer session.EndSession(ctx)

	//util.LogError(session.StartTransaction())
	//defer util.LogError(session.AbortTransaction(ctx))

	log.Println("Session ID : ", session.ID().String())
	log.Println("Total Session :", m.client.NumberSessionsInProgress())

	// Get database collection
	coll := m.GetCollection(DB_NAME, options.Database(), COLL_EVENT)
	if coll == nil {
		//return model.Event{}, errors.New("there are no such collections")
		return nil, errors.New("there are no such collections")
	}

	// Find data
	if res := coll.FindOne(ctx, filter_); res != nil && res.Err() == nil {
		var event model.Event
		err := res.Decode(&event)
		//fmt.Printf("Address of event in FindOne : %p\n", &event)
		return &event, err
	} else {
		return nil, res.Err()
	}

}

func (m *MongoDBService) FindAllEvents() ([]model.Event, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Creating session
	session := util.PackReturn(m.client.StartSession())
	defer session.EndSession(ctx)

	//util.LogError(session.StartTransaction())
	//defer util.LogError(session.AbortTransaction(ctx))

	log.Println("Session ID : ", session.ID().String())
	log.Println("Total Session :", m.client.NumberSessionsInProgress())

	// Get database collection
	coll := m.GetCollection(DB_NAME, options.Database(), COLL_EVENT)
	if coll == nil {
		return []model.Event{}, errors.New("there are no such collections")
	}

	// Find data
	var events []model.Event
	if cursor, err := coll.Find(ctx, bson.M{}); err != nil {
		return events, err
	} else {
		err = cursor.All(ctx, &events)
		return events, err
	}
}

func (m *MongoDBService) FindAllAvailableEvents() ([]model.Event, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Creating session
	session := util.PackReturn(m.client.StartSession())
	defer session.EndSession(ctx)

	//util.LogError(session.StartTransaction())
	//defer util.LogError(session.AbortTransaction(ctx))

	log.Println("Session ID : ", session.ID().String())
	log.Println("Total Session :", m.client.NumberSessionsInProgress())

	// Get database collection
	coll := m.GetCollection(DB_NAME, options.Database(), COLL_EVENT)
	if coll == nil {
		return []model.Event{}, errors.New("there are no such collections")
	}

	// Find data
	var events []model.Event
	now := time.Now().Second()
	if cursor, err := coll.Find(ctx, bson.M{"startDate": bson.M{"$gt": now}}); err != nil {
		return events, err
	} else {
		err = cursor.All(ctx, events)
		return events, err
	}
}

func (m *MongoDBService) FindLocationById(any) (*model.Location, error) {
	return nil, nil
}

func (m *MongoDBService) FindLocationByName(string) (*model.Location, error) {
	return nil, nil
}

func (m *MongoDBService) FindAllLocations() ([]model.Location, error) {
	return nil, nil
}
