package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//哈希将字节映射到uint32
//定义了函数类型 Hash，采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换，
//默认为 crc32.ChecksumIEEE 算法。
type Hash func(data []byte) uint32

//Map constains all hashed keys
//Map 是一致性哈希算法的主数据结构，包含 4 个成员变量：
//Hash 函数 hash；
//虚拟节点倍数 replicas；
//哈希环 keys；
//虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

//新建创建一个Map实例,构造函数 New() 允许自定义虚拟节点倍数和 Hash 函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//Add 函数允许传入 0 或 多个真实节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		//对每一个真实节点 key，对应创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			//虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
			//使用 m.hash() 计算虚拟节点的哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//使用 append(m.keys, hash) 添加到环上。
			m.keys = append(m.keys, hash)
			//在 hashMap 中增加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	//最后一步，环上的哈希值排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//第一步，计算 key 的哈希值。
	hash := int(m.hash([]byte(key)))
	//第二步，顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。
	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//第三步，通过 hashMap 映射得到真实的节点。
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
