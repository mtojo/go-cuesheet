// generated by stringer -type Flags cuesheet/cuesheet.go; DO NOT EDIT

package cuesheet

import "fmt"

const _Flags_name = "NoneDcpFour_chPreScms"

var _Flags_index = [...]uint8{0, 4, 7, 14, 17, 21}

func (i Flags) String() string {
	if i < 0 || i+1 >= Flags(len(_Flags_index)) {
		return fmt.Sprintf("Flags(%d)", i)
	}
	return _Flags_name[_Flags_index[i]:_Flags_index[i+1]]
}
