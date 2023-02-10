#### Usage

see [example](https://github.com/Vnaki/gt/tree/master/example)

```go 
package main

import (
	"fmt"
	"github.com/vnaki/gt"
	"time"
)

type Model struct {
	Id        int32      `db:"id,omitempty" gen:"length:10,pk,ai"`
	SDK       bool       `db:"sdk" gen:"type:tinyint,length:1"`
	CreatedAt time.Time  `db:"created_at"`
	UpdateAt  *time.Time `db:"updated_at"`
}

type ThreeStudentModel struct {
	Model
	Num     uint64  `db:"num" gen:"notnull,default:0"`
	Name    string  `db:"name" gen:"notnull,default:"`
	Content string  `db:"content" gen:"type:text"`
	Score   float32 `db:"score" gen:"length:1,decimal:1,default:1,notnull,unsigned"`
	Money   float64 `db:"money" gen:"length:10,decimal:2,default:1,notnull,unsigned"`
}

type TwoStudent struct {
	Model
	Name    string `db:"name" gen:"notnull"`
	Content string `db:"content"`
	Score   int    `db:"score" gen:"length:1,decimal:1,default:1,notnull,unsigned"`
}

func main() {
	b := gt.New()
	b.SetSchema("stu")
	b.SetWrap(true)
	sql, err := b.Model(ThreeStudentModel{})
	fmt.Println(sql, err)

	sql, err = b.Model(TwoStudent{}, "twostu")
	fmt.Println(sql, err)

	b = gt.New()
	b.SetWrap(true)
	b.SetMode(gt.MYSQL)

	sql, err = b.Model(ThreeStudentModel{})
	fmt.Println(sql, err)
	sql, err = b.Model(TwoStudent{}, "twostu")
	fmt.Println(sql, err)
}

```

result output

```sql

-- sqlite
CREATE TABLE stu.three_student(
    'id' integer PRIMARY KEY AUTOINCREMENT,
    'sdk' tinyint(1),
    'created_at' datetime,
    'updated_at' datetime,
    'num' bigint NOT NULL DEFAULT 0,
    'name' varchar NOT NULL DEFAULT '',
    'content' text,
    'score' float(1,1) NOT NULL DEFAULT 1,
    'money' double(10,2) NOT NULL DEFAULT 1
);

CREATE TABLE stu.twostu(
    'id' integer PRIMARY KEY AUTOINCREMENT,
    'sdk' tinyint(1),
    'created_at' datetime,
    'updated_at' datetime,
    'name' varchar NOT NULL,
    'content' varchar,
    'score' bigint(1,1) NOT NULL DEFAULT 1
);

-- mysql
CREATE TABLE three_student(
  `id` int(10) PRIMARY KEY AUTO_INCREMENT,
  `sdk` tinyint(1),
  `created_at` datetime,
  `updated_at` datetime,
  `num` bigint NOT NULL DEFAULT 0,
  `name` varchar NOT NULL DEFAULT '',
  `content` text,
  `score` float(1,1) UNSIGNED NOT NULL DEFAULT 1,
  `money` double(10,2) UNSIGNED NOT NULL DEFAULT 1
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

CREATE TABLE twostu(
    `id` int(10) PRIMARY KEY AUTO_INCREMENT,
   `sdk` tinyint(1),
   `created_at` datetime,
   `updated_at` datetime,
   `name` varchar NOT NULL,
   `content` varchar,
   `score` bigint(1,1) UNSIGNED NOT NULL DEFAULT 1
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
```
#### Mode 

- MYSQL
- SQLITE

#### Tag `db`

Corresponding data table column name, `Id` to `id`, `Content` to `content` 

```go 
type People struct {
    Id        int32  `db:"id,omitempty" gen:"pk,ai"`
    Content   string `db:"content" gen:"type:text"`
}


```

#### Tag `gen`

| 属性 | 默认值 | 说明 |
| --- | --- | --- |
| type | | 原生sql数据类型:char,text,mediumint,timestamp,datetime 等 |
| length | | 数据长度 |
| decimal | 2 | 浮点类型精度 |
| default | | 默认值 |
| pk | | 主键 |
| ai | | 自增 |
| comment | | 注释 |
| unsigned | | 无符号 |
| notnull | | not null |

#### Integer Data Type

| 数据库数据类型 | 范围 | 无符号范围 | 数据类型 |
| --- | --- | --- | --- |
| TINYINT | -128〜127 | 0 〜255 | int8/uint8 |
| SMALLINT | -32768〜32767 | 0〜65535 | int16/uint16|
| INT (INTEGER) | -2147483648〜2147483647 | 0〜4294967295 | int32/uint32|
| BIGINT | -9223372036854775808〜9223372036854775807 | 0〜18446744073709551615 | int64 int / uint64 uint|

#### String Data Type

``` 
string -> varchar
```

```
// int int8 int16 int32 int64 byte rune
// uint uint8 uint16 uint32 uint64 byte rune
// float32 float64
// char varchar text
// datetime timestamp

// TINYINT	-128〜127	0 〜255        int8
// SMALLINT	-32768〜32767	0〜65535   int16
// MEDIUMINT	-8388608〜8388607	0〜16777215
// INT (INTEGER)	-2147483648〜2147483647	0〜4294967295   int32
// BIGINT	-9223372036854775808〜9223372036854775807	0〜18446744073709551615 int64 int
//

```
