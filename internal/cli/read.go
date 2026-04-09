package cli

import (
	"fmt"

	"github.com/hueryang/practice/internal/input"
	"github.com/spf13/cobra"
)

func newReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "read <file>",
		Short:         "Print meeting minutes file contents to stdout",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			data, err := input.ReadMeetingFile(path)
			if err != nil {
				return fmt.Errorf("read 命令 path=%q 读取文件失败: %w", path, err)
			}
			out := cmd.OutOrStdout()
			if _, werr := out.Write(data); werr != nil {
				return fmt.Errorf("read 命令 path=%q 写入标准输出失败: %w", path, werr)
			}
			return nil
		},
	}
}
