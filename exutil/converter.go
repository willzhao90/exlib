package exutil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"

	pbtypes "github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson"
)

func checkField(t reflect.Type, paths []string) (err error) {
	if len(paths) == 0 {
		return nil
	}
	field, ok := t.FieldByName(paths[0])
	if !ok {
		return fmt.Errorf("cannot find field %v", paths[0])
	}
	st := field.Type
	switch st.Kind() {
	case reflect.Ptr:
		st = field.Type.Elem()
		fallthrough
	case reflect.Struct:
		return checkField(st, paths[1:])
	default:
		if len(paths) != 1 {
			return fmt.Errorf("invalid sub field detected: %v is not a struct type", paths)
		}
		return nil
	}
}

func parseField(t reflect.Type, v reflect.Value, paths []string) (val interface{}, err error) {
	if len(paths) == 0 {
		err = fmt.Errorf("failed to parse out value for path %v", paths)
		return
	}

	field, ok := t.FieldByName(paths[0])
	if ok {
		val = v.FieldByName(paths[0]).Interface()
		return
	}
	st := field.Type
	if st.Kind() == reflect.Ptr {
		st = field.Type.Elem()
	}
	switch st.Kind() {
	case reflect.Struct:
		val, err = parseField(st, v, paths[1:])
	default:
		err = fmt.Errorf("failed to parse out value for path %v", paths)
	}
	return
}

// GenerateFieldMask new a fieldmask from comparing a fieldmask from frontend and a golang struct
func GenerateFieldMask(paths []string, dest interface{}) (*pbtypes.FieldMask, error) {
	// note: mask from frontend use json tags while pb fieldmask use struct name, not pb field name
	// todo: json name -> struct name
	fm := &pbtypes.FieldMask{Paths: paths}

	k := reflect.TypeOf(dest).Elem()
	fmt.Printf("field type %v", k)
	for _, path := range paths {
		// quick fix to let 1 level struct checking pass. todo, fix me properly.
		subFields := strings.Split(path, ".")
		err := checkField(k, subFields)
		if err != nil {
			return nil, fmt.Errorf("invalid sub field detected: %v", subFields)
		}
	}
	return fm, nil
}

//ApplyFieldMaskToBson using fieldmask to get the key/value of a struct
// and apply to a bson struct. The param "src" should be a struct pointer.
func ApplyFieldMaskToBson(src interface{}, mask *pbtypes.FieldMask) (objs bson.M, err error) {
	if mask == nil {
		return nil, errors.New("field mask is null")
	}
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = r.(error)
			default:
				err = errors.New("ConvertProtoToBson: Unknown panic")
			}
			objs = nil
		}
	}()

	var splittedMask [][]string

	for _, path := range mask.GetPaths() {
		splittedMask = append(splittedMask, strings.Split(path, "."))
	}

	return applyFieldMaskToBson(src, splittedMask)
}

func applyFieldMaskToBson(src interface{}, mask [][]string) (out bson.M, err error) {
	out = bson.M{}
	k := reflect.TypeOf(src).Elem()
	v := reflect.ValueOf(src).Elem()
	var subFields = make(map[string][][]string)

	for _, path := range mask {
		field, ok := k.FieldByName(path[0])
		if !ok {
			return nil, fmt.Errorf("Trying to access %v in type %v", path[0], k)
		}
		if len(path) == 1 {
			bsonName, err := getBsonName(field)
			if err != nil {
				return nil, err
			}
			out[bsonName] = v.FieldByIndex(field.Index).Interface()
		} else {
			subFields[path[0]] = append(subFields[path[0]], path[1:])
		}
	}

	for subField, subMask := range subFields {
		// Already handled in the previous loop
		field, _ := k.FieldByName(subField)
		bsonName, err := getBsonName(field)
		if err != nil {
			return nil, err
		}
		subOutput, err := applyFieldMaskToBson(
			v.FieldByIndex(field.Index).Interface(),
			subMask,
		)
		if err != nil {
			return nil, err
		}
		for subKey, subVal := range subOutput {
			out[bsonName+"."+subKey] = subVal
		}
	}
	return out, nil
}

func getBsonName(field reflect.StructField) (string, error) {
	key, err := bsoncodec.DefaultStructTagParser(field)
	if err != nil {
		return "", err
	}
	return key.Name, nil
}
