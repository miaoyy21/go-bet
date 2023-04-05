package xmd

// 每个数字1000次出现的标准次数
var stds map[int]int

func calc() {
	stds = make(map[int]int)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				stds[i+j+k]++
			}
		}
	}
}
