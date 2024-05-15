package sanbox

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/xissg/userManageSystem/core/docker"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type JudgeContext struct {
	//题目id
	ID string `json:"id"`
	// "编程语言"
	Language string `json:"language" `
	//"用户代码"
	Code string `json:"code" `
	// "判题用例json数组"
	JudgeCase []modelquestion.JudgeCase `json:"judge_case" `
}

func ToJudgeContext(submit *modelquestion.QuestionSubmit, question *modelquestion.Question) *JudgeContext {
	var judgeCase []modelquestion.JudgeCase
	err := json.Unmarshal([]byte(question.JudgeCase), &judgeCase)
	if err != nil {
		return nil
	}

	return &JudgeContext{
		ID:        question.ID,
		Language:  submit.Language,
		Code:      submit.Code,
		JudgeCase: judgeCase,
	}
}

type JudgeResult []docker.Result

// 代码逻辑，接收数据，将对象中code字段保存到文件，读取文件并编译，将编译后的文件复制到docker中进行运行，返回运行结果
type SanBox struct {
	filePath string
	codePath string
}

func NewSanBox() *SanBox {
	return &SanBox{}
}

func (s *SanBox) Start(ctx *JudgeContext) (JudgeResult, error) {
	//数据预处理
	err := s.preProcess(ctx)
	if err != nil {
		return nil, err
	}

	//开始编译
	err = s.compile()
	if err != nil {
		return nil, err
	}

	//开始运行代码
	res := s.run(ctx)

	//删除文件
	err = s.deleteFile()
	log.Printf("delete file: %v", err)

	//处理返回结果
	return res, nil
}

// 预处理，对代码进行安全校验，根据题目id分类保存文件，日志记录，错误处理等
func (s *SanBox) preProcess(ctx *JudgeContext) error {
	//校验数据
	err := s.checkData(ctx)
	if err != nil {
		log.Println("data validate error:", err)
		return err
	}
	//创建文件夹
	folderPath := s.mkdir(ctx)
	if folderPath == "" {
		log.Println("directory create error")
		return errors.New("create directory error")
	}

	//创建文件并写入数据
	err = s.touchAndWrite(folderPath, ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// 编译并保存可执行文件
func (s *SanBox) compile() error {
	if s.filePath == "" || s.codePath == "" {
		return errors.New("invalid file or code path")
	}

	//执行编译命令
	cmd := exec.Command("go", "build", "-o", s.codePath, s.filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (s *SanBox) run(ctx *JudgeContext) JudgeResult {

	if s.codePath == "" || ctx.JudgeCase == nil {
		return nil
	}
	err := os.Chmod(s.codePath, 0755)
	if err != nil {
		return nil
	}

	var results JudgeResult
	for _, v := range ctx.JudgeCase {
		//多次执行结果
		result, err := docker.Docker(s.codePath, v.Input)
		if err != nil {
			return nil
		}
		results = append(results, result)
	}

	return results
}

// 数据校验
func (s *SanBox) checkData(ctx *JudgeContext) error {
	if ctx.ID == "" {
		return errors.New("invalid submit id")
	}
	if ctx.Code == "" {
		return errors.New("invalid code")
	}
	if ctx.JudgeCase == nil {
		return errors.New("invalid judge case")
	}
	return nil
}

func (s *SanBox) mkdir(ctx *JudgeContext) string {
	// 指定文件夹路径
	path, err := os.Getwd()
	parentDir := filepath.Dir(path)
	folderPath := filepath.Join(parentDir, "tmp", ctx.ID)
	_, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		// 创建文件夹
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return ""
		}
	}
	return folderPath
}

// 创建并写入代码
func (s *SanBox) touchAndWrite(folderPath string, ctx *JudgeContext) error {
	//生成文件名
	id := strings.Split(uuid.NewString(), "-")[0]
	fileName := fmt.Sprintf("%s.%s", id, strings.ToLower(ctx.Language))

	//保存源文件,以及以后生成可执行文件名
	s.filePath = filepath.Join(folderPath, fileName)
	s.codePath = filepath.Join(folderPath, id)
	// 创建文件
	file, err := os.Create(s.filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	//读取并写入文件
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(ctx.Code)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer:", err)
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	return nil
}

// 删除临时文件
func (s *SanBox) deleteFile() error {
	err := os.Remove(s.filePath)
	err = os.Remove(s.codePath)
	if err != nil {
		return err
	}

	return nil
}
