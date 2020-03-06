package convert

import (
	"database/sql"
	"time"
)

//因为[]byte不能接收数据库为空的值,这个思路废弃了......
type columnValue []byte

func (v columnValue) Bytes() []byte {
	return []byte(v)
}

func (v columnValue) String() string {
	return string(v)
}

func (v columnValue) NullString() sql.NullString {
	if v == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	} else {
		return sql.NullString{
			String: string(v),
			Valid:  true,
		}
	}
}

func (v columnValue) Bool() bool {
	return Bool(v)
}

func (v columnValue) NullBool() sql.NullBool {
	if v == nil {
		return sql.NullBool{
			Bool:  false,
			Valid: false,
		}
	} else {
		return sql.NullBool{
			Bool:  Bool(v),
			Valid: true,
		}
	}
}

func (v columnValue) Int() int {
	return Int(v)
}

func (v columnValue) Int8() int8 {
	return Int8(v)
}

func (v columnValue) Int16() int16 {
	return Int16(v)
}

func (v columnValue) Int32() int32 {
	return Int32(v)
}

func (v columnValue) NullInt32() sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	} else {
		return sql.NullInt32{
			Int32: Int32(v),
			Valid: true,
		}
	}
}

func (v columnValue) Int64() int64 {
	return Int64(v)
}

func (v columnValue) NullInt64() sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	} else {
		return sql.NullInt64{
			Int64: Int64(v),
			Valid: true,
		}
	}
}

func (v columnValue) Uint() uint {
	return Uint(v)
}

func (v columnValue) Uint8() uint8 {
	return Uint8(v)
}

func (v columnValue) Uint16() uint16 {
	return Uint16(v)
}

func (v columnValue) Uint32() uint32 {
	return Uint32(v)
}

func (v columnValue) Uint64() uint64 {
	return Uint64(v)
}

func (v columnValue) Float32() float32 {
	return Float32(v)
}

func (v columnValue) Float64() float64 {
	return Float64(v)
}

func (v columnValue) NullFloat64() sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	} else {
		return sql.NullFloat64{
			Float64: Float64(v),
			Valid:   true,
		}
	}
}

func (v columnValue) Time(format string, TZLocation ...*time.Location) time.Time {
	return Time(v, format, TZLocation...)
}

func (v columnValue) NullTime(format string, TZLocation ...*time.Location) sql.NullTime {
	if v == nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	} else {
		return sql.NullTime{
			Time:  Time(v, format, TZLocation...),
			Valid: true,
		}
	}

}

func (v columnValue) TimeDuration() time.Duration {
	return TimeDuration(v)
}
