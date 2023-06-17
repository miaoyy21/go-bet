package xmd

import (
	"strconv"
	"strings"
)

func (o *Cache) Sync(size int) error {
	items, err := hAnalyseHistory(size, o.user)
	if err != nil {
		return err
	}

	histories := make([]IssueResult, 0, len(items))
	for i, item := range items {
		issue, err := strconv.Atoi(item.Issue)
		if err != nil {
			return err
		}

		result, err := strconv.Atoi(item.Result)
		if err != nil {
			return err
		}

		money, err := strconv.Atoi(strings.ReplaceAll(item.Money, ",", ""))
		if err != nil {
			return err
		}

		member := item.Member

		// 最新结果
		if i == 0 {
			o.issue = issue
			o.result = result
			o.money = money
			o.member = member
		}

		histories = append(histories,
			IssueResult{
				issue:  issue,
				result: result,
				money:  money,
				member: member,
			},
		)
	}
	o.histories = histories

	return nil
}
