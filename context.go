package decimal

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ericlagergren/decimal/internal/arith"
)

// Precision and scale limits.
const (
	MaxScale = math.MaxInt32 // largest allowed scale.
	MinScale = math.MinInt32 // smallest allowed scale.

	MaxPrecision = math.MaxInt32 // largest allowed Context precision.
	MinPrecision = 1             // smallest allowed Context precision.
)

// DefaultPrecision is the default precision used for decimals created as
// literals or using new.
const DefaultPrecision = 16

// Context is a per-decimal contextual object that governs specific operations
// such as how lossy operations (e.g. division) round.
type Context struct {
	// OperatingMode which dictates how the decimal operates under certain
	// conditions. See OperatingMode for more information.
	OperatingMode OperatingMode

	precision int32

	// Traps are a set of exceptional conditions that should result in an error.
	Traps Condition

	// Conditions are a set of the most recent exceptional conditions to occur
	// during an operation.
	Conditions Condition

	// Err is the most recent error to occur during an operation.
	Err error

	// RoundingMode instructs how an infinite (repeating) decimal expansion
	// (digits following the radix) should be rounded. This can occur during
	// "lossy" operations like division.
	RoundingMode RoundingMode
}

// New is shorthand to create a Big from a Context.
func (c Context) New(value int64, scale int32) *Big {
	x := New(value, scale)
	x.Context = c
	return x
}

const noPrecision = -1

// Precision returns the Context's precision.
func (c Context) Precision() int32 {
	switch c.precision {
	case 0:
		return DefaultPrecision
	case noPrecision:
		return 0
	default:
		return c.precision
	}
}

// SetPrecision sets c's precision.
func (c *Context) SetPrecision(prec int32) {
	if prec == 0 {
		c.precision = noPrecision
	} else {
		c.precision = prec
	}
}

// The following are called ContextXX instead of DecimalXX
// to reserve the DecimalXX namespace for future decimal types.

// The following Contexts are based on IEEE 754R. Each Context's RoundingMode is
// ToNearestEven, OperatingMode is GDA, and traps are set to every exception
// other than Inexact, Rounded, and Subnormal.
var (
	// Context32 is the IEEE 754R Decimal32 format.
	Context32 = Context{
		precision:     7,
		RoundingMode:  ToNearestEven,
		OperatingMode: GDA,
		Traps:         ^(Inexact | Rounded | Subnormal),
	}

	// Context64 is the IEEE 754R Decimal64 format.
	Context64 = Context{
		precision:     16,
		RoundingMode:  ToNearestEven,
		OperatingMode: GDA,
		Traps:         ^(Inexact | Rounded | Subnormal),
	}

	// Context128 is the IEEE 754R Decimal128 format.
	Context128 = Context{
		precision:     34,
		RoundingMode:  ToNearestEven,
		OperatingMode: GDA,
		Traps:         ^(Inexact | Rounded | Subnormal),
	}
)

// RoundingMode determines how a decimal will be rounded if the exact result
// cannot accurately be represented.
type RoundingMode uint8

// The following rounding modes are supported.
const (
	ToNearestEven RoundingMode = iota // == IEEE 754-2008 roundTiesToEven
	ToNearestAway                     // == IEEE 754-2008 roundTiesToAway
	ToZero                            // == IEEE 754-2008 roundTowardZero
	AwayFromZero                      // no IEEE 754-2008 equivalent
	ToNegativeInf                     // == IEEE 754-2008 roundTowardNegative
	ToPositiveInf                     // == IEEE 754-2008 roundTowardPositive
)

//go:generate stringer -type RoundingMode

// OperatingMode dictates how the decimal approaches specific non-numeric
// operations like conversions to strings and panicking on NaNs.
type OperatingMode uint8

const (
	// Go adheres to typical Go idioms. In particular:
	//
	//  - it panics on NaN values
	//  - has lossless (i.e., without rounding) addition, subtraction, and
	//    multiplication
	//  - traps are ignored; it does not set Context.Err or Context.Conditions
	//  - its string forms of qNaN, sNaN, +Inf, and -Inf are "NaN", "NaN",
	//     "+Inf", and "-Inf", respectively
	//
	Go OperatingMode = iota

	// GDA strictly adheres to the General Decimal Arithmetic Specification
	// Version 1.70. In particular:
	//
	//  - at does not panic
	//  - all arithmetic operations will be rounded down to the proper precision
	//    if necessary
	//  - it utilizes traps to set both Context.Err and Context.Conditions
	//  - its string forms of qNaN, sNaN, +Inf, and -Inf are "NaN", "sNaN",
	//    "Infinity", and "-Infinity", respectively
	//
	GDA
)

