package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/xiaoma03xf/sharddoc/kv"
)

func assert(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}

// table schema
type TableDef struct {
	// user defined
	Name    string
	Types   []uint32   // column types
	Cols    []string   // column names
	Indexes [][]string // the first index is the primary key
	// auto-assigned B-tree key prefixes for different tables and indexes
	Prefixes []uint32
}

const (
	TYPE_ERROR = 0 // uninitialized
	TYPE_BYTES = 1
	TYPE_INT64 = 2
	TYPE_INF   = 0xff // do not use
)

// table cell
type Value struct {
	Type uint32
	I64  int64
	Str  []byte
}

// table row
type Record struct {
	Cols []string
	Vals []Value
}

func (rec *Record) AddStr(col string, val []byte) *Record {
	rec.Cols = append(rec.Cols, col)
	rec.Vals = append(rec.Vals, Value{Type: TYPE_BYTES, Str: val})
	return rec
}
func (rec *Record) AddInt64(col string, val int64) *Record {
	rec.Cols = append(rec.Cols, col)
	rec.Vals = append(rec.Vals, Value{Type: TYPE_INT64, I64: val})
	return rec
}

func (rec *Record) Get(key string) *Value {
	for i, c := range rec.Cols {
		if c == key {
			return &rec.Vals[i]
		}
	}
	return nil
}

// extract multiple column values
func getValues(tdef *TableDef, rec Record, cols []string) ([]Value, error) {
	vals := make([]Value, len(cols))
	for i, c := range cols {
		v := rec.Get(c)
		if v == nil {
			return nil, fmt.Errorf("missing column: %s", tdef.Cols[i])
		}
		if v.Type != tdef.Types[slices.Index(tdef.Cols, c)] {
			return nil, fmt.Errorf("bad column type: %s", c)
		}
		vals[i] = *v
	}
	return vals, nil
}

// escape the null byte so that the string contains no null byte.
func escapeString(in []byte) []byte {
	toEscape := bytes.Count(in, []byte{0}) + bytes.Count(in, []byte{1})
	if toEscape == 0 {
		return in // fast path: no escape
	}

	out := make([]byte, len(in)+toEscape)
	pos := 0
	for _, ch := range in {
		if ch <= 1 {
			// using 0x01 as the escaping byte:
			// 00 -> 01 01
			// 01 -> 01 02
			out[pos+0] = 0x01
			out[pos+1] = ch + 1
			pos += 2
		} else {
			out[pos] = ch
			pos += 1
		}
	}
	return out
}

func unescapeString(in []byte) []byte {
	if bytes.Count(in, []byte{1}) == 0 {
		return in // fast path: no unescape
	}

	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] == 0x01 {
			// 01 01 -> 00
			// 01 02 -> 01
			i++
			assert(in[i] == 1 || in[i] == 2)
			out = append(out, in[i]-1)
		} else {
			out = append(out, in[i])
		}
	}
	return out
}

// order-preserving encoding
func encodeValues(out []byte, vals []Value) []byte {
	for _, v := range vals {
		out = append(out, byte(v.Type)) // doesn't start with 0xff
		switch v.Type {
		case TYPE_INT64:
			var buf [8]byte
			u := uint64(v.I64) + (1 << 63)        // flip the sign bit
			binary.BigEndian.PutUint64(buf[:], u) // big endian
			out = append(out, buf[:]...)
		case TYPE_BYTES:
			out = append(out, escapeString(v.Str)...)
			out = append(out, 0) // null-terminated
		default:
			panic("what?")
		}
	}
	return out
}

// for primary keys and indexes
func encodeKey(out []byte, prefix uint32, vals []Value) []byte {
	// 4-byte table prefix
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], prefix)
	out = append(out, buf[:]...)
	// order-preserving encoded keys
	out = encodeValues(out, vals)
	return out
}

// for the input range, which can be a prefix of the index key.
func encodeKeyPartial(
	out []byte, prefix uint32, vals []Value, cmp int,
) []byte {
	out = encodeKey(out, prefix, vals)
	if cmp == kv.CMP_GT || cmp == kv.CMP_LE { // encode missing columns as infinity
		out = append(out, 0xff) // unreachable +infinity
	} // else: -infinity is the empty string
	return out
}

func decodeValues(in []byte, out []Value) {
	for i := range out {
		assert(out[i].Type == uint32(in[0]))
		in = in[1:]
		switch out[i].Type {
		case TYPE_INT64:
			u := binary.BigEndian.Uint64(in[:8])
			out[i].I64 = int64(u - (1 << 63))
			in = in[8:]
		case TYPE_BYTES:
			idx := bytes.IndexByte(in, 0)
			assert(idx >= 0)
			out[i].Str = unescapeString(in[:idx])
			in = in[idx+1:]
		default:
			panic("what?")
		}
	}
	assert(len(in) == 0)
}

func decodeKey(in []byte, out []Value) {
	decodeValues(in[4:], out)
}

// internal table: metadata
var TDEF_META = &TableDef{
	Name:     "@meta",
	Types:    []uint32{TYPE_BYTES, TYPE_BYTES},
	Cols:     []string{"key", "val"},
	Indexes:  [][]string{{"key"}},
	Prefixes: []uint32{1},
}

// internal table: table schemas
var TDEF_TABLE = &TableDef{
	Name:     "@table",
	Types:    []uint32{TYPE_BYTES, TYPE_BYTES},
	Cols:     []string{"name", "def"},
	Indexes:  [][]string{{"name"}},
	Prefixes: []uint32{2},
}

var INTERNAL_TABLES map[string]*TableDef = map[string]*TableDef{
	"@meta":  TDEF_META,
	"@table": TDEF_TABLE,
}
