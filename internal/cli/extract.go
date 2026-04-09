package cli

import (
	"fmt"

	"github.com/hueryang/practice/internal/extract"
	"github.com/hueryang/practice/internal/input"
	"github.com/hueryang/practice/internal/llm"
	"github.com/spf13/cobra"
)

func newExtractCmd() *cobra.Command {
	var model string

	cmd := &cobra.Command{
		Use:           "extract <file>",
		Short:         "使用大模型从会议纪要文本中提取主题、结论、待办等关键信息并输出到标准输出",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			data, err := input.ReadMeetingFile(path)
			if err != nil {
				return fmt.Errorf("extract 命令 path=%q 读取文件失败: %w", path, err)
			}

			client, err := llm.NewClientFromEnv()
			if err != nil {
				return fmt.Errorf("extract 命令: %w", err)
			}

			m := model
			if m == "" {
				m = extract.DefaultModel
			}
			minutes, err := extract.FromMinutesText(client, m, string(data))
			if err != nil {
				return fmt.Errorf("extract 命令 path=%q 调用大模型提取失败: %w", path, err)
			}

			out := cmd.OutOrStdout()
			s := extract.FormatHumanReadable(minutes)
			if _, werr := out.Write([]byte(s)); werr != nil {
				return fmt.Errorf("extract 命令 path=%q 写入标准输出失败: %w", path, werr)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&model, "model", extract.DefaultModel, "BigModel 模型名，例如 glm-4-flash-250414、glm-4.7-flash")
	return cmd
}
