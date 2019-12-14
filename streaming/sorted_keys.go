package streaming

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/cube2222/octosql"
	"github.com/pkg/errors"
)

func SortedMarshal(v *octosql.Value) []byte {
	switch v.GetType() {
	case octosql.TypeString:
		return SortedMarshalString(v.AsString())
	case octosql.TypeInt:
		return SortedMarshalInt(v.AsInt())
	case octosql.TypeBool, octosql.TypePhantom, octosql.TypeNull:
		return SortedMarshalBool(v.AsBool())
	case octosql.TypeTime:
		return SortedMarshalTime(v.AsTime())
	case octosql.TypeDuration:
		return SortedMarshalDuration(v.AsDuration())
	case octosql.TypeFloat: // TODO - fix
		return SortedMarshalFloat(v.AsFloat())
	default: //TODO: add other types - tuple/object
		return nil
	}
}

func SortedUnmarshal(b []byte, v *octosql.Value) error {
	switch v.GetType() {
	case octosql.TypeString:
		u, err := SortedUnmarshalString(b)
		if err != nil {
			return err
		}

		*v = octosql.MakeString(u)
		return nil
	case octosql.TypeInt:
		u, err := SortedUnmarshalInt(b)
		if err != nil {
			return err
		}

		*v = octosql.MakeInt(u)
		return nil
	case octosql.TypeBool, octosql.TypePhantom, octosql.TypeNull:
		u, err := SortedUnmarshalBool(b)
		if err != nil {
			return err
		}

		*v = octosql.MakeBool(u)
	case octosql.TypeTime:
		u, err := SortedUnmarshalTime(b)
		if err != nil {
			return err
		}

		*v = octosql.MakeTime(u)
	case octosql.TypeDuration:
		u, err := SortedUnmarshalDuration(b)
		if err != nil {
			return err
		}

		*v = octosql.MakeDuration(u)
	case octosql.TypeFloat: // TODO - fix
		u, err := SortedUnmarshalFloat(b)
		if err != nil {
			return err
		}

		*v = octosql.MakeFloat(u)
	default: //TODO: add other types
		return nil
	}

	panic("unreachable")
}

/* Marshal string */
func SortedMarshalString(s string) []byte {
	return []byte(s)
}

func SortedUnmarshalString(b []byte) (string, error) {
	return string(b), nil
}

/* Marshal int and int64 */
func SortedMarshalInt(i int) []byte {
	return SortedMarshalInt64(int64(i))
}

func SortedMarshalInt64(i int64) []byte {
	return SortedMarshalUint64(uint64(i), i >= 0)
}

func SortedMarshalUint64(ui uint64, sign bool) []byte {
	b := make([]byte, 9)

	binary.LittleEndian.PutUint64(b, ui)

	if sign {
		b[8] = 1
	} else {
		b[8] = 0
	}

	return reverseByteSlice(b)
}

func SortedUnmarshalInt(b []byte) (int, error) {
	if len(b) != 9 {
		return 0, errors.New("incorrect int key size")
	}
	return int(binary.LittleEndian.Uint64(reverseByteSlice(b[1:]))), nil
}

func SortedUnmarshalInt64(b []byte) (int64, error) {
	value, err := SortedUnmarshalUint64(b)
	if err != nil {
		return 0, errors.Wrap(err, "incorrect int64 key size")
	}

	return int64(value), nil
}

func SortedUnmarshalUint64(b []byte) (uint64, error) {
	if len(b) != 9 {
		return 0, errors.New("incorrect uint64 key size")
	}
	return binary.LittleEndian.Uint64(reverseByteSlice(b[1:])), nil
}

func reverseByteSlice(b []byte) []byte {
	c := make([]byte, len(b))

	for i, j := 0, len(b)-1; i <= j; i, j = i+1, j-1 {
		c[i] = b[j]
		c[j] = b[i]
	}

	return c
}

/* Marshal bool */
func SortedMarshalBool(b bool) []byte {
	if b {
		return []byte{1}
	}

	return []byte{0}
}

func SortedUnmarshalBool(b []byte) (bool, error) {
	if len(b) != 1 {
		return false, errors.New("incorrect bool key size")
	}

	switch b[0] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, errors.New("incorrect bool key value")
	}
}

/* Marshal Timestamp */
func SortedMarshalTime(t time.Time) []byte {
	return SortedMarshalInt64(t.UnixNano())
}

func SortedUnmarshalTime(b []byte) (time.Time, error) {
	value, err := SortedUnmarshalInt64(b)
	if err != nil {
		return time.Now(), errors.Wrap(err, "incorrect time key representation")
	}

	return time.Unix(value/int64(time.Second), value%int64(time.Second)), nil
}

/* Marshal Duration */
func SortedMarshalDuration(d time.Duration) []byte {
	return SortedMarshalInt64(d.Nanoseconds())
}

func SortedUnmarshalDuration(b []byte) (time.Duration, error) {
	value, err := SortedUnmarshalInt64(b)
	if err != nil {
		return time.Duration(0), errors.Wrap(err, "incorrect duration key representation")
	}

	return time.Duration(value), nil
}

/* Marshal float */
func SortedMarshalFloat(f float64) []byte {
	val := math.Float64bits(f)

	return SortedMarshalUint64(val, f >= 0.0)
}

func SortedUnmarshalFloat(b []byte) (float64, error) {
	value, err := SortedUnmarshalUint64(b)
	if err != nil {
		return 0.0, errors.Wrap(err, "incorrect float key representation")
	}

	return math.Float64frombits(value), nil
}