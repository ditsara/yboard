package main

import (
	"os/exec"
	"strings"
)

// executeClipboardCopy copies buffer to the system clipboard.
// Probes available tools in priority order: clip.exe (WSL) → wl-copy (Wayland) → xclip (X11).
// Returns a user-facing status message.
func executeClipboardCopy(buffer string) string {
	if len(buffer) == 0 {
		return "Buffer empty"
	}

	var cmd *exec.Cmd
	if _, err := exec.LookPath("clip.exe"); err == nil {
		cmd = exec.Command("clip.exe")
	} else if _, err := exec.LookPath("wl-copy"); err == nil {
		cmd = exec.Command("wl-copy")
	} else {
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	cmd.Stdin = strings.NewReader(buffer)
	if err := cmd.Run(); err == nil {
		return "📋 Copied text block to clipboard!"
	}
	return "❌ Clipboard process pipeline error"
}
