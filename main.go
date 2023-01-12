package main

import (
	"fmt"
	"reflect"
	"strings"
)

type People struct {
	Id        int32  `db:"id,omitempty" gen:"pk,ai"`
	Content   string `db:"content" gen:"text"`
	CreatedAt string `db:"created_at"`
}

// smallint/int/big/int
// char/varchar
type ThreeStudent struct {
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

func main() {
	b := New("user")

	sql, err := b.Bind(ThreeStudent{})

	fmt.Println(sql, err)

	q := New("user")
	q.SetMode(MYSQL)

	sql, err = q.Bind(ThreeStudent{})
	fmt.Println(sql, err)
}

type Mode int8

const (
	SQLITE Mode = iota
	MYSQL
)

type Schema struct {
	sql    []string
	schema string
	quote  string
	mode   Mode
}

func New(schema string) *Schema {
	return &Schema{
		sql:    []string{},
		mode:   SQLITE,
		schema: schema,
		quote:  "'",
	}
}

func (sc *Schema) SetMode(mode Mode) {
	sc.mode = mode

	if mode == MYSQL {
		sc.quote = "`"
	} else if mode == SQLITE {
		sc.quote = "'"
	}
}

func (sc *Schema) Bind(i interface{}, table ...string) (string, error) {
	// NewBinder()
	b := NewBinder(sc.mode, sc.quote)

	if len(table) > 0 {
		b.SetTable(table[0])
	}

	return b.Generate(i)
}

type Binder struct {
	mode  Mode
	quote string
	table string
	sql   []string
}

func NewBinder(mode Mode, quote string) *Binder {
	return &Binder{
		sql:   []string{},
		mode:  mode,
		quote: quote,
	}
}

func (b *Binder) SetTable(table string) {
	b.table = table
}

func (b *Binder) Generate(model interface{}) (string, error) {
	t := reflect.TypeOf(model)

	if kind := t.Kind().String(); kind != "struct" {
		return "", fmt.Errorf("unsupported type %v, only type struct is supported", kind)
	}

	if t.NumField() == 0 {
		return "", fmt.Errorf("struct %v empty field", t.Name())
	}

	b.parse(t)

	sf := ""

	if b.mode == MYSQL {
		sf = " ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4"
	}

	if b.table == "" {
		b.table = b.snake(t.Name())
	}

	sql := fmt.Sprintf(
		"CREATE TABLE %v(%v)%v;",
		fmt.Sprintf("%v%v%v", b.quote, b.table, b.quote),
		strings.Join(b.sql, ","),
		sf,
	)

	return sql, nil
}

func (b *Binder) parse(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			b.parse(field.Type)
		} else {
			b.sql = append(b.sql, b.parseField(field))
		}
	}
}

func (b *Binder) parseField(field reflect.StructField) string {
	t := field.Tag.Get("db")
	if t == "" {
		return ""
	}

	// name
	name := strings.SplitN(t, ",", 2)[0]
	if name == "omitempty" {
		return ""
	}

	name = fmt.Sprintf("%v%v%v", b.quote, name, b.quote)

	// parse sql params
	gen := b.parseGen(field.Type.Name(), field.Tag.Get("gen"))

	return fmt.Sprintf("%v %v", name, gen)
}

func (b *Binder) parseGen(typ, gen string) string {
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
	} else if b.isInt(typ) {
		var length string

		r = b.covert(typ)

		if v, ok := kv["length"]; ok && v != "" {
			length = v
		}

		if length != "" {
			r = fmt.Sprintf("%v(%v)", r, length)
		}
	} else if b.isFloat(typ) {
		var (
			length  string
			decimal = "2"
		)

		r = b.covert(typ)

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

		r = b.covert(typ)

		if v, ok := kv["length"]; ok && v != "" {
			length = v
		}

		if length != "" {
			r = fmt.Sprintf("%v(%v)", r, length)
		}
	}

	if b.contain("unsigned", ex) {
		r = fmt.Sprintf("%v UNSIGNED", r)
	}

	if b.contain("notnull", ex) {
		r = fmt.Sprintf("%v NOT NULL", r)
	}

	if b.contain("pk", ex) {
		r = fmt.Sprintf("%v PRIMARY KEY", r)
	}

	if b.contain("ai", ex) {
		r = fmt.Sprintf("%v AUTO_INCREMENT", r)
	}

	if v, ok := kv["default"]; ok && v != "" {
		r = fmt.Sprintf("%v DEFAULT %v", r, v)
	}

	return r
}

func (b *Binder) isInt(v string) bool {
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

func (b *Binder) isFloat(v string) bool {
	switch v {
	case "float32":
		fallthrough
	case "float64":
		return true
	}

	return false
}

func (b *Binder) covert(v string) string {
	var kv = map[string]string{
		"int":     "bigint",
		"int8":    "tinyint",
		"int16":   "smallint",
		"int32":   "int",
		"int64":   "bigint",
		"byte":    "tinyint",
		"rune":    "int",
		"float32": "float",  // 单精度
		"float64": "double", // 双精度
		"string":  "varchar",
		"char":    "char",
	}

	return kv[v]
}

func (b *Binder) contain(v string, arr []string) bool {
	for _, v1 := range arr {
		if v == v1 {
			return true
		}
	}
	return false
}

func (b *Binder) snake(v string) string {
	v = strings.TrimRight(v, "Model")

	d := make([]byte, len(v))

	for i := 0; i < len(v); i++ {
		if v[i] >= 'A' && v[i] <= 'Z' {
			if i > 0 {
				d = append(d, '_')
			}

			d = append(d, v[i]+'a'-'A')
		} else {
			d = append(d, v[i])
		}
	}

	return string(d)
}
