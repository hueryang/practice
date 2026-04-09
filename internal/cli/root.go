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
		Short:         "从会议纪要生成结构化 PNG 报告",
		Long:          "使用智谱 BigModel 从纪要文本提取主题、结论、待办等，并在当前目录写入 PNG 图片。API Key 请设置环境变量 BIGMODEL_API_KEY（或 ZHIPU_API_KEY）。中文字体可通过 MEETING_REPORT_FONT 指定。",
		SilenceErrors: true,
	}
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
