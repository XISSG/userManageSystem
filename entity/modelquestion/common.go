package modelquestion

import "time"

// JudgeInfo Message 字段的值
const (
	Accepted            = "Accepted"
	WrongAnswer         = "Wrong Answer"
	CompileError        = "Compile Error"
	MemoryLimitExceeded = "Memory Limit Exceeded"
	TimeLimitExceeded   = "Time Limit Exceeded"
	PresentationError   = "Presentation Error"
	OutputLimitExceeded = "Output Limit Exceeded"
	Waiting             = "Waiting"
	DangerousOperation  = "Dangerous Operation"
	RuntimeError        = "Runtime Error"
	SystemError         = "System Error"
)

// status字段的值，用户提交题目状态
const (
	WAITING = 1
	JUDGING = 2
	SUCCESS = 3
	FAIL    = 4
)

// language字段的值
const (
	Java   = "Java"
	Python = "Python"
	C      = "C"
	Cpp    = "Cpp"
	Go     = "Go"
	Js     = "Js"
	Php    = " Php"
	Ts     = "Ts"
)

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
