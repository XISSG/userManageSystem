package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"io"
	"log"
	"path"
	"strings"
	"time"
)

type Result struct {
	Logs     string
	CostTime time.Duration
	Memory   int64
}

const (
	imageName = "my-golang-image" //执行的镜像
	dstDir    = "/app"            //docker中执行的路径
)

func Docker(codePath string, input string) (Result, error) {
	var result Result

	//创建连接客户端
	ctx := context.Background()
	//初始化
	cli, err := initClient()
	if err != nil {
		log.Printf("client initialization error: %v", err)
		return Result{}, err
	}
	defer cli.Close()

	//初始化配置
	resp, err := initContainer(cli, codePath, input)
	if err != nil {
		log.Printf("container initialization error: %v", err)
		return Result{}, err
	}
	// 将文件编译后复制到容器中
	tarReader, err := archive.Tar(codePath, archive.Uncompressed)
	if err != nil {
		log.Fatal("compress file error:", err)
		return Result{}, err
	}
	err = cli.CopyToContainer(context.Background(), resp.ID, dstDir, tarReader, types.CopyToContainerOptions{})
	if err != nil {
		log.Fatal("copy file error", err)
		return Result{}, err
	}

	//开始执行容器
	if err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Println("container start:", err)
		return Result{}, err
	}

	//获取执行结果
	result.Logs, err = getLogs(resp, cli)
	result.CostTime, result.Memory, err = getStats(resp, cli)
	if err != nil {
		log.Println("get stats error", err)
		return Result{}, err
	}

	//删除容器
	err = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Println("container remove error", err)
		return Result{}, err
	}
	return result, nil
}

func initClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func initContainer(cli *client.Client, filePath string, input string) (container.CreateResponse, error) {
	//初始化配置
	filename := path.Base(filePath)
	file := path.Join(dstDir, filename)

	var strSlice []string
	args := input
	str := strings.Split(args, " ")
	strSlice = append(strSlice, file)
	strSlice = append(strSlice, str...)
	timeout := new(int)
	*timeout = 10
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        imageName,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		StopTimeout:  timeout,

		Cmd: strSlice,
	}, &container.HostConfig{ReadonlyRootfs: true}, nil, nil, "")
	if err != nil {
		return container.CreateResponse{}, err
	}

	return resp, nil
}

func getLogs(resp container.CreateResponse, cli *client.Client) (string, error) {

	// 容器 ID
	containerID := resp.ID

	// 获取容器日志
	logOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	logReader, err := cli.ContainerLogs(context.Background(), containerID, logOptions)
	if err != nil {
		panic(err)
	}
	defer logReader.Close()

	// 读取容器日志
	logs, err := io.ReadAll(logReader)
	if err != nil {
		return "", err
	}
	return string(logs), nil
}

func getStats(resp container.CreateResponse, cli *client.Client) (time.Duration, int64, error) {
	// 容器 ID
	containerID := resp.ID

	// 获取容器信息
	inspect, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return 0, 0, err
	}

	// 计算执行时间
	createdAt, err := time.Parse(time.RFC3339Nano, inspect.Created)
	if err != nil {
		return 0, 0, err
	}
	executionTime := time.Since(createdAt)

	// 获取内存使用量
	memoryUsage := inspect.HostConfig.Memory
	return executionTime, memoryUsage, nil
}
