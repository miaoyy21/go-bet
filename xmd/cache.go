package xmd

import (
	"sort"
	"strconv"
)

type UserBase struct {
	isDebug bool

	gold   int
	unix   string
	code   string
	device string
	id     string
	token  string
}

func NewUserBase(debug bool, gold int, unix string, code string, device string, id string, token string) UserBase {
	return UserBase{
		isDebug: debug,

		gold:   gold,
		unix:   unix,
		code:   code,
		device: device,
		id:     id,
		token:  token,
	}
}

type IssueResult struct {
	issue  int
	result int
}

type Cache struct {
	user UserBase

	issue  int // 最新期数
	result int // 最新开奖结果

	histories []IssueResult // 每期存在数据库的开奖记录
}

func NewCache(user UserBase) (*Cache, error) {
	cache := &Cache{
		user: user,

		issue:  -1,
		result: -1,

		histories: make([]IssueResult, 0),
	}

	return cache, nil
}

func (o *Cache) Sync(size int) error {
	items, err := hGetHistories(size, o.user)
	if err != nil {
		return err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Issue <= items[j].Issue
	})

	histories := make([]IssueResult, 0, len(items))
	for _, item := range items {
		issue, err := strconv.Atoi(item.Issue)
		if err != nil {
			return err
		}

		result, err := strconv.Atoi(item.Result)
		if err != nil {
			return err
		}

		o.issue = issue
		o.result = result

		histories = append(histories, IssueResult{issue: issue, result: result})
	}
	o.histories = histories

	return nil
}
