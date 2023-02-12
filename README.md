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

```

result output

```sql

-- sqlite
CREATE TABLE student.user(
  'id' integer PRIMARY KEY AUTOINCREMENT,
  'created_at' datetime NOT NULL,
  'updated_at' datetime,
  'deleted_at' datetime,
  'username' varchar(10) NOT NULL, -- '用户名称'
  'content' text,
  'email' varchar(100) NOT NULL,
  'phone' char(11) NOT NULL,
  'score' float(10,2) NOT NULL DEFAULT 1,
  'money' double(10,2) NOT NULL DEFAULT 1,
  'status' tinyint(2) NOT NULL
);

-- mysql

CREATE TABLE student.user(
    `id` int(10) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `created_at` datetime NOT NULL,
    `updated_at` datetime,
    `deleted_at` datetime,
    `username` varchar(10) NOT NULL COMMENT '用户名称',
    `content` text,
    `email` varchar(100) NOT NULL,
    `phone` char(11) NOT NULL,
    `score` float(10,2) UNSIGNED NOT NULL DEFAULT 1,
    `money` double(10,2) UNSIGNED NOT NULL DEFAULT 1,
    `status` tinyint(2) UNSIGNED NOT NULL
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