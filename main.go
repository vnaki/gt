package main

import (
	"fmt"
	"reflect"
	"strings"
)

type People struct {
	Id int32 `yaml:"id,omitempty"`
}

// smallint/int/big/int
// char/varchar
type Student struct {
	People
	Name  string `db:"name" gen:"notnull"`
	Score int    `db:"score" gen:"length:1,decimal:1,default:1,notnull,unsigned"`
}

// length:"1" decimal:"2" default:"1" required:"true" notnull:"true"

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

// CREATE TABLE NAME(
//
//
//
//);

func main() {

	sc := Schema{}

	sc.Model(Student{})

	//s := Student{}
	//
	//t := reflect.TypeOf(s)
	//
	//// 读取字段
	//fields := t.Field(1)
	//// 读取Tag
	//tag := fields.Tag.Get("ini")
	//
	//v, _ := json.Marshal(fields)
	//
	//fmt.Println(t.Name(), t.Kind(), string(v), tag)
}

type Mode int8

const (
	SQLITE Mode = iota
	MYSQL
)

type Schema struct {
	sql []string
	schema string
	mode Mode
}

func New(schema string) *Schema {
	return &Schema{
		sql: []string{},
		mode: SQLITE,
		schema: schema,
	}
}

func (sc *Schema) Generate() {

}

func (sc *Schema) Mode(mode Mode) {
	sc.mode = mode
}

func (sc *Schema) Model(i interface{}) {
	if err := sc.model(i); err != nil {
		panic(err)
	}
}

func (sc *Schema) model(model interface{}) error {
	t := reflect.TypeOf(model)

	if kind := t.Kind().String(); kind != "struct" {
		return fmt.Errorf("unsupported type %v, only type struct is supported", kind)
	}

	if t.NumField() == 0 {
		return fmt.Errorf("struct %v empty field", t.Name())
	}

	fmt.Println("table:", t.Name(), t)


	for i:=0;i<t.NumField();i++ {
		v := sc.parseField(t.Field(i))
		fmt.Println(">>>>>", v)

		sc.sql = append(sc.sql, v)
	}
	//template := fmt.Sprintf("CREATE TABLE %v;", t.Name())

	// PRIMARY KEY (`id`),
	// ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='工作节点'
	// KEY `idx_deleted` (`deleted_at`) USING BTREE,
	return nil
}

func (sc *Schema) parseField(field reflect.StructField) string {
	//field.Tag.Get()
	// column2 datatype attributes,

	t := field.Tag.Get("db")
	if t == "" {
		return ""
	}

	// name
	name := strings.SplitN(t, ",", 2)[0]
	if name == "omitempty" {
		return ""
	}

	// type


	gen := sc.parseGen(field.Type.Name(), field.Tag.Get("gen"))

	return fmt.Sprintf("%v %v", name, gen)
}

func (sc *Schema) parseGen(typ, gen string) string {
	var (
		ex []string
		kv = make(map[string]string)
	)

	for _, v := range strings.Split(gen, ",") {
		sn := strings.SplitN(v, ":", 2)

		if len(sn) == 2 {
			kv[sn[0]] = sn[1]
		} else {
			ex = append(ex, sn[0])
		}
	}

	var r string

	if v, ok := kv["type"]; ok && v != "" {
		r = v
	} else if sc.isInt(typ) {
		var length string

		r = sc.covert(typ)

		if v, ok := kv["length"]; ok && v != "" {
			length = v
		}

		r = fmt.Sprintf("%v(%v)", r, length)
	} else if sc.isFloat(typ) {
		var (
			length string
			decimal = "2"
		)

		r = sc.covert(typ)

		if v, ok := kv["length"]; ok && v != "" {
			length = v

			if v, ok := kv["decimal"]; ok && v != "" {
				decimal = v
			}
		}

		if length != "" {
			r = fmt.Sprintf("%v(%v,%v)", r, length, decimal)
		}
	} else if typ == "string" || typ == "char" {
		var (
			length string
		)

		r = sc.covert(typ)

		if v, ok := kv["length"]; ok && v != "" {
			length = v
		}

		if length != "" {
			r = fmt.Sprintf("%v(%v)", r, length)
		}
	}

	if sc.contain("unsigned", ex) {
		r = fmt.Sprintf("%v UNSIGNED", r)
	}

	if sc.contain("notnull", ex) {
		r = fmt.Sprintf("%v NOT NULL", r)
	}

	if v, ok := kv["default"]; ok && v != "" {
		r = fmt.Sprintf("%v DEFAULT %v", r, v)
	}

	return r
}

func (sc *Schema) isInt(v string) bool {
	switch v {
	case "int":
		fallthrough
	case "int8":
		fallthrough
	case "int16":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		fallthrough
	case "byte":
		fallthrough
	case "rune":
		return true
	}

	return false
}

func (sc *Schema) isFloat(v string) bool {
	switch v {
	case "float32":
		fallthrough
	case "float64":
		return true
	}

	return false
}

func (sc *Schema) covert(v string) string {
	var kv = map[string]string{
		"int": "BIGINT",
		"int8": "TINYINT",
		"int16": "SMALLINT",
		"int32": "INT",
		"int64": "BIGINT",
		"byte": "TINYINT",
		"rune": "INT",
		"float32": "FLOAT", // 单精度
		"float64": "DOUBLE", // 双精度
		"string": "VARCHAR",
		"char": "CHAR",
	}

	return kv[v]
}


func (sc *Schema) covertFloat(v string) string {
	var kv = map[string]string{
		"float32": "FLOAT",
		"float64": "DOUBLE",
	}

	return kv[v]
}

func (sc *Schema) contain(v string, arr []string) bool {
	for _, v1 := range arr {
		if v == v1 {
			return true
		}
	}
	return false
}