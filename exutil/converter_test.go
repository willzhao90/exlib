package exutil

import (
	"fmt"
	"reflect"
	"testing"

	pbtypes "github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson"
)

func TestExutil_ApplyFieldMaskToBson(t *testing.T) {
	type args struct {
		src  interface{}
		mask *pbtypes.FieldMask
	}
	type spec struct {
		Foo string `bson:"foo"`
		Bar int    `bson:"bar"`
	}
	type testElem struct {
		name     string
		args     args
		wantResp bson.M
		wantErr  bool
	}
	em := &spec{
		Foo: "111",
		Bar: 222,
	}
	tests := []testElem{
		{
			name: "test struct with general member type",
			args: args{
				src: &struct {
					A string  `bson:"aaa"`
					B int     `bson:"bbb"`
					C float32 `bson:"ccc"`
					D byte
				}{
					A: "hello",
					B: 123,
					C: 456.9,
					D: 12,
				},
				mask: &pbtypes.FieldMask{
					Paths: []string{
						"A", "B", "C",
					},
				},
			},
			wantResp: bson.M{
				"aaa": "hello",
				"bbb": int(123),
				"ccc": float32(456.9),
			},
			wantErr: false,
		},
		{
			name: "test struct with embedded member type",
			args: args{
				src: &struct {
					A string  `bson:"aaa"`
					B int     `bson:"bbb"`
					C float32 `bson:"ccc"`
					D *spec   `bson:"ddd"`
				}{
					A: "hello",
					B: 123,
					C: 456.9,
					D: em,
				},
				mask: &pbtypes.FieldMask{
					Paths: []string{
						"A", "B", "C", "D",
					},
				},
			},
			wantResp: bson.M{
				"aaa": "hello",
				"bbb": int(123),
				"ccc": float32(456.9),
				"ddd": em,
			},
			wantErr: false,
		},
		{
			name: "test struct with partial embedded",
			args: args{
				src: &struct {
					A string  `bson:"aaa"`
					B int     `bson:"bbb"`
					C float32 `bson:"ccc"`
					D *spec   `bson:"ddd"`
				}{
					A: "hello",
					B: 123,
					C: 456.9,
					D: em,
				},
				mask: &pbtypes.FieldMask{
					Paths: []string{
						"A", "B", "C", "D.Bar",
					},
				},
			},
			wantResp: bson.M{
				"aaa":     "hello",
				"bbb":     int(123),
				"ccc":     float32(456.9),
				"ddd.bar": em.Bar,
			},
			wantErr: false,
		},
		{
			name: "miss a mask",
			args: args{
				src: &struct {
					A string  `bson:"aaa"`
					B int     `bson:"bbb"`
					C float32 `bson:"ccc"`
					D byte
				}{
					A: "hello",
					B: 123,
					C: 456.9,
					D: 12,
				},
				mask: &pbtypes.FieldMask{
					Paths: []string{
						"A", "B", "C", "E",
					},
				},
			},
			wantResp: nil,
			wantErr:  true,
		},
		{
			name: "miss a tag, using field name",
			args: args{
				src: &struct {
					A   string  `bson:"aaa"`
					B   int     `bson:"bbb"`
					C   float32 `bson:"ccc"`
					DDD byte
				}{
					A:   "hello",
					B:   123,
					C:   456.9,
					DDD: 12,
				},
				mask: &pbtypes.FieldMask{
					Paths: []string{
						"A", "B", "C", "DDD",
					},
				},
			},
			wantResp: bson.M{
				"aaa": "hello",
				"bbb": int(123),
				"ccc": float32(456.9),
				"ddd": byte(12),
			},
			wantErr: false,
		},
		{
			name: "pass a struct instead of a pointer",
			args: args{
				src: struct {
					A string  `bson:"aaa"`
					B int     `bson:"bbb"`
					C float32 `bson:"ccc"`
					D byte
				}{
					A: "hello",
					B: 123,
					C: 456.9,
				},
				mask: &pbtypes.FieldMask{
					Paths: []string{
						"A", "B", "C", "D",
					},
				},
			},
			wantResp: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objs, err := ApplyFieldMaskToBson(tt.args.src, tt.args.mask)
			if (err != nil) != tt.wantErr {
				fmt.Printf("field type %v %v", tt.args.src, tt.args.mask)
				t.Errorf("ApplyFieldMaskToBson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(objs, tt.wantResp) {
				t.Errorf("ApplyFieldMaskToBson() = %v, want %v", objs, tt.wantResp)
				return
			}
		})
	}
}

func TestGenerateFieldMask(t *testing.T) {
	type args struct {
		paths []string
		dest  interface{}
	}
	type Foo struct {
		Boo string
	}
	tests := []struct {
		name    string
		args    args
		want    *pbtypes.FieldMask
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				paths: []string{"A", "B", "C", "D.Boo", "E.Boo"},
				dest: &struct {
					A string
					B int
					C byte
					D Foo
					E *Foo
				}{},
			},
			want: &pbtypes.FieldMask{
				Paths: []string{"A", "B", "C", "D.Boo", "E.Boo"},
			},
			wantErr: false,
		}, {
			name: "neg1",
			args: args{
				paths: []string{"D.Boo.Bar"},
				dest: &struct {
					D Foo
				}{},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "neg2",
			args: args{
				paths: []string{"E.Boo.Bar"},
				dest: &struct {
					E *Foo
				}{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFieldMask(tt.args.paths, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFieldMask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateFieldMask() = %v, want %v", got, tt.want)
			}
		})
	}
}
