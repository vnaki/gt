package main

import (
	"fmt"
	"github.com/vnaki/gt"
	"time"
)

type Model struct {
	Id        uint32     `db:"id,omitempty" gen:"length:10,pk,ai,unsigned"`
	CreatedAt time.Time  `db:"created_at" gen:"notnull"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type UserModel struct {
	Model
	Username string  `db:"username" gen:"length:10,comment:'用户名称',notnull"`
	Content  string  `db:"content" gen:"type:text"`
	Email    string  `db:"email" gen:"length:100,notnull"`
	Phone    string  `db:"phone" gen:"type:char,length:11,notnull"`
	Score    float32 `db:"score" gen:"length:10,decimal:2,default:1,notnull,unsigned"`
	Money    float64 `db:"money" gen:"length:10,decimal:2,default:1,notnull,unsigned"`
	Status   uint8   `db:"status" gen:"length:2,notnull,unsigned"`
}

// gen: length:1,decimal:2,default:111,pk,ai,unsigned,notnull

// int int8 int16 int32 int64 byte rune
// uint uint8 uint16 uint32 uint64 byte rune
// float32 float64
// char varchar text
// datetime

// TINYINT	-128〜127	0 〜255        int8
// SMALLINT	-32768〜32767	0〜65535   int16
// MEDIUMINT	-8388608〜8388607	0〜16777215
// INT (INTEGER)	-2147483648〜2147483647	0〜4294967295   int32
// BIGINT	-9223372036854775808〜9223372036854775807	0〜18446744073709551615 int64 int
//

func main() {
	b := gt.New()
	b.SetSchema("student")
	b.SetWrap(true)
	b.SetDrop(false)
	b.SetMode(gt.MYSQL)
	b.SetSuffix("Model")
	ss, err := b.Model(UserModel{})
	if err != nil {
		panic(err)
	} else {
		for _, s := range ss {
			fmt.Println(s)
		}
	}
}
