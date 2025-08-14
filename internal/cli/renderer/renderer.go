package renderer

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
)

// Renderer 输出渲染器
type Renderer struct {
	termRenderer *glamour.TermRenderer
}

// NewRenderer 创建新的渲染器
func NewRenderer() (*Renderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(150),
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		termRenderer: r,
	}, nil
}

// RenderMarkdown 渲染Markdown内容
func (r *Renderer) RenderMarkdown(content string) error {
	out, err := r.termRenderer.Render(content)
	if err != nil {
		return err
	}

	fmt.Print(out)
	return nil
}

// RenderPlain 渲染纯文本内容
func (r *Renderer) RenderPlain(content string) {
	fmt.Print(content)
}

// RenderError 渲染错误信息
func (r *Renderer) RenderError(message string) {
	fmt.Fprintf(os.Stderr, "❌ %s\n", message)
}

// RenderSuccess 渲染成功信息
func (r *Renderer) RenderSuccess(message string) {
	fmt.Fprintf(os.Stderr, "✅ %s\n", message)
}

// RenderInfo 渲染信息
func (r *Renderer) RenderInfo(message string) {
	fmt.Fprintf(os.Stderr, "ℹ️  %s\n", message)
}

// RenderWarning 渲染警告信息
func (r *Renderer) RenderWarning(message string) {
	fmt.Fprintf(os.Stderr, "⚠️  %s\n", message)
}

// RenderDiff 渲染diff内容
func (r *Renderer) RenderDiff(diff string) {
	if diff == "" {
		r.RenderInfo("无 diff 变更。")
		return
	}
	r.RenderPlain(diff)
}

// RenderConfig 渲染配置信息
func (r *Renderer) RenderConfig(key, value string) {
	fmt.Printf("%s: %s\n", key, value)
}

// RenderVersion 渲染版本信息
func (r *Renderer) RenderVersion(name, version string) {
	fmt.Printf("%s version %s\n", name, version)
}
