package gitutil

import (
	"fmt"
	"os/exec"
)

// GetGitDiff 获取指定分支或默认的 diff 内容
// sourceRef、targetRef 均可为空
func GetGitDiff(sourceRef, targetRef string) (string, error) {
	var cmd *exec.Cmd
	if sourceRef == "" && targetRef == "" {
		cmd = exec.Command("git", "diff")
	} else if sourceRef != "" && targetRef != "" {
		cmd = exec.Command("git", "diff", targetRef+".."+sourceRef)
	} else if sourceRef != "" {
		cmd = exec.Command("git", "diff", sourceRef)
	} else {
		cmd = exec.Command("git", "diff", targetRef)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%v\n%s", err, string(out))
	}
	return string(out), nil
}
