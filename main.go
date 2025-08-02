package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/BourgeoisBear/rasterm"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxSize      = 1200.0
	minSize      = 800.0
	charWidthPx  = 8
	charHeightPx = 16
)

type errMsg struct{ err error }

type imageLoadedMsg struct {
	index int
	image string
}

func loadImageCmd(imgPath string, index, termW, termH int) tea.Cmd {
	return func() tea.Msg {
		pxW := termW * charWidthPx
		pxH := (termH / 2) * charHeightPx
		out, err := renderImage(imgPath, pxW, pxH, termW)
		if err != nil {
			return errMsg{err}
		}
		return imageLoadedMsg{index: index, image: out}
	}
}

type model struct {
	imgs     []string
	cur      int
	w, h     int
	status   string
	ready    bool
	loading  bool
	imgPanel string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h, m.ready, m.loading = msg.Width, msg.Height, true, true
		return m, tea.Batch(tea.ClearScreen, loadImageCmd(m.imgs[m.cur], m.cur, m.w, m.h))
	case imageLoadedMsg:
		if msg.index == m.cur {
			m.loading, m.imgPanel = false, msg.image
		}
		return m, nil
	case errMsg:
		m.loading, m.status = false, msg.err.Error()
		return m, nil
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		return handleKeys(m, msg)
	}
	return m, nil
}

func handleKeys(m model, k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	case "l", "right":
		if m.cur < len(m.imgs)-1 {
			m.cur, m.loading, m.status = m.cur+1, true, ""
			return m, tea.Batch(tea.ClearScreen, loadImageCmd(m.imgs[m.cur], m.cur, m.w, m.h))
		}
	case "h", "left":
		if m.cur > 0 {
			m.cur, m.loading, m.status = m.cur-1, true, ""
			return m, tea.Batch(tea.ClearScreen, loadImageCmd(m.imgs[m.cur], m.cur, m.w, m.h))
		}
	case "r":
		if n := len(m.imgs); n > 1 {
			idx := rand.Intn(n - 1)
			if idx >= m.cur {
				idx++
			}
			m.cur, m.loading, m.status = idx, true, ""
			return m, tea.Batch(tea.ClearScreen, loadImageCmd(m.imgs[m.cur], m.cur, m.w, m.h))
		}
	case "s", "enter":
		setWallpaper(m.imgs[m.cur])
		m.status = "Wallpaper set → " + filepath.Base(m.imgs[m.cur])
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return ""
	}

	topH := m.h / 2
	statH := m.h - topH

	var top string
	if m.loading {
		top = lipgloss.NewStyle().Width(m.w).Height(topH).Align(lipgloss.Center, lipgloss.Center).Render("Loading…")
	} else {
		top = m.imgPanel
	}

	var sb strings.Builder
	if m.status != "" {
		sb.WriteString(m.status)
	} else {
		sb.WriteString(fmt.Sprintf("[%d/%d] %s", m.cur+1, len(m.imgs), filepath.Base(m.imgs[m.cur])))
	}
	sb.WriteString("  |  h/j ←  l/k →  r rand  s set  q quit")

	status := lipgloss.NewStyle().Width(m.w).Height(statH).Padding(1, 2).Render(sb.String())
	return top + status
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wallpeek <images-folder>")
		return
	}
	if !rasterm.IsKittyCapable() && !rasterm.IsItermCapable() {
		fmt.Println("Terminal does not support kitty/iTerm graphics.")
		return
	}
	rand.Seed(time.Now().UnixNano())

	imgs, err := scan(os.Args[1])
	if err != nil || len(imgs) == 0 {
		log.Fatalf("%v", err)
	}

	p := tea.NewProgram(model{imgs: imgs}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func scan(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		switch strings.ToLower(filepath.Ext(e.Name())) {
		case ".jpg", ".jpeg", ".png":
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(out)
	return out, nil
}

func renderImage(p string, w, h, termW int) (string, error) {
	fi, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	cache := filepath.Join(os.TempDir(), fmt.Sprintf("wp-%x-%dx%d.esq", fi.ModTime().UnixNano(), w, h))
	if b, err := os.ReadFile(cache); err == nil {
		return string(b), nil
	}

	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	img = resize(img, w, h)
	out, err := encode(img, w)
	if err == nil {
		_ = os.WriteFile(cache, []byte(out), 0o600)
	}
	return out, err
}

func resize(img image.Image, maxW, maxH int) image.Image {
	w, h := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	r := math.Min(float64(maxW)/w, float64(maxH)/h)
	nw, nh := int(w*r), int(h*r)
	if nw == 0 || nh == 0 {
		nw, nh = 1, 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	for y := 0; y < nh; y++ {
		for x := 0; x < nw; x++ {
			srcX := int(float64(x) / float64(nw) * w)
			srcY := int(float64(y) / float64(nh) * h)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

func encode(img image.Image, containerPxW int) (string, error) {
	var buf bytes.Buffer
	imgCharW := int(math.Ceil(float64(img.Bounds().Dx()) / charWidthPx))
	contCharW := int(float64(containerPxW) / charWidthPx)
	if pad := (contCharW - imgCharW) / 2; pad > 0 {
		buf.WriteString(strings.Repeat(" ", pad))
	}

	var err error
	switch {
	case rasterm.IsKittyCapable():
		err = rasterm.KittyWriteImage(&buf, img, rasterm.KittyImgOpts{})
	case rasterm.IsItermCapable():
		err = rasterm.ItermWriteImage(&buf, img)
	default:
		err = fmt.Errorf("no graphics protocol")
	}
	buf.WriteByte('\n')
	return buf.String(), err
}

func setWallpaper(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder" to set desktop picture to POSIX file "%s"`, path))
	case "linux":
		const waypaper = "waypaper"
		if _, err := exec.LookPath(waypaper); err == nil {
			cmd = exec.Command(waypaper, "--wallpaper", path)
		} else {
			cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+path)
		}
	default:
		fmt.Println("Wallpaper change not implemented on", runtime.GOOS)
		return
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
