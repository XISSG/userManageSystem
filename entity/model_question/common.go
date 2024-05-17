package model_question

type JudgeCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type JudgeConfig struct {
	TimeLimit   int64  `json:"time_limit"`
	MemoryLimit uint64 `json:"memory_limit"`
}

type JudgeInfo struct {
	Message string `json:"message"` //值为以上枚举值
	Time    int64  `json:"time"`    //单位为ms
	Memory  uint64 `json:"memory"`  //单位为kb
}