//go:generate stringer -type OperatingMode

// Condition is a bitmask value raised after or during specific operations. For
// example, dividing by zero is undefined so a DivisionByZero Condition flag
// will be set in the decimal's Context.
type Condition uint32

const (
	// Clamped occurs if the scale has been modified to fit the constraints of
	// the decimal representation.
	Clamped Condition = 1 << iota
	// ConversionSyntax occurs when a string is converted to a decimal and does
	// not have a valid syntax.
	ConversionSyntax
	// DivisionByZero occurs when division is attempted with a finite,
	// non-zero dividend and a divisor with a value of zero.
	DivisionByZero
	// DivisionImpossible occurs when the result of integer division would
	// contain too many digits (i.e. be longer than the specified precision).
	DivisionImpossible
	// DivisionUndefined occurs when division is attempted with in which both
	// the divided and divisor are zero.
	DivisionUndefined
	// Inexact occurs when the result of an operation (e.g. division) is not
	// exact, or when the Overflow/Underflow Conditions occur.
	Inexact
	// InsufficientStorage occurs when the system doesn't have enough storage
	// (i.e. memory) to store the decimal.
	InsufficientStorage
	// InvalidContext occurs when an invalid context was detected during an
	// operation. This might occur if, for example, an invalid RoundingMode was
	// passed to a Context.
	InvalidContext
	// InvalidOperation occurs when:
	//
	// 	- an operand to an operation is a signaling NaN
	// 	- an attempt is made to add or subtract infinities of opposite signs
	// 	- an attempt is made to multiply zero by an infinity of either sign
	// 	- an attempt is made to divide an infinity by an infinity
	// 	- the divisor for a remainder operation is zero
	// 	- the dividend for a remainder operation is an infinity
	// 	- either operand of the quantize operation is an infinity, or the result
	// 	  of a quantize operation would require greater precision than is
	// 	  available
	// 	- the operand of the ln or the log10 operation is less than zero
	// 	- the operand of the square-root operation has a sign of 1 and a
	// 	  non-zero coefficient
	// 	- both operands of the power operation are zero, or if the left-hand
	// 	  operand is less than zero and the right-hand operand does not have an
	// 	  integral value or is an infinity
	//
	InvalidOperation
	// Overflow occurs when the adjusted scale, after rounding, would be
	// greater than MaxScale. (Inexact and Rounded will also be raised.)
	Overflow
	// Rounded occurs when the result of an operation is rounded, or if an
	// Overflow/Underflow occurs.
	Rounded
	// Subnormal ocurs when the result of a conversion or operation is subnormal
	// (i.e. the adjusted scale is less than MinScale before any rounding).
	Subnormal
	// Underflow occurs when the result is inexact and the adjusted scale would
	// be smaller (more negative) than MinScale.
	Underflow
)

func (c Condition) String() string {
	if c == 0 {
		return ""
	}

	// Each condition is one bit, so this saves some allocations.
	a := make([]string, 0, arith.Popcnt32(uint32(c)))
	for i := Condition(1); c != 0; i <<= 1 {
		if c&i == 0 {
			continue
		}
		switch c ^= i; i {
		case Clamped:
			a = append(a, "clamped")
		case ConversionSyntax:
			a = append(a, "conversion syntax")
		case DivisionByZero:
			a = append(a, "division by zero")
		case DivisionImpossible:
			a = append(a, "division impossible")
		case Inexact:
			a = append(a, "inexact")
		case InsufficientStorage:
			a = append(a, "insufficient storage")
		case InvalidContext:
			a = append(a, "invalid context")
		case InvalidOperation:
			a = append(a, "invalid operation")
		case Overflow:
			a = append(a, "overflow")
		case Rounded:
			a = append(a, "rounded")
		case Subnormal:
			a = append(a, "subnormal")
		case Underflow:
			a = append(a, "underflow")
		default:
			a = append(a, fmt.Sprintf("unknown(%d)", i))
		}
	}
	return strings.Join(a, ", ")
}

var (
	one = New(1, 0)

	zeroInt = big.NewInt(0)
	oneInt  = big.NewInt(1)
	twoInt  = big.NewInt(2)
	tenInt  = big.NewInt(10)

	tenFloat = big.NewFloat(10)
)
