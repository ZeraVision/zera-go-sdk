package convert

import (
	"encoding/json"

	"math/big"

	"github.com/jackc/pgtype"
	"gopkg.in/inf.v0"
)

// Helper function to convert various types to *big.Float
func ToBigFloat(val interface{}) *big.Float {
	switch v := val.(type) {
	case *inf.Dec:
		// Convert *inf.Dec to *big.Float
		bigFloatValue := new(big.Float)
		bigFloatValue.SetString(v.String())
		return bigFloatValue
	case *big.Int:
		return new(big.Float).SetInt(v)
	case int:
		return big.NewFloat(float64(v))
	case int16:
		return big.NewFloat(float64(v))
	case int32:
		return big.NewFloat(float64(v))
	case int64:
		return big.NewFloat(float64(v))
	case float64:
		return big.NewFloat(v)
	case *float64:
		if v == nil {
			return big.NewFloat(1.0 / 1e18)
		}

		return big.NewFloat(*v)
	case string:
		bigFloatValue := new(big.Float)
		_, success := bigFloatValue.SetString(v)
		if !success {
			return nil
		}
		return bigFloatValue
	case *string:
		if v == nil {
			return nil
		}
		bigFloatValue := new(big.Float)
		_, success := bigFloatValue.SetString(*v)
		if !success {
			return nil
		}
		return bigFloatValue
	default:
		return nil
	}
}

// Helper function to convert various types to *big.Int
func ToBigInt(val interface{}) *big.Int {
	switch v := val.(type) {
	case *inf.Dec:
		// Convert *inf.Dec to *big.Int
		bigIntValue := new(big.Int)
		bigIntValue.SetString(v.String(), 10)
		return bigIntValue
	case *big.Float:
		// Convert *big.Float to *big.Int
		bigIntValue := new(big.Int)
		v.Int(bigIntValue)
		return bigIntValue
	case int:
		return big.NewInt(int64(v))
	case int16:
		return big.NewInt(int64(v))
	case int32:
		return big.NewInt(int64(v))
	case int64:
		return big.NewInt(v)
	case string:
		bigIntValue := new(big.Int)
		bigIntValue, success := bigIntValue.SetString(v, 10)
		if !success {
			var bigIntNullable *big.Int
			return bigIntNullable
		}
		return bigIntValue
	case *string:
		if v == nil {
			return nil
		}
		bigIntValue := new(big.Int)
		bigIntValue, success := bigIntValue.SetString(*v, 10)
		if !success {
			return nil
		}
		return bigIntValue
	case pgtype.Numeric:
		if v.Int == nil {
			return nil
		}
		bigInt := new(big.Int).Set(v.Int)
		if v.Exp > 0 {
			bigInt.Mul(bigInt, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(v.Exp)), nil))
		} else if v.Exp < 0 {
			return nil // fractional numbers cannot be converted to big.Int
		}

		return bigInt
	case json.Number:
		bigIntValue := new(big.Int)
		bigIntValue, success := bigIntValue.SetString(string(v), 10)
		if !success {
			return nil
		}
		return bigIntValue
	default:
		return nil
	}
}

func ToInfDecimal(val interface{}) *inf.Dec {
	stringVal := ""

	switch v := val.(type) {
	case *string:
		if v == nil {
			return nil
		}
		stringVal = *v
	case string:
		stringVal = v
	case *big.Int:
		if v == nil {
			return nil
		}
		stringVal = v.String()
	case *inf.Dec:
		return v
	case *pgtype.Numeric:
		if v == nil || v.Status != pgtype.Present {
			return nil
		}
		dec := new(inf.Dec)
		if v.Int.IsInt64() {
			dec.SetUnscaled(v.Int.Int64())
		}
		dec.SetScale(inf.Scale(-v.Exp))
		return dec
	case pgtype.Numeric:
		if v.Status != pgtype.Present {
			return nil
		}
		dec := new(inf.Dec)
		if v.Int.IsInt64() {
			dec.SetUnscaled(v.Int.Int64())
		}
		dec.SetScale(inf.Scale(-v.Exp))
		return dec
	default:
		return nil
	}

	dec := new(inf.Dec)
	if _, ok := dec.SetString(stringVal); !ok {
		return nil
	}
	return dec
}

func InfDecToFloat64(d *inf.Dec) float64 {
	f := new(big.Float)
	f.SetPrec(64)
	f.SetString(d.String())

	// Convert big.Float to float64
	result, _ := f.Float64()
	return result
}
