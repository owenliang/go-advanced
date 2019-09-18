package main

import (
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
func Query(result interface{}, sql string, values ...interface{}) error {
	// type1是*[]*User
	type1 := reflect.TypeOf(result)
	if type1.Kind() != reflect.Ptr {
		return errors.New("第一个参数必须是指针")
	}

	// type2是[]*User
	type2 := type1.Elem()	// 解指针后的类型
	if type2.Kind() != reflect.Slice {
		return errors.New("第一个参数必须指向切片")
	}

	// type3是*User
	type3 := type2.Elem()
	if type3.Kind() != reflect.Ptr {
		return errors.New("切片元素必须是指针类型")
	}

	// 发起SQL查询
	rows, _ := db.Raw(sql, values...).Rows()
	for rows.Next() {
		//  type3.Elem()是User, elem是*User
		elem := reflect.New(type3.Elem())
		// 传入*User
		db.ScanRows(rows, elem.Interface())
		// reflect.ValueOf(result).Elem()是[]*User，Elem是*User，newSlice是[]*User
		newSlice := reflect.Append(reflect.ValueOf(result).Elem(), elem)
		// 扩容后的slice赋值给*result
		// reflect.ValueOf(result).Elem()是[]User
		reflect.ValueOf(result).Elem().Set(newSlice)
	}
	return nil
}

func main() {
	db, _ = gorm.Open("mysql", "root:baidu@/gin?charset=utf8&parseTime=True&loc=Local")

	result := []*User{}
	if err := Query(&result, "select * from user where name=?", "owen"); err == nil {
		for i := 0; i < len(result); i++ {
			fmt.Println(*result[i])
		}
	}
}