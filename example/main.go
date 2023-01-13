package main

import (
	"fmt"
	"gt"
)

type People struct {
	Id        int32  `db:"id,omitempty" gen:"pk,ai"`
	Content   string `db:"content" gen:"type:text"`
	CreatedAt string `db:"created_at"`
}

type ThreeStudentModel struct {
	People
	Name  string `db:"name" gen:"notnull"`
	Score int    `db:"score" gen:"length:1,decimal:1,default:1,notnull,unsigned"`
}

type TwoStudent struct {
	People
	Name  string `db:"name" gen:"notnull"`
	Score int    `db:"score" gen:"length:1,decimal:1,default:1,notnull,unsigned"`
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
	b.SetSchema("stu")
	b.SetWrap(false)
	sql, err := b.Model(ThreeStudentModel{})
	fmt.Println(sql, err)

	sql, err = b.Model(TwoStudent{}, "twostu")
	fmt.Println(sql, err)

	b = gt.New()
	b.SetMode(gt.MYSQL)

	sql, err = b.Model(TwoStudent{})
	fmt.Println(sql, err)
}
