package mongo

import (
	"bytes"
	"context"
	"testing"

	"gitlab.com/sdce/exlib/exutil"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	pb "gitlab.com/sdce/protogo"
)

type TypicalModel struct {
	ID    *pb.UUID `protobuf:"bytes,1,opt,name=id" bson:"_id"`
	extra *pb.UUID
}

func init() {
	ctx := context.Background()
	conf := Config{URI: "mongodb://localhost:27017",
		DbName: "testing_db"}
	db = Connect(ctx, conf)
}

var (
	uuidTests = []struct {
		id *pb.UUID
	}{
		{exutil.NewUUID()},
	}

	db *Database
)

func TestUUID(t *testing.T) {
	ctx := context.Background()

	c := db.CreateCollection("test")

	for _, tt := range uuidTests {
		in := TypicalModel{ID: tt.id}
		res, err := c.InsertOne(ctx, in)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		insertedID := res.InsertedID.(primitive.ObjectID)
		if err != nil {
			t.Fatal(err.(mongo.MarshalError).Err)
		}
		cur, err := c.Find(ctx, bson.M{})
		var out []*TypicalModel
		err = DecodeCursorToSlice(ctx, cur, &out)
		c.DeleteMany(ctx, bson.M{})
		if err != nil {
			t.Errorf("Got error decoding: %v", err)
			t.FailNow()
		} else if len(out) != 1 {
			t.Errorf("Expecting 1 item. Got %v", len(out))
			t.FailNow()
		}

		t.Logf("Inserted %v", res.InsertedID)
		if err != nil {
			t.Errorf("Failed to decode model: %v", err)
			t.FailNow()
		}
		if !bytes.Equal(insertedID[:], out[0].ID.Bytes) {
			t.Errorf("Inserted ID %x doesn't match with read ID %x", res.InsertedID, out[0].ID.Bytes)
			t.FailNow()
		}
		if !bytes.Equal(tt.id.Bytes, out[0].ID.Bytes) {
			t.Errorf("Written %v Got %v.", tt.id, out[0].ID)
			t.FailNow()
		}

	}
}

func TestNilID(t *testing.T) {
	ctx := context.Background()

	c := db.CreateCollection("test")
	defer c.DeleteMany(ctx, bson.M{})

	res, err := c.InsertOne(ctx, &TypicalModel{})
	if err != nil {
		t.Errorf("Error inserting: %v", err)
		t.FailNow()
	}
	if res.InsertedID != nil {
		t.Errorf("Expecting nil insertedId. Got %v", res.InsertedID)
	}

	cur, err := c.Find(ctx, bson.M{})
	var out []*TypicalModel
	err = DecodeCursorToSlice(ctx, cur, &out)
	if err != nil {
		t.Errorf("Got error decoding: %v", err)
		t.FailNow()
	} else if len(out) != 1 {
		t.Errorf("Expecting 1 item. Got %v", len(out))
		t.FailNow()
	}

	if out[0].extra != nil || out[0].ID != nil {
		t.Errorf("Expecting nil ID and extra. Got %v", out[0])
	}
}
