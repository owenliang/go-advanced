package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
)

type User struct {
	Id int
	Name string
}

var (
	db *gorm.DB
)

// 查询SQL
func Query(ctx context.Context, result interface{}, sql string, values ...interface{}) error {
	// result必须是数组的指针, 例如*[]User
	type1 := reflect.TypeOf(result)
	if type1.Kind() != reflect.Ptr {
		return errors.New("第一个参数必须是指针")
	}

	// *result必须是数组或者切片
	type2 := type1.Elem()	// 解指针后的类型
	if type2.Kind() != reflect.Slice {
		return errors.New("第一个参数必须指向切片")
	}

	// 得到*result[i]的类型
	type3 := type2.Elem()
	if type3.Kind() != reflect.Ptr {
		return errors.New("切片元素必须是指针类型")
	}

	// 发起SQL查询
	rows, _ := db.Raw(sql, values...).Rows()
	for rows.Next() {
		// 新建一个User，返回其指针
		elem := reflect.New(type3.Elem())
		// 传入*User
		db.ScanRows(rows, elem.Interface())
		// 结果append到*result的切片中
		newSlice := reflect.Append(reflect.ValueOf(result).Elem(), elem)
		// 扩容后的slice赋值给*result
		reflect.ValueOf(result).Elem().Set(newSlice)
	}
	return nil
}

func main() {
	db, _ = gorm.Open("mysql", "root:baidu@123@/gin?charset=utf8&parseTime=True&loc=Local")

	result := []*User{}
	if err := Query(context.Background(), &result, "select * from user where name=?", "owen"); err == nil {
		for i := 0; i < len(result); i++ {
			fmt.Println(*result[i])
		}
	}
}