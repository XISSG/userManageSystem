package judge

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/core/sanbox"
	"github.com/xissg/userManageSystem/entity/model_question"
	mysql2 "github.com/xissg/userManageSystem/service/mysql"
	"log"
)

type JudgeService struct {
	questionService       *mysql2.QuestionService
	questionSubmitService *mysql2.QuestionSubmitService
}

func NewJudgeService(questionService *mysql2.QuestionService, questionSubmitService *mysql2.QuestionSubmitService) *JudgeService {
	return &JudgeService{
		questionService:       questionService,
		questionSubmitService: questionSubmitService,
	}
}
func (s *JudgeService) Judge(submitId string) {
	//判断提交判题状态
	submit, err := s.questionSubmitService.GetSubmitQuestion(submitId)
	if err != nil {
		return
	}
	if submit.Status != constant.WAITING {
		return
	}
	res, err := s.questionService.GetQuestion(submit.QuestionId)
	if err != nil || res.ID == "" {
		return
	}

	var update model_question.UpdateQuestionSubmitRequest
	var judgeInfo model_question.JudgeInfo

	update.ID = submit.ID
	update.Status = constant.JUDGING
	judgeContext := sanbox.ToJudgeContext(&submit, &res)

	//开始沙箱判题
	box := sanbox.NewSanBox()
	result, err := box.Start(judgeContext)

	//沙箱初始化异常处理
	if err != nil {
		update.Status = constant.FAIL
		judgeInfo.Message = constant.CompileError
		update.JudgeInfo = append(update.JudgeInfo, judgeInfo)

		return
	}

	//程序执行内存溢出，超时等
	question := model_question.QuestionToReturnQuestion(res)

	for i := range result {
		switch result[i].ExitCode {
		case -1:
			update.Status = constant.FAIL
			judgeInfo.Message = constant.SystemError
			update.JudgeInfo = append(update.JudgeInfo, judgeInfo)
		case 0:
			if result[i].CostTime > question.JudgeConfig.TimeLimit {
				update.Status = constant.FAIL
				judgeInfo.Message = constant.TimeLimitExceeded

			} else if result[i].Memory > question.JudgeConfig.MemoryLimit {
				update.Status = constant.FAIL
				judgeInfo.Message = constant.MemoryLimitExceeded

			} else if result[i].ExecResult != question.Answer[i] {
				update.Status = constant.FAIL
				judgeInfo.Message = constant.WrongAnswer

			} else {
				update.Status = constant.SUCCESS
				judgeInfo.Message = constant.Accepted

			}
		default:
			update.Status = constant.FAIL
			judgeInfo.Message = constant.RuntimeError
		}
		judgeInfo.Time = result[i].CostTime
		judgeInfo.Memory = result[i].Memory
		update.JudgeInfo = append(update.JudgeInfo, judgeInfo)
	}

	//判题成功更新数据
	common := model_question.UpdateQSToCommonQS(update)
	err = s.questionSubmitService.UpdateSubmitQuestion(common)

	if err != nil {
		log.Printf("update submit question %v", err)
	}
}
