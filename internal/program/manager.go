package program

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getProgramPath 获取当前目录下指定程序的绝对路径。
func getProgramPath(path string) (string, error) {
	//fmt.Printf("获取程序路径: %s", path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("获取程序路径时出错: %v\n", err)
		return "", err
	}
	//fmt.Printf("程序绝对路径: %s", absPath)
	return absPath, nil
}

// isRunning 检查指定的程序是否正在运行
// 通过调用 pgrep 命令来实现检查
func isRunning(program string) (bool, error) {
	/*programPath, err := getProgramPath(program)
	if err != nil {
		return false, err
	}*/

	//fmt.Printf("检查程序是否运行: %s", program)

	// 使用 pgrep 命令来查找正在运行的程序
	cmd := exec.Command("pgrep", "-fl", program)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		// 检查 pgrep 的退出状态码，如果是 1 表示没有找到进程，不算错误
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			//fmt.Printf("pgrep 未找到匹配进程: %s", program)
			return false, nil
		}
		fmt.Printf("运行 pgrep 时出错: %v\n", err)
		return false, err
	}

	//fmt.Printf("pgrep 输出:\n%s", out.String())

	// 获取当前程序的 PID
	currentPID := os.Getpid()
	//fmt.Printf("当前程序 PID: %d", currentPID)

	// 检查 pgrep 输出，过滤掉当前运行的命令
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		// 跳过空行
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 将每行按空格分割成字段
		fields := strings.Fields(line)
		// 如果字段数小于2，跳过
		if len(fields) < 2 {
			continue
		}

		// 获取进程ID (PID) 和命令
		pid := fields[0]
		cmd := strings.Join(fields[1:], " ")

		// 过滤掉当前运行的命令
		if pid == string(rune(currentPID)) || (strings.Contains(cmd, "ecctl") && strings.Contains(cmd, "status")) {
			continue
		}

		// 检查命令行是否包含程序名
		if strings.Contains(cmd, program) {
			//fmt.Printf("找到匹配进程: %s", line)
			return true, nil
		}
	}

	//fmt.Printf("未找到匹配进程: %s", program)
	return false, nil
}

// Start 启动指定的程序
func Start(path string) {
	//fmt.Printf("尝试启动程序: %s", path)

	// 获取程序的绝对路径。
	absPath, err := getProgramPath(path)
	if err != nil {
		log.Fatalln("获取程序路径时出错：", err)
		return
	}

	// 检查程序是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("程序不存在：%s\n", absPath)
		return
	}

	// 获取程序的工作目录
	workDir := filepath.Dir(absPath)
	//fmt.Printf("目录: %s", workDir)

	program := filepath.Base(absPath)
	// 创建一个新的命令来运行指定的程序
	var cmd *exec.Cmd
	if strings.HasSuffix(program, ".jar") {
		cmd = exec.Command("nohup java -jar " + program + " &")
	} else {
		cmd = exec.Command("./" + program)
	}

	// 检查程序是否已经在运行。
	running, err := isRunning(program)
	if err != nil {
		log.Fatalln("检查程序运行状态时出错：", err)
		return
	}

	if running {
		fmt.Printf("%s is already running.\n", program)
		return
	}

	// 设置命令的工作目录。
	cmd.Dir = workDir

	// 启动命令执行的进程
	err = cmd.Start()
	if err != nil {
		log.Fatalln("启动程序时出错：", err)
		return
	}

	fmt.Printf("%s started successfully.\n", program)

	// 启动一个新的 goroutine 来等待进程的结束
	go func() {
		err := cmd.Wait()
		if err != nil {
			return
		}
		fmt.Printf("%s has been launched.\n", program)
	}()
}

// Stop 停止指定的程序
// 使用 pkill 命令来终止运行中的程序
func Stop(program string) {
	//fmt.Printf("尝试停止程序: %s", program)

	// 获取程序的绝对路径。
	/*programPath, err := getProgramPath(program)
	if err != nil {
		logger.Log.Error("获取程序路径时出错：", err)
		return
	}*/

	// 检查程序是否已经在运行。
	running, err := isRunning(program)
	if err != nil {
		log.Fatalln("检查程序运行状态时出错：", err)
		return
	}

	if !running {
		fmt.Printf("%s is not running.\n", program)
		return
	}

	// 使用 pkill 命令终止运行中的程序
	cmd := exec.Command("pkill", "-f", program)
	err = cmd.Run()
	if err != nil {
		log.Fatalln("停止程序时出错：", err)
		return
	}

	// 记录程序停止成功的信息。
	fmt.Printf("%s has stoped.\n", program)
}

// Status 检查指定程序的运行状态
// 使用 pgrep 命令来检查程序是否在运行
func Status(program string) {
	// 获取程序的绝对路径。
	/*programPath, err := getProgramPath(program)
	if err != nil {
		logger.Log.Error("获取程序路径时出错：", err)
		return
	}*/

	// 检查程序是否正在运行
	running, err := isRunning(program)
	if err != nil {
		log.Fatalln("检查程序运行状态时出错：", err)
		return
	}

	if running {
		fmt.Printf("%s is running.\n", program)
	} else {
		fmt.Printf("%s is not running.\n", program)
	}
}
