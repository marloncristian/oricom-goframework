package database

import (
	"context"
	"reflect"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// ServiceBase base service struct
type ServiceBase struct {
	collectionName string
}

// fill parses and fill the collection documents
func (base ServiceBase) fill(slice interface{}, cursor mongo.Cursor) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter slice must be a pointer" }
	}

	for cursor.Next(context.Background()) {

		spt := reflect.ValueOf(slice)
		svl := spt.Elem()

		sl := reflect.Indirect(spt)
		tT := sl.Type().Elem()

		ptr := reflect.New(tT).Interface()

		err := cursor.Decode(ptr)
		if err != nil {
			return err
		}

		s := reflect.ValueOf(ptr).Elem()

		svl.Set(reflect.Append(svl, s))
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

// query retrieves documents by query or all
func (base ServiceBase) query(query interface{}, slice interface{}) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter slice must be a pointer" }
	}

	col := database.Collection(base.collectionName)
	cur, err := col.Find(nil, query)
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())
	if err := base.fill(slice, cur); err != nil {
		return err
	}

	return nil
}

// queryAndPage retrieves an specific page of a document query
func (base ServiceBase) queryAndPage(query interface{}, slice interface{}, skip int64, limit int64) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter slice must be a pointer" }
	}

	opt := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}

	col := database.Collection(base.collectionName)
	cur, err := col.Find(nil, query, &opt)
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())
	if err := base.fill(slice, cur); err != nil {
		return err
	}

	return nil
}

// GetOne : returns a single instance of an object
func (base ServiceBase) GetOne(query interface{}, res interface{}) error {
	if reflect.ValueOf(res).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter res must be a pointer" }
	}

	opts := options.Find()
	opts.SetLimit(1)

	col := database.Collection(base.collectionName)
	cur, err := col.Find(context.Background(), query, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())

	if hasNext := cur.Next(context.Background()); !hasNext {
		return EntityNotFound{}
	}

	if err := cur.Decode(res); err != nil {
		return err
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

// GetByHexID returns an especific element its hexa string representation
func (base ServiceBase) GetByHexID(hexID string, res interface{}) error {
	if reflect.ValueOf(res).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter res must be a pointer" }
	}
	objID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return err
	}
	if err := base.GetOne(bson.M{"_id": objID}, res); err != nil {
		return err
	}
	return nil
}

// GetByObjID returns an especific element its objectiid
func (base ServiceBase) GetByObjID(objID primitive.ObjectID, res interface{}) error {
	if reflect.ValueOf(res).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter res must be a pointer" }
	}
	if err := base.GetOne(bson.M{"_id": objID}, res); err != nil {
		return err
	}
	return nil
}

// GetAll : returns all documents from collection
func (base ServiceBase) GetAll(slice interface{}) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return InvalidArgument{ Description : "parameter slice must be a pointer" }
	}
	if err := base.query(bson.D{{}}, slice); err != nil {
		return err
	}
	return nil
}

// GetAllWithSkipLimit retrieves chunks of information defined by parameters skip and limit
func (base ServiceBase) GetAllWithSkipLimit(slice interface{}, skip int64, limit int64) error {
	return base.queryAndPage(bson.D{{}}, slice, skip, limit)
}

// GetWithSkipLimit returns a filtered and paged document list from repository
func (base ServiceBase) GetWithSkipLimit(query interface{}, slice interface{}, skip int64, limit int64) error {
	return base.queryAndPage(query, slice, skip, limit)
}

// CountAll returns a count of all documents in repository
func (base ServiceBase) CountAll() (int64, error) {
	col := database.Collection(base.collectionName)
	cnt, err := col.Count(nil, bson.D{{}})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// CountWithFilter returns a count of filtered documents
func (base ServiceBase) CountWithFilter(query interface{}) (int64, error) {
	col := database.Collection(base.collectionName)
	cnt, err := col.Count(context.Background(), query)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

//InsertOne : inserts a new object in repository
func (base ServiceBase) InsertOne(value interface{}) (primitive.ObjectID, error) {
	c := database.Collection(base.collectionName)
	res, err := c.InsertOne(context.Background(), value)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

// UpdateOne : updates an document
func (base ServiceBase) UpdateOne(id primitive.ObjectID, values map[string]interface{}, result interface{}) error {
	col := database.Collection(base.collectionName)
	doc := col.FindOneAndUpdate(context.Background(), bson.M{"_id": id}, bson.M{"$set": values})
	if result == nil {
		return nil
	}
	if err := doc.Decode(result); err != nil {
		return err
	}
	return nil
}

// DeleteOne removes an elemento from database
func (base ServiceBase) DeleteOne(id primitive.ObjectID) error {
	col := database.Collection(base.collectionName)
	_, err := col.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// NewServiceBase creates a new service base
func NewServiceBase(collectionName string) ServiceBase {
	return ServiceBase{
		collectionName: collectionName,
	}
}
