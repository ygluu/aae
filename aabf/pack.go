package aabf

import (
	"errors"
	"fmt"
)

func unpack(data []byte, on_pack func(data []byte)) (ret []byte, err error) {
	index := 0
	sum := len(data)

	for index < sum {
		curr := index
		if sum-curr < 3 {
			break
		}
		plen := (int)(((uint32)(data[index])) + ((uint32)(data[index+1]))<<8 + ((uint32)(data[index+2]))<<16)
		curr += 3
		if sum-curr < plen+3 {
			break
		}

		iFlag := index + plen
		flag := (((uint32)(data[iFlag])) + ((uint32)(data[iFlag+1]))<<8 + ((uint32)(data[iFlag+2]))<<16)
		if (uint32)(plen) != ^flag {
			return ret, errors.New(fmt.Sprintf("Serious flag error >> %d-%d", plen, flag))
		}

		on_pack(data[index:plen])
		index += 6 + plen
	}

	if index < sum {
		ret = data[index : sum-index]
	}

	return ret, nil
}
