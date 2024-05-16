package constant

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
	Java   = "java"
	Python = "python"
	C      = "c"
	Cpp    = "cpp"
	Go     = "go"
	Js     = "js"
	Php    = " php"
	Ts     = "ts"
)