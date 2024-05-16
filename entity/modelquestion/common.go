package modelquestion

import "time"

type Tags struct {
	Easy   string `json:"easy"`
	Medium string `json:"medium"`
	Hard   string `json:"hard"`
}

type JudgeCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type JudgeConfig struct {
	TimeLimit   time.Duration `json:"time_limit"`
	MemoryLimit int64         `json:"memory_limit"`
	StackLimit  int64         `json:"stack_limit"`
}

type JudgeInfo struct {
	Message string        `json:"message"` //值为以上枚举值
	Time    time.Duration `json:"time"`    //单位为ms
	Memory  int64         `json:"memory"`  //单位为kb
}
