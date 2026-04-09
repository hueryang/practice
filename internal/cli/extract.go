package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hueryang/practice/internal/extract"
	"github.com/hueryang/practice/internal/input"
	"github.com/hueryang/practice/internal/llm"
	"github.com/hueryang/practice/internal/render"
	"github.com/spf13/cobra"
)

func newExtractCmd() *cobra.Command {
	var model string
	var outputPath string

	cmd := &cobra.Command{
		Use:           "extract <file>",
		Short:         "从会议纪要文本提取要点，并在当前目录生成 PNG 图片报告",
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

			out, err := resolveOutputPath(path, outputPath)
			if err != nil {
				return fmt.Errorf("extract 命令: %w", err)
			}

			if err := render.WriteMeetingReportPNG(out, minutes); err != nil {
				return fmt.Errorf("extract 命令 path=%q 写入 PNG %q 失败: %w", path, out, err)
			}
			slog.Info("已生成会议纪要图片", "output", out)
			return nil
		},
	}
	cmd.Flags().StringVar(&model, "model", extract.DefaultModel, "BigModel 模型名，例如 glm-4-flash-250414、glm-4.7-flash")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "输出 PNG 文件路径（默认：当前目录下与输入文件同主文件名的 .png）")
	return cmd
}

func resolveOutputPath(inputPath, flagOut string) (string, error) {
	if strings.TrimSpace(flagOut) != "" {
		p := filepath.Clean(flagOut)
		if filepath.IsAbs(p) {
			return p, nil
		}
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("获取当前目录: %w", err)
		}
		return filepath.Join(wd, p), nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取当前目录: %w", err)
	}
	base := filepath.Base(inputPath)
	stem := strings.TrimSuffix(base, filepath.Ext(base))
	if stem == "" {
		stem = "meeting-report"
	}
	return filepath.Join(wd, stem+".png"), nil
}
