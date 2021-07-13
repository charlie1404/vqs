package message

import "reflect"

// type DataTypeType string

// const (
// 	IntType        DataTypeType = "int"
// 	StringType     DataTypeType = "string"
// 	ListStringType DataTypeType = "array_string"
// 	ListIntType    DataTypeType = "array_int"
// )

type MessageAttribute struct {
	DataType  reflect.Type
	Value     string
	ListValue []interface{}
}

type MessageAttributes map[string]MessageAttribute
