package utils

import (
	"ecctl/internal/program"
	"github.com/spf13/cobra"
)

// init 函数用于在包初始化时添加子命令到根命令
func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
}

// Execute 执行根命令，这是程序的入口点
func Execute() error {
	return rootCmd.Execute()
}

// 定义根命令 rootCmd，它是所有子命令的父命令
var rootCmd = &cobra.Command{
	Use:   "ecctl",
	Short: "A tool to manage programs",
	Long:  "A command-line tool to monitor, start, and stop programs.",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

// 定义 startCmd 子命令，用于启动程序
var startCmd = &cobra.Command{
	Use:   "start [program]",
	Short: "Start a program",
	// 参数数量必须为 1
	Args: cobra.ExactArgs(1),
	// 默认运行的 Run 函数，当没有子命令时显示帮助信息
	Run: func(cmd *cobra.Command, args []string) {
		program.Start(args[0])
	},
}

// 定义 stopCmd 子命令，用于停止程序
var stopCmd = &cobra.Command{
	Use:   "stop [program]",
	Short: "Stop a program",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		program.Stop(args[0])
	},
}

// 定义 statusCmd 子命令，用于检查程序状态
var statusCmd = &cobra.Command{
	Use:   "status [program]",
	Short: "Check status of a program",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		program.Status(args[0])
	},
}
