// Code generated by "stringer -type RoundingMode"; DO NOT EDIT.

package decimal

import "fmt"

const _RoundingMode_name = "ToNearestEvenToNearestAwayToZeroAwayFromZeroToNegativeInfToPositiveInf"

var _RoundingMode_index = [...]uint8{0, 13, 26, 32, 44, 57, 70}

func (i RoundingMode) String() string {
	if i >= RoundingMode(len(_RoundingMode_index)-1) {
		return fmt.Sprintf("RoundingMode(%d)", i)
	}
	return _RoundingMode_name[_RoundingMode_index[i]:_RoundingMode_index[i+1]]
}
