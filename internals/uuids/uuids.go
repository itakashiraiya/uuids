package uuids

import (
	"fmt"
	"math/big"
	"strings"
)

type uuid struct {
	number *big.Int
}

func Test() {
	a := new(uuid)
	a.number = MaxEntropyNum()
	a.number.Xor(a.number, new(big.Int).Mul(big.NewInt(999999324189000006), big.NewInt(6983245123485468)))
	fmt.Printf("hex: %x, num: %d, bin: %08b\n", a.number, a.number, a.number)
	fmt.Println(a.toString())
	fmt.Printf("kex: %x\n", toUuid(a.toString()).number)
	b := a.sliceBits(7, 32)
	fmt.Printf("hex: %x, num: %d, bin: %08b\n", b, b, b)
}

func GetUuids(start *big.Int, length int) []string {
	ret := make([]string, length)
	for i := range ret {
		id := uuid{new(big.Int).Add(start, big.NewInt(int64(i)))}
		ret[i] = id.toString()
	}
	return ret
}

func (u *uuid) toString() string {
	ret := ""
	ret += fmt.Sprintf("%08s", u.sliceBits(0, 32).Text(16))
	ret += "-"
	ret += fmt.Sprintf("%04s", u.sliceBits(32, 16).Text(16))
	ret += "-4"
	ret += fmt.Sprintf("%03s", u.sliceBits(48, 12).Text(16))
	ret += "-"
	ret += fmt.Sprintf("%04s", new(big.Int).Add(new(big.Int).Lsh(big.NewInt(2), 14), u.sliceBits(60, 14)).Text(16))
	ret += "-"
	ret += fmt.Sprintf("%012s", u.sliceBits(74, 48).Text(16))
	return ret
}

func toUuid(sUuid string) uuid {
	parse := func(subString string) *big.Int {
		num, ok := new(big.Int).SetString(subString, 16)
		if !ok {
			panic("aaa")
		}
		return num
	}

	list := strings.Split(sUuid, "-")
	ret := uuid{number: big.NewInt(0)}

	ret.number.Add(ret.number, new(big.Int).Lsh(parse(list[0]), 90))
	ret.number.Add(ret.number, new(big.Int).Lsh(parse(list[1]), 74))
	ret.number.Add(ret.number, new(big.Int).Lsh(new(big.Int).Sub(parse(list[2]), new(big.Int).Lsh(big.NewInt(4), 12)), 62))
	ret.number.Add(ret.number, new(big.Int).Lsh(new(big.Int).Sub(parse(list[3]), new(big.Int).Lsh(big.NewInt(2), 14)), 48))
	ret.number.Add(ret.number, new(big.Int).Lsh(parse(list[4]), 0))

	return ret
}

func MaxEntropyNum() *big.Int {
	return new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(122), nil), big.NewInt(1))
}

func (u *uuid) sliceBits(start, length uint) (ret *big.Int) {
	ret = new(big.Int).Lsh(u.number, start)
	ret = ret.And(ret, MaxEntropyNum())
	ret.Rsh(ret, 122-length)
	return
}
