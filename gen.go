package gt

import (
	"fmt"
	"reflect"
	"strings"
)

type Mode int8

type GTable struct {
	mode  Mode
	quote string
	schema string
	suffix string
	wrap bool
}

const (
	SQLITE Mode = iota
	MYSQL
)

func New() *GTable {
	return &GTable{
		mode:  SQLITE,
		quote: "'",
		suffix: "Model",
		wrap: true,
	}
}

func (b *GTable) SetWrap(wrap bool) {
	b.wrap = wrap
}

func (b *GTable) SetSuffix(suffix string) {
	b.suffix = suffix
}

func (b *GTable) SetSchema(schema string) {
	b.schema = schema
}

func (b *GTable) SetMode(mode Mode) {
	b.mode = mode

	if mode == MYSQL {
		b.quote = "`"
	} else if mode == SQLITE {
		b.quote = "'"
	}
}

func (b *GTable) Model(model interface{}, table ...string) (string, error) {
	t := reflect.TypeOf(model)

	if kind := t.Kind().String(); kind != "struct" {
		return "", fmt.Errorf("unsupported type %v, only type struct is supported", kind)
	}

	if t.NumField() == 0 {
		return "", fmt.Errorf("struct %v empty field", t.Name())
	}

	columns := b.parse(t)

	sf := ""

	if b.mode == MYSQL {
		sf = " ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4"
	}

	if len(table) == 0 || table[0] == "" {
		table = []string{ b.snake(t.Name()) }
	}

	sep := ","

	if b.wrap {
		sep = ",\n"
	}

	sql := strings.Join(columns, sep)
	if b.wrap {
		sql = fmt.Sprintf("%v%v%v", "\n", sql, "\n")
	}

	tb := fmt.Sprintf("%v%v%v", b.quote, table[0], b.quote)
	if b.schema != "" {
		tb = fmt.Sprintf("%v%v%v.%v", b.quote, b.schema, b.quote, tb)
	}

	return fmt.Sprintf("CREATE TABLE %v(%v)%v;", tb,	sql, sf), nil
}

func (b *GTable) parse(t reflect.Type) (columns []string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			columns = append(columns, b.parse(field.Type)...)
		} else {
			columns = append(columns, b.parseField(field))
		}
	}
	return
}

func (b *GTable) parseField(field reflect.StructField) string {
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

	// parse gen
	gen := b.parseGen(field.Type.Name(), field.Tag.Get("gen"))

	return fmt.Sprintf("%v %v", name, gen)
}

func (b *GTable) parseGen(typ, gen string) string {
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

	if v, ok := kv["default"]; ok {
		r = fmt.Sprintf("%v DEFAULT %v", r, v)
	}

	return r
}

func (b *GTable) isInt(v string) bool {
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

func (b *GTable) isFloat(v string) bool {
	switch v {
	case "float32":
		fallthrough
	case "float64":
		return true
	}

	return false
}

func (b *GTable) covert(v string) string {
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
	}
	return kv[v]
}

func (b *GTable) contain(v string, arr []string) bool {
	for _, v1 := range arr {
		if v == v1 {
			return true
		}
	}
	return false
}

func (b *GTable) snake(v string) string {
	v = strings.TrimRight(v, b.suffix)

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
