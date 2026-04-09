package cli

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd builds the root command (exported for tests).
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "meetingreport",
		Short:         "Meeting minutes utilities: read files and (later) generate image reports",
		Long:          "会议纪要工具：read 打印原文；extract 使用智谱 BigModel 提取主题、结论、待办等并输出到标准输出。API Key 请设置环境变量 BIGMODEL_API_KEY（或 ZHIPU_API_KEY）。",
		SilenceErrors: true,
	}
	root.AddCommand(newReadCmd())
	root.AddCommand(newExtractCmd())
	return root
}

// Execute runs the CLI using os.Args[1:], logs failures to stderr via slog, and returns an exit code.
func Execute() int {
	return ExecuteWithArgs(os.Args[1:])
}

// ExecuteWithArgs runs the root command with the given arguments (excluding program name) and logs any failure.
func ExecuteWithArgs(argv []string) int {
	cmd := NewRootCmd()
	cmd.SetArgs(argv)
	if err := cmd.Execute(); err != nil {
		slog.Error("命令执行失败", "error", err)
		return 1
	}
	return 0
}
