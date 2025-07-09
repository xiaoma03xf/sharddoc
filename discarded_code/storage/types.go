package storage

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type InsertRes struct {
	TableName string
	Rec       *Record
	Err       error
}
type UpdateRes struct {
	TableName string
	UpdateMp  map[string]string
	Scan      *Scanner
	Err       error
}
type QueryResult struct {
	Recs []Record
	Err  error
}
type DelRes struct {
	Scan      *Scanner
	TableName string
}

func (qr *QueryResult) Scan(dest any) error {
	if qr.Err != nil {
		return qr.Err
	}
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		// 必须是指针, 因为必须要往里面写入数据
		return fmt.Errorf("dest must be a pointer")
	}
	destElem := destVal.Elem()
	switch destElem.Kind() {
	// 判断传入的是切片还是结构体
	case reflect.Slice:
		elemType := destElem.Type().Elem()
		slice := reflect.MakeSlice(destElem.Type(), 0, len(qr.Recs))

		for _, rec := range qr.Recs {
			elemPtr := reflect.New(elemType).Interface()
			fields := mapStructFieldsOnce(elemPtr)
			if err := scanRecordToStruct(rec, elemPtr, fields); err != nil {
				return err
			}
			slice = reflect.Append(slice, reflect.ValueOf(elemPtr).Elem())
		}
		destElem.Set(slice)
	case reflect.Struct:
		if len(qr.Recs) == 0 {
			return nil // 空数据
		}
		fields := mapStructFieldsOnce(dest)
		return scanRecordToStruct(qr.Recs[0], dest, fields)
	default:
		return fmt.Errorf("unsupported dest type: %v", destElem.Kind())
	}
	return nil
}

func mapStructFieldsOnce(ptr any) map[string]reflect.Value {
	val := reflect.ValueOf(ptr).Elem() // 获取结构体值
	typ := val.Type()

	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // 非导出字段跳过
			continue
		}
		fieldMap[strings.ToLower(field.Name)] = val.Field(i)
	}
	return fieldMap
}
func convertValue(v Value) (any, error) {
	switch v.Type {
	case TYPE_INT64:
		// 返回int64而不是直接转为int，让后续的类型转换更灵活
		return v.I64, nil
	case TYPE_BYTES:
		return string(v.Str), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %d", v.Type)
	}
}
func scanRecordToStruct(record Record, ptr any, fields map[string]reflect.Value) error {
	for i, col := range record.Cols {
		f, ok := fields[strings.ToLower(col)]
		if !ok {
			continue
		}
		val, err := convertValue(record.Vals[i])
		if err != nil {
			return err
		}
		rv := reflect.ValueOf(val)

		// 尝试智能类型转换
		switch {
		case rv.Type().AssignableTo(f.Type()):
			// 类型完全匹配
			f.Set(rv)
		case rv.Type().ConvertibleTo(f.Type()):
			// 类型可以转换（如int64->int, string->int等）
			f.Set(rv.Convert(f.Type()))
		case rv.Type().Kind() == reflect.Int64 && f.Type().Kind() == reflect.String:
			// 特殊处理：int64转string
			f.SetString(fmt.Sprintf("%d", rv.Int()))
		case rv.Type().Kind() == reflect.String && f.Type().Kind() == reflect.Int64:
			// 特殊处理：string转int64
			if i, err := strconv.ParseInt(rv.String(), 10, 64); err == nil {
				f.SetInt(i)
			} else {
				return fmt.Errorf("cannot convert string %q to int64", rv.String())
			}
		default:
			return fmt.Errorf("cannot convert %v (%v) to field %v (%v)",
				rv.Interface(), rv.Type(), f.Interface(), f.Type())
		}
	}
	return nil
}
func scanRecordsToStructs(records []Record, ptrs []any) error {
	fields := mapStructFieldsOnce(ptrs[0])
	for i, record := range records {
		if i >= len(ptrs) {
			break
		}
		if err := scanRecordToStruct(record, ptrs[i], fields); err != nil {
			return err
		}
	}
	return nil
}
