package sanbox

import (
	"fmt"
	"github.com/xissg/userManageSystem/entity/model_question"
	"testing"
)

func TestPreProcess(t *testing.T) {
	ctx := &JudgeContext{
		ID:        "test",
		Language:  "Go",
		Code:      "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"strconv\"\n)\n\nfunc main() {\n\targs := os.Args[1:]\n\tsum := 0\n\tfor _, arg := range args {\n\t\tnum, _ := strconv.Atoi(arg)\n\t\tsum += num\n\t\t\t}\n\tfmt.Println(sum)\n}\n",
		JudgeCase: []model_question.JudgeCase{{Input: "1 2", Output: "3"}},
	}

	san := NewSanBox()
	results, err := san.Start(ctx)
	if err != nil {
		_ = fmt.Errorf("%v", err)
	}
	for _, test := range results {
		fmt.Printf("exit code: %v\n", test.ExitCode)
		fmt.Printf("exec result: %v\n", test.ExecResult)
		fmt.Printf("cost time: %v\n", test.CostTime)
		fmt.Printf("cost memory: %v\n", test.Memory)
	}
}
