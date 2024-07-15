package program

import (
	"bytes"
	"ecctl/internal/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getProgramPath 获取当前目录下指定程序的绝对路径。
func getProgramPath(path string) (string, error) {
	logger.Log.Debugf("获取程序路径: %s", path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.Log.Errorf("获取程序路径时出错: %v", err)
		return "", err
	}
	logger.Log.Debugf("程序绝对路径: %s", absPath)
	return absPath, nil
}

// isRunning 检查指定的程序是否正在运行
// 通过调用 pgrep 命令来实现检查
func isRunning(program string) (bool, error) {
	/*programPath, err := getProgramPath(program)
	if err != nil {
		return false, err
	}*/

	logger.Log.Debugf("检查程序是否运行: %s", program)

	// 使用 pgrep 命令来查找正在运行的程序
	cmd := exec.Command("pgrep", "-fl", program)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		// 检查 pgrep 的退出状态码，如果是 1 表示没有找到进程，不算错误
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			logger.Log.Debugf("pgrep 未找到匹配进程: %s", program)
			return false, nil
		}
		logger.Log.Errorf("运行 pgrep 时出错: %v", err)
		return false, err
	}

	logger.Log.Debugf("pgrep 输出:\n%s", out.String())

	// 获取当前程序的 PID
	currentPID := os.Getpid()
	logger.Log.Debugf("当前程序 PID: %d", currentPID)

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
			logger.Log.Debugf("找到匹配进程: %s", line)
			return true, nil
		}
	}

	logger.Log.Debugf("未找到匹配进程: %s", program)
	return false, nil
}

// Start 启动指定的程序
func Start(path string) {
	logger.Log.Infof("尝试启动程序: %s", path)

	// 获取程序的绝对路径。
	absPath, err := getProgramPath(path)
	if err != nil {
		logger.Log.Error("获取程序路径时出错：", err)
		return
	}

	// 检查程序是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		logger.Log.Errorf("程序不存在：%s", absPath)
		return
	}

	// 获取程序的工作目录
	workDir := filepath.Dir(absPath)
	logger.Log.Debugf("目录: %s", workDir)

	program := filepath.Base(absPath)
	// 创建一个新的命令来运行指定的程序
	cmd := exec.Command("./" + program)

	// 检查程序是否已经在运行。
	running, err := isRunning(program)
	if err != nil {
		logger.Log.Error("检查程序运行状态时出错：", err)
		return
	}

	if running {
		logger.Log.Warnf("程序 '%s' 已经在运行", program)
		return
	}

	// 设置命令的工作目录。
	cmd.Dir = workDir
	// 将命令的标准输出和标准错误输出重定向到日志记录器
	cmd.Stdout = logger.Log.Out
	cmd.Stderr = logger.Log.Out

	// 启动命令执行的进程
	err = cmd.Start()
	if err != nil {
		logger.Log.Error("启动程序时出错：", err)
		return
	}

	logger.Log.Infof("程序 '%s' 启动成功", program)

	// 启动一个新的 goroutine 来等待进程的结束
	go func() {
		err := cmd.Wait()
		if err != nil {
			return
		}
		logger.Log.Infof("程序 '%s' 已退出", program)
	}()
}

// Stop 停止指定的程序
// 使用 pkill 命令来终止运行中的程序
func Stop(program string) {
	logger.Log.Infof("尝试停止程序: %s", program)

	// 获取程序的绝对路径。
	/*programPath, err := getProgramPath(program)
	if err != nil {
		logger.Log.Error("获取程序路径时出错：", err)
		return
	}*/

	// 检查程序是否已经在运行。
	running, err := isRunning(program)
	if err != nil {
		logger.Log.Error("检查程序运行状态时出错：", err)
		return
	}

	if !running {
		logger.Log.Warnf("程序 '%s' 没有运行", program)
		return
	}

	// 使用 pkill 命令终止运行中的程序
	cmd := exec.Command("pkill", "-f", program)
	err = cmd.Run()
	if err != nil {
		logger.Log.Error("停止程序时出错：", err)
		return
	}

	// 记录程序停止成功的信息。
	logger.Log.Infof("程序 '%s' 已停止", program)
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
		logger.Log.Error("检查程序运行状态时出错：", err)
		return
	}

	if running {
		logger.Log.Infof("程序 '%s' 正在运行", program)
	} else {
		logger.Log.Infof("程序 '%s' 未运行", program)
	}
}
