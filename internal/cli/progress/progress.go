package progress

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ProgressBar 进度条结构
type ProgressBar struct {
	total     int
	current   int
	width     int
	startTime time.Time
	message   string
}

// NewProgressBar 创建新的进度条
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		total:     total,
		current:   0,
		width:     50,
		startTime: time.Now(),
		message:   message,
	}
}

// Update 更新进度
func (p *ProgressBar) Update(current int) {
	p.current = current
	p.render()
}

// Increment 增加进度
func (p *ProgressBar) Increment() {
	p.current++
	p.render()
}

// render 渲染进度条
func (p *ProgressBar) render() {
	if p.total == 0 {
		return
	}

	percentage := float64(p.current) / float64(p.total)
	filled := int(float64(p.width) * percentage)

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)

	elapsed := time.Since(p.startTime)
	eta := time.Duration(0)
	if p.current > 0 {
		eta = time.Duration(float64(elapsed) * float64(p.total-p.current) / float64(p.current))
	}

	fmt.Fprintf(os.Stderr, "\r%s [%s] %d/%d (%.1f%%) %s ETA: %s",
		p.message, bar, p.current, p.total, percentage*100, elapsed.Round(time.Second), eta.Round(time.Second))
}

// Finish 完成进度条
func (p *ProgressBar) Finish() {
	p.current = p.total
	p.render()
	fmt.Fprintln(os.Stderr)
}

// Spinner 旋转指示器
type Spinner struct {
	chars    []string
	current  int
	message  string
	stopChan chan bool
}

// NewSpinner 创建新的旋转指示器
func NewSpinner(message string) *Spinner {
	return &Spinner{
		chars:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		current:  0,
		message:  message,
		stopChan: make(chan bool),
	}
}

// Start 开始旋转
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(os.Stderr, "\r%s %s", s.message, s.chars[s.current])
				s.current = (s.current + 1) % len(s.chars)
			case <-s.stopChan:
				return
			}
		}
	}()
}

// Stop 停止旋转
func (s *Spinner) Stop() {
	s.stopChan <- true
	fmt.Fprintln(os.Stderr)
}

// SimpleProgress 简单进度显示
type SimpleProgress struct {
	message string
}

// NewSimpleProgress 创建简单进度显示
func NewSimpleProgress(message string) *SimpleProgress {
	return &SimpleProgress{message: message}
}

// Show 显示进度消息
func (s *SimpleProgress) Show(message string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", s.message, message)
}

// Success 显示成功消息
func (s *SimpleProgress) Success(message string) {
	fmt.Fprintf(os.Stderr, "✅ %s: %s\n", s.message, message)
}

// Error 显示错误消息
func (s *SimpleProgress) Error(message string) {
	fmt.Fprintf(os.Stderr, "❌ %s: %s\n", s.message, message)
}

// Info 显示信息消息
func (s *SimpleProgress) Info(message string) {
	fmt.Fprintf(os.Stderr, "ℹ️  %s: %s\n", s.message, message)
}
