// Code generated by "stringer -type OperatingMode"; DO NOT EDIT.

package decimal

import "fmt"

const _OperatingMode_name = "GoGDA"

var _OperatingMode_index = [...]uint8{0, 2, 5}

func (i OperatingMode) String() string {
	if i >= OperatingMode(len(_OperatingMode_index)-1) {
		return fmt.Sprintf("OperatingMode(%d)", i)
	}
	return _OperatingMode_name[_OperatingMode_index[i]:_OperatingMode_index[i+1]]
}
