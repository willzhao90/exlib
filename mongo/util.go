package mongo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	pb "gitlab.com/sdce/protogo"
)

// IDFilter returns a filter based on mongo's `_id` field
func IDFilter(uuid *pb.UUID) bson.M {
	return bson.M{"_id": uuid}
}

// SearchFilter returns filter that search collection.
// Requires text index
func SearchFilter(searchInput string) bson.M {
	if searchInput == "" {
		return bson.M{}
	}
	return bson.M{"$text": bson.M{"$search": searchInput}}
}

// DecodeCursorToIterator sends items in mongo.Cursor one by one
func DecodeCursorToIterator(ctx context.Context, rt reflect.Type, cur *mongo.Cursor, send func(message interface{})) error {
	retPtr := false
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		retPtr = true
	}
	for cur.Next(ctx) {
		elem := reflect.New(rt)
		if err := cur.Decode(elem.Interface()); err != nil {
			return fmt.Errorf("cannot decode cursor to iterator: %v", err)
		}
		if !retPtr {
			elem = elem.Elem()
		}
		send(elem.Interface())
	}
	return nil
}

// DecodeCursorToSlice gets items in mongo.Cursor to slice
// out should be *[]*TYPE
func DecodeCursorToSlice(ctx context.Context, cur *mongo.Cursor, out interface{}) error {

	rt := reflect.TypeOf(out)
	rv := reflect.ValueOf(out)

	// out must be a ptr to pass value out
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	} else {
		return fmt.Errorf("Input param is not a pointer to slice. Got %v", rt.Kind())
	}

	if rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	} else {
		return fmt.Errorf("Input param is not a slice. got %v", rt.Kind())
	}

	isElemPtr := false

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		isElemPtr = true
	}

	for cur.Next(ctx) {
		newItem := reflect.New(rt)
		if err := cur.Decode(newItem.Interface()); err != nil {
			return fmt.Errorf("cannot decode cursor to slice: %v", err)
		}
		if !isElemPtr {
			newItem = newItem.Elem()
		}
		rv.Set(reflect.Append(rv, newItem))
	}
	return nil
}

// UniqueFields make a list of fields unique in a collection
func UniqueFields(ctx context.Context, c *mongo.Collection, fields []string) error {
	for _, key := range fields {
		_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bsonx.Doc{{key, bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true)})
		if err != nil {
			return err
		}
	}
	return nil
}

//NewPaginationOptions
func NewPaginationOptions(pageIndex int64, pageSize int64) *options.FindOptions {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	opts.SetSkip(pageSize * pageIndex)
	return opts
}
