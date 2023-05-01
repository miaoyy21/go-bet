package xmd

import (
	"sort"
	"strconv"
)

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
