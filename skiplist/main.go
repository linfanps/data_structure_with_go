package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	SKIPLIST_MAX_LEVEL = 10
	SKIPLIST_P         = 0.5
)

type SkipListNode struct {
	key   string
	score int
	next  []*SkipListNode
}

type SkipList struct {
	head  *SkipListNode
	level int
}

func (sl *SkipList) Insert(key string, score int) int {
	// 每层需要更新的结点
	updateNodes := make([]*SkipListNode, SKIPLIST_MAX_LEVEL)
	node := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].score < score {
			node = node.next[i]
		}
		updateNodes[i] = node
	}

	newNodeLevel := randomLevel()
	if newNodeLevel > sl.level {
		for i := sl.level; i < newNodeLevel; i++ {
			updateNodes[i] = sl.head
		}
		sl.level = newNodeLevel
	}
	newNode := createNode(newNodeLevel, score, key)

	for i := 0; i < newNodeLevel; i++ {
		newNode.next[i] = updateNodes[i].next[i]
		updateNodes[i].next[i] = newNode
	}

	return newNodeLevel
}

func (sl *SkipList) Search(score int) *SkipListNode {
	node := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].score <= score {
			node = node.next[i]
		}
	}
	if node != nil && node.score == score {
		return node
	}
	return nil
}

func (sl *SkipList) Delete(score int) *SkipListNode {
	node := sl.head
	updateNodes := make([]*SkipListNode, SKIPLIST_MAX_LEVEL)

	for i := sl.level - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].score < score {
			node = node.next[i]
		}
		updateNodes[i] = node
	}

	node = node.next[0]
	if node != nil && node.score == score {
		for i := 0; i < sl.level-1; i++ {
			if updateNodes[i].next[i] == node {
				updateNodes[i].next[i] = node.next[i]
			}
		}
	}
	return node
}

func (sl *SkipList) Show() {
	head := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		fmt.Printf("level %d:", i)
		for node := head; node != nil; node = node.next[i] {
			if node == head {
				fmt.Print("head->")
			} else {
				fmt.Printf("%s(%d)->", node.key, node.score)
			}
		}
		fmt.Print("nil\n")
	}
}

func randomLevel() int {
	level := 1
	for rand.Float64() < SKIPLIST_P && level < SKIPLIST_MAX_LEVEL {
		level += 1
	}
	return level
}

func createNode(level, score int, key string) *SkipListNode {
	sln := &SkipListNode{
		key:   key,
		score: score,
		next:  make([]*SkipListNode, level),
	}

	return sln
}

func createSkipList() *SkipList {
	sl := &SkipList{
		head:  createNode(SKIPLIST_MAX_LEVEL, 0, ""),
		level: 0,
	}
	return sl
}

func main() {
	rand.Seed(time.Now().Unix())

	skiplist := createSkipList()
	v := map[string]int{
		"one":   1,
		"two":   2,
		"four":  4,
		"five":  5,
		"seven": 7,
		"eight": 8,
		"nine":  9,
		"ten":   10,
	}

	for key, score := range v {
		level := skiplist.Insert(key, score)
		fmt.Printf("Insert %s(%d) Level:%d\n", key, score, level)
		skiplist.Show()
	}

	for score := 1; score <= 11; score++ {
		node := skiplist.Search(score)
		if node == nil {
			fmt.Printf("Search %d Not Found\n", score)
		} else {
			fmt.Printf("Search %d Found Key:%s\n", score, node.key)
		}
	}

	deteleScores := []int{2, 22, 8}
	for _, score := range deteleScores {
		skiplist.Delete(score)
		fmt.Printf("Delete %d\n", score)
		skiplist.Show()
	}
}
