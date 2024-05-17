package docker

import (
	"fmt"
	"github.com/xissg/userManageSystem/core/sanbox"
	"github.com/xissg/userManageSystem/entity/model_question"
	"testing"
)

func TestDocker(t *testing.T) {
	ctx := &sanbox.JudgeContext{
		ID:        "test",
		Language:  "Go",
		Code:      "package main\nimport (\n\"fmt\"\n\"os\"\n)\nfunc main() {\nargs := os.Args\nfor _, arg := range args {\n fmt.Println(arg+\"\\n\")\n }\npanic(\"error\")\n}",
		JudgeCase: []model_question.JudgeCase{{Input: "a b", Output: "a b"}, {Input: "c d", Output: "c d"}},
	}

	docker, err := Docker("", ctx.JudgeCase[0].Input)
	if err != nil {
		return
	}
	fmt.Println(docker)
}
