package mongo

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	log "github.com/sirupsen/logrus"
	"gitlab.com/sdce/exlib/exutil"
	pb "gitlab.com/sdce/protogo"
)

var (
	testColl *mongo.Collection
)

func init() {
	ctx := context.Background()
	conf := Config{URI: "mongodb://localhost:27017",
		DbName: "testing_db"}
	db := Connect(ctx, conf)

	testColl = db.CreateCollection("testColl")
}

func SampleTradesCursor(ctx context.Context, testColl *mongo.Collection, n int) (cur *mongo.Cursor, err error) {
	err = testColl.Drop(ctx)
	if err != nil {
		return
	}
	for i := 1; i <= n; i++ {
		trade := &pb.TradeDefined{
			Id:     exutil.NewUUID(),
			Price:  3.14,
			Value:  "314",
			Volume: "100",
			Time:   time.Now().UnixNano()}
		_, err = testColl.InsertOne(ctx, trade)
		if err != nil {
			return
		}
	}

	return testColl.Find(ctx, bson.M{})
}

func TestCursorIteratorWrapper(t *testing.T) {
	ctx := context.Background()
	tradesCursor, err := SampleTradesCursor(ctx, testColl, 2)
	if err != nil {
		t.Log("Failed to get sample trades: ", err)
		t.FailNow()
	}
	counter := 0
	err = DecodeCursorToIterator(ctx, reflect.TypeOf((*pb.TradeDefined)(nil)), tradesCursor, func(message interface{}) {
		got := message.(*pb.TradeDefined)
		log.Info(got)
		counter++
	})
	if err != nil {
		log.Errorf("Error decoding cursor: %v", err)
		t.Fail()
	}
	if counter != 2 {
		log.Errorf("Expecting two iterator calls. Got %v", counter)
		t.Fail()
	}
}

func TestCursorSliceWrapper(t *testing.T) {

	ctx := context.Background()

	tradesCursor, err := SampleTradesCursor(ctx, testColl, 2)
	if err != nil {
		t.Log("Failed to get sample trades")
		t.FailNow()
	}
	out := make([]*pb.TradeDefined, 0, 2)
	err = DecodeCursorToSlice(ctx, tradesCursor, &out)
	if err != nil {
		log.Errorf("%v", err)
		t.FailNow()
	}
	log.Info(out)
	if len(out) != 2 {
		t.Fail()
	}

	tradesCursor, err = SampleTradesCursor(ctx, testColl, 2)
	if err != nil {
		t.Log("Failed to get sample trades")
		t.FailNow()
	}
	out2 := make([]pb.TradeDefined, 0, 2)
	err = DecodeCursorToSlice(ctx, tradesCursor, &out2)
	if err != nil {
		log.Errorf("%v", err)
		t.FailNow()
	}
	log.Info(out2)
	if len(out2) != 2 {
		t.Fail()
	}
}

func BenchmarkCursorWrapper(b *testing.B) {
	ctx := context.Background()

	tradesCursor, err := SampleTradesCursor(ctx, testColl, 100)
	if err != nil {
		b.Log("Failed to get sample trades")
		b.FailNow()
	}
	out := make([]*pb.TradeDefined, 0, 2)
	b.ResetTimer()
	DecodeCursorToSlice(ctx, tradesCursor, &out)
}

func TestUniqueFields(t *testing.T) {
	type args struct {
		ctx    context.Context
		c      *mongo.Collection
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Normal",
			args: args{
				ctx:    context.Background(),
				c:      testColl,
				fields: []string{"code"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UniqueFields(tt.args.ctx, tt.args.c, tt.args.fields); (err != nil) != tt.wantErr {
				t.Errorf("UniqueFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
