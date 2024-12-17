package uuids

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var entropySize = new(big.Int).Exp(big.NewInt(2), big.NewInt(122), nil)
var slope *big.Int = new(big.Int).Exp(big.NewInt(3), big.NewInt(77), nil)
var mod_inv_slope *big.Int

func init() {
	var err error
	mod_inv_slope, err = modularInverse(slope, entropySize)
	if err != nil {
		panic(err)
	}
}

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
		id.shuffle()
		ret[i] = id.toString()
	}
	return ret
}

func (u *uuid) shuffle() {
	u.number = linearEncode(slope, u.number, big.NewInt(30), entropySize)
	u.number = linearEncode(slope, u.number, big.NewInt(0), entropySize)
}

func (u *uuid) unshuffle() {
	u.number = linearDecode(mod_inv_slope, u.number, big.NewInt(0), entropySize)
	u.number = linearDecode(mod_inv_slope, u.number, big.NewInt(30), entropySize)
}

func linearEncode(a, x, b, N *big.Int) *big.Int {
	ret := new(big.Int).Add(new(big.Int).Mul(a, x), b)
	ret.Mod(ret, N)
	return ret
}

func linearDecode(inv_a, x, b, N *big.Int) *big.Int {
	ret := new(big.Int).Sub(x, b)
	ret.Mul(ret, inv_a)
	ret.Mod(ret, N)
	return ret
}

func modularInverse(a, N *big.Int) (*big.Int, error) {
	// Variables to store results of GCD computation
	gcd := new(big.Int)
	x := new(big.Int)

	// Compute GCD and the coefficients x and y
	gcd.GCD(x, nil, a, N)

	// Check if GCD is 1 (a and N must be coprime)
	if gcd.Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("no modular inverse exists (a and N are not coprime)")
	}

	// Ensure the result (x) is positive
	if x.Sign() < 0 {
		x.Add(x, N)
		if x.Sign() < 0 {
			return nil, errors.New("???, x still negative?")
		}
	}

	return x, nil
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
	return new(big.Int).Sub(entropySize, big.NewInt(1))
}

func (u *uuid) sliceBits(start, length uint) (ret *big.Int) {
	ret = new(big.Int).Lsh(u.number, start)
	ret = ret.And(ret, MaxEntropyNum())
	ret.Rsh(ret, 122-length)
	return
}
