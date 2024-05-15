package judge

import (
	"github.com/xissg/userManageSystem/core/sanbox"
	"github.com/xissg/userManageSystem/entity/modelquestion"
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
	//TODO:判断执行结果和答案是否一致
	//判断提交判题状态
	submit, err := s.questionSubmitService.GetSubmitQuestion(submitId)
	if err != nil {
		return
	}
	if submit.Status != modelquestion.WAITING {
		return
	}
	res, err := s.questionService.GetQuestion(submit.QuestionId)
	if err != nil || res.ID == "" {
		return
	}

	var update modelquestion.UpdateQuestionSubmitRequest
	update.Status = modelquestion.JUDGING
	judgeContext := sanbox.ToJudgeContext(&submit, &res)

	//开始沙箱判题
	box := sanbox.NewSanBox()
	result, err := box.Start(judgeContext)

	//判题结果处理
	if err != nil {
		update.JudgeInfo[0].Message = err.Error()
		update.Status = modelquestion.FAIL
	}

	//程序执行内存溢出，超时等
	question := modelquestion.QuestionToReturnQuestion(res)
	for i := range result {
		if result[i].CostTime > question.JudgeConfig.TimeLimit {
			update.Status = modelquestion.FAIL
			update.JudgeInfo[i].Message = modelquestion.TimeLimitExceeded

		} else if result[i].Memory > question.JudgeConfig.MemoryLimit {
			update.Status = modelquestion.FAIL
			update.JudgeInfo[i].Message = modelquestion.MemoryLimitExceeded

		} else {
			update.Status = modelquestion.SUCCESS
			update.JudgeInfo[i].Message = modelquestion.Accepted
		}
		update.JudgeInfo[i].Time = result[i].CostTime
		update.JudgeInfo[i].Memory = result[i].Memory
	}

	//判题成功更新数据
	common := modelquestion.UpdateQSToCommonQS(update)
	err = s.questionSubmitService.UpdateSubmitQuestion(common)

	if err != nil {
		log.Printf("update submit question %v", err)
	}
}
