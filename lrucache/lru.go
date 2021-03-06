package lrucache

import (
	"fmt"
	"math"
	"sync"
)

// LRUキャッシュの定義
type LRUCache struct {
	limit   int
	values     map[int]*item
	currentAge int
	mutex      *sync.Mutex
}

// LRUキャッシュの生成
func NewLRU(limit int) (*LRUCache, error){
	if limit < 1 {
		return nil, fmt.Errorf("nonsensical LRU cache size specified\n")
	}

	return &LRUCache{
		limit: limit,
		values: make(map[int]*item, limit),
		currentAge: 0,
		mutex: new(sync.Mutex),
	}, nil
}

func (c *LRUCache)IsEmpty()bool{
	if len(c.values) == 0{
		return true
	}
	return false
}

// keyの値が存在していれば取り出してageをインクリメントする
func (c *LRUCache) Get(key int) int {
	i, ok := c.values[key]
	if !ok {
		return -1
	}
	c.mutex.Lock()
	i.age = c.currentAge
	c.currentAge++
	c.mutex.Unlock()
	return i.value
}

func (c *LRUCache) Put(key int, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	i, ok := c.values[key]
	// キーが存在する時は値を更新する
	if ok {
		i.value = value
		i.age = c.currentAge
		c.currentAge++
	}else {
		// limitを超えたときはageの低いkeyを探す
		if len(c.values) >= c.limit {
			leastAge := math.MaxInt32
			leastAgeKey := 0
			for key, item := range c.values {
				if item.age < leastAge {
					leastAge = item.age
					leastAgeKey = key
				}
			}
			// 最もageの低いキーを削除する
			if leastAgeKey != 0 {
				delete(c.values, leastAgeKey)
			}
		}
		c.values[key] = &item{
			value: value,
			age:   c.currentAge,
		}
		c.currentAge++
	}
}