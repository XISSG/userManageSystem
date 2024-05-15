package sanbox

import (
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"testing"
)

func TestPreProcess(t *testing.T) {
	ctx := &JudgeContext{
		ID:        "test",
		Language:  "Go",
		Code:      "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"strconv\"\n)\n\nfunc main() {\n\targs := os.Args[1:]\n\tsum := 0\n\tfor _, arg := range args {\n\t\tnum, _ := strconv.Atoi(arg)\n\t\tsum += num\n\t\t\t}\n\tfmt.Println(sum)\n}\n",
		JudgeCase: []modelquestion.JudgeCase{{Input: "1 2", Output: "3"}},
	}

	san := NewSanBox()
	err := san.preProcess(ctx)
	if err != nil {
		t.Error(err)
	}
	err = san.compile()
	if err != nil {
		t.Error(err)
	}
	san.run(ctx)

}
