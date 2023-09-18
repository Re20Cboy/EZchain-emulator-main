package core

///////////////////////////
/////// by XLD-9-18 ///////
///////////////////////////

type BitMap struct {
	Bits []byte
	Vmax uint
}

func NewBitMap(max_val ...uint) *BitMap {
	var max uint = 8192
	if len(max_val) > 0 && max_val[0] > 0 {
		max = max_val[0]
	}

	bm := &BitMap{}
	bm.Vmax = max
	sz := (max + 7) / 8
	bm.Bits = make([]byte, sz, sz)
	return bm
}

func (bm *BitMap)Set(num uint) {
	if num > bm.Vmax {
		bm.Vmax += 1024
		if bm.Vmax < num {
			bm.Vmax = num
		}

		dd := int(num+7)/8 - len(bm.Bits)
		if dd > 0 {
			tmp_arr := make([]byte, dd, dd)
			bm.Bits = append(bm.Bits, tmp_arr...)
		}
	}

	//将1左移num%8后，然后和以前的数据做|，这样就替换成1了
	bm.Bits[num/8] |= 1 << (num%8)
}

func (bm *BitMap)UnSet(num uint) {
	if num > bm.Vmax {
		return
	}
	//&^:将1左移num%8后，然后进行与非运算，将运算符左边数据相异的位保留，相同位清零
	bm.Bits[num/8] &^= 1 << (num%8)
}

func (bm *BitMap)Check(num uint) bool {
	if num > bm.Vmax {
		return false
	}
	//&:与运算符，两个都是1，结果为1
	return bm.Bits[num/8] & (1 << (num%8)) != 0
}


type BloomFilter struct {
	Bset *BitMap
	Size uint
}

func NewBloomFilter(size_val ...uint) *BloomFilter { //bloom的大小缺省值为1024*1024
	var size uint = 1024*1024
	if len(size_val) > 0 && size_val[0] > 0 {
		size = size_val[0]
	}

	bf := &BloomFilter{}
	bf.Bset = NewBitMap(size)
	bf.Size = size
	return bf
}

//hash函数
var seeds = []uint{3011, 3017, 3031}
func (bf *BloomFilter)hashFun(seed uint, value string) uint64 {
	hash := uint64(seed)
	for i := 0; i < len(value); i++ {
		hash = hash*33 + uint64(value[i])
	}
	return hash
}

//添加元素
func (bf *BloomFilter)Set(value string) {
	for _, seed := range seeds {
		hash := bf.hashFun(seed, value)
		hash = hash % uint64(bf.Size)
		bf.Bset.Set(uint(hash))
	}
}

//判断元素是否存在
func (bf *BloomFilter)Check(value string) bool {
	for _, seed := range seeds {
		hash := bf.hashFun(seed, value)
		hash = hash % uint64(bf.Size)
		ret := bf.Bset.Check(uint(hash))
		if !ret {
			return false
		}
	}
	return true
}