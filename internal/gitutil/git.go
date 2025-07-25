package gitutil

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetGitDiff 获取指定分支或默认的 diff 内容
// sourceRef、targetRef 均可为空

func GetGitDiff(sourceRef, targetRef string) (string, error) {
	if isEmptyRef(sourceRef) {
		sourceRef = "--"
	}
	if isEmptyRef(targetRef) {
		targetRef = "HEAD"
	}

	var result strings.Builder

	// 1. 获取常规diff
	regularDiff, err := runGitCommand("diff", targetRef, sourceRef)
	if err != nil {
		return "", fmt.Errorf("获取diff失败: %w", err)
	}
	result.WriteString(regularDiff)

	// 2. 获取新增文件diff
	untrackedFiles, err := runGitCommand("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", fmt.Errorf("获取未跟踪文件失败: %w", err)
	}

	for _, file := range strings.Split(strings.TrimSpace(untrackedFiles), "\n") {
		if file == "" {
			continue
		}

		// 执行diff命令
		fileDiff, err := runGitDiffForNewFile(file)
		if err != nil {
			// 即使diff失败也继续处理其他文件
			// continue
			return "", fmt.Errorf("获取未跟踪文件diff失败 %s: %w", file, err)
		}
		result.WriteString(fileDiff)
	}

	return result.String(), nil
}

func runGitDiffForNewFile(file string) (string, error) {
	// 使用更可靠的方式执行diff
	cmd := exec.Command("git", "diff", "--no-index", "--", "/dev/null", file)

	// 设置完整的环境
	cmd.Dir, _ = os.Getwd()
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1", // 忽略系统级配置
		"GIT_TERMINAL_PROMPT=0", // 禁用终端提示
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// 检查是否是我们可以忽略的错误
		if isBenignDiffError(err, output) {
			return string(output), nil
		}
		return "", fmt.Errorf("获取diff失败 %s: %w\n%s", file, err, string(output))
	}
	return string(output), nil
}

func isBenignDiffError(err error, output []byte) bool {
	// Git diff --no-index 对于新文件有时会返回1但输出是有效的
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 {
			// 检查输出是否包含有效的diff内容
			return strings.Contains(string(output), "diff --git")
		}
	}
	return false
}

func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir, _ = os.Getwd()
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git %s: %w\n%s",
			strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}

func isEmptyRef(ref string) bool {
	return ref == "" || ref == "."
}
