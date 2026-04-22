package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	white  = "\033[37m"
	bold   = "\033[1m"
	bright = "\033[1m\033[38;2;78;205;196m"
)

type Config struct {
	htmlPath    string
	port        string
	proxyTo     string
	tunnelToken string
}

var config Config

func printMenu() {
	menu := green + bold + `
  ╔══════════════════════════════════════════════════════════╗
  ║                    GO-SERVER                           ║
  ║              All-in-One Local & Public Server            ║
  ╠══════════════════════════════════════════════════════════╣
  ║  [1] Serve HTML Files        - Host your HTML website   ║
  ║  [2] Proxy to Server       - Forward to Flask/other   ║
  ║  [3] Cloudflare Tunnel    - Make server PUBLIC online  ║
  ║  [4] Exit                 - Quit program             ║
  ╚══════════════════════════════════════════════════════════╝
` + reset
	fmt.Println(menu)
}

func ask(prompt string) string {
	fmt.Print(cyan + "  > " + reset + prompt + reset)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func askPort() string {
	fmt.Print(cyan + "  > Enter port [8080]: " + reset)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return "8080"
	}
	return input
}

func runInteractive() {
	printBanner()
	fmt.Println()

	fmt.Println(purple + "  ┌─ SELECT MODE ───────────────────────────────┐" + reset)
	fmt.Println(purple + "  │ " + reset + "1. Serve HTML Files" + strings.Repeat(" ", 23) + purple + "│" + reset)
	fmt.Println(purple + "  │ " + reset + "2. Proxy to Server (Flask/Python)" + strings.Repeat(" ", 16) + purple + "│" + reset)
	fmt.Println(purple + "  │ " + reset + "3. Make Public (Cloudflare)" + strings.Repeat(" ", 20) + purple + "│" + reset)
	fmt.Println(purple + "  └────────────────────────────────────────────┘" + reset)
	fmt.Println()

	choice := ask("Choose [1-3]: ")

	config.port = askPort()

	switch choice {
	case "1":
		config.htmlPath = ask("Enter HTML folder path: ")
	case "2":
		config.proxyTo = ask("Enter server to proxy [localhost:5000]: ")
		if config.proxyTo == "" {
			config.proxyTo = "localhost:5000"
		}
	case "3":
		fmt.Println()
		fmt.Println(yellow + "  ── Cloudflare Setup ──" + reset)
		config.htmlPath = ask("Enter HTML folder path: ")
		if config.htmlPath == "" {
			config.htmlPath = ask("Enter proxy server: ")
		}
		config.tunnelToken = ask("Enter Cloudflare tunnel token: ")
	default:
		fmt.Println(red + "  ⚠ Invalid choice!" + reset)
	}

	fmt.Println()
	startServer()
}

func printBanner() {
	banner := green + bold + `
  ██████╗ ███████╗    ██████╗  ██████╗ ██████╗ 
██╔════╝ ██╔════╝    ██╔══██╗██╔═══██╗██╔══██╗
██║  ███╗█████╗      ██║  ██║██║   ██║██████╔╝
██║  ██║██╔══╝      ██║  ██║██║   ██║██╔══██╗
╚██████╔╝███████╗    ██████╔╝╚██████╔╝██║  ██║
 ╚═════╝ ╚══════╝    ╚═════╝  ╚═════╝ ╚═╝  ╚═╝
` + reset
	fmt.Println(banner)
	fmt.Println(cyan + "  Go-Server v1.0.0 - Your All-in-One Local & Public Server" + reset)
}

func printStatusBox(status, details string) {
	box := purple + "┌" + strings.Repeat("─", 44) + "┐" + reset
	fmt.Println(box)
	fmt.Println(purple + "│ " + reset + green + bold + "STATUS" + reset + strings.Repeat(" ", 35) + purple + "│" + reset)
	fmt.Println(purple + "│ " + reset + "  " + status + strings.Repeat(" ", 41-len(status)) + purple + "│" + reset)
	fmt.Println(purple + "│" + reset + strings.Repeat(" ", 44) + purple + "│" + reset)
	fmt.Println(purple + "│ " + reset + yellow + bold + "DETAILS" + reset + strings.Repeat(" ", 35) + purple + "│" + reset)
	fmt.Println(purple + "│ " + reset + "  " + details + strings.Repeat(" ", 41-len(details)) + purple + "│" + reset)
	fmt.Println(box)
}

func checkCloudflareInstalled() bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cloudflared.exe version")
	case "darwin":
		cmd = exec.Command("cloudflared", "version")
	case "linux":
		cmd = exec.Command("cloudflared", "version")
	}

	if cmd != nil {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		return cmd.Run() == nil
	}
	return false
}

func getCloudflareInstallURL() string {
	switch runtime.GOOS {
	case "windows":
		return "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe"
	case "darwin":
		return "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-darwin-amd64"
	case "linux":
		return "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64"
	}
	return ""
}

func installCloudflare() error {
	fmt.Println(cyan + "  ⬇ Installing Cloudflare Tunnel..." + reset)

	installURL := getCloudflareInstallURL()
	if installURL == "" {
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	resp, err := http.Get(installURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := "cloudflared"
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		os.Chmod(filename, 0755)
	}

	fmt.Println(green + "  ✓ Cloudflare Tunnel installed!" + reset)
	return nil
}

func startCloudflareTunnel(token string) error {
	if !checkCloudflareInstalled() {
		if err := installCloudflare(); err != nil {
			return err
		}
	}

	fmt.Println(cyan + "  🚀 Starting Cloudflare Tunnel..." + reset)

	filename := "cloudflared"
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, filename, "tunnel", "run", "--token", token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println(green + "  ✓ Cloudflare Tunnel started! Your server is now PUBLIC!" + reset)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sigChan
		fmt.Println(yellow + "  ⏹ Stopping Cloudflare Tunnel..." + reset)
		cmd.Process.Kill()
		cancel()
	}()

	return cmd.Wait()
}

func createProxyHandler(target string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return proxy
}

func serveHTML(htmlDir string) http.Handler {
	fs := http.FileServer(http.Dir(htmlDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/" || path == "" {
			indexFiles := []string{"index.html", "index.htm", "default.html"}
			for _, indexFile := range indexFiles {
				indexPath := filepath.Join(htmlDir, indexFile)
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
			}
		}

		filePath := filepath.Join(htmlDir, strings.TrimPrefix(path, "/"))

		if info, err := os.Stat(filePath); err == nil && info.IsDir() {
			for _, indexFile := range []string{"index.html", "index.htm"} {
				indexPath := filepath.Join(filePath, indexFile)
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
			}
		}

		fs.ServeHTTP(w, r)
	})
}

func startServer() {
	var handler http.Handler

	if config.proxyTo != "" {
		handler = createProxyHandler(config.proxyTo)
		fmt.Println(cyan + "  📡 Proxying to: " + config.proxyTo + reset)
		fmt.Println(green + "  🌐 Server running at: http://localhost:" + config.port + reset)
	} else if config.htmlPath != "" {
		absPath, err := filepath.Abs(config.htmlPath)
		if err != nil {
			log.Fatal(err)
		}
		handler = serveHTML(absPath)
		fmt.Println(cyan + "  📂 Serving HTML from: " + absPath + reset)
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Println(yellow + "  ⚠ Warning: Folder does not exist! Create your HTML files first." + reset)
		}
		fmt.Println(green + "  🌐 Server running at: http://localhost:" + config.port + reset)
	} else {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Go-Server</title></head>
<body style="font-family: Arial; background: linear-gradient(135deg, #1A1A2E 0%, #16213E 100%); color: #FFFFFF; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0;">
<h1>🎉 Go-Server Running!</h1>
<p>Use -path flag to serve HTML files</p>
</body>
</html>`))
		})
		fmt.Println(green + "  🌐 Server running at: http://localhost:" + config.port + reset)
	}

	details := fmt.Sprintf("Port: %s | HTML: %s | Proxy: %s", config.port, config.htmlPath, config.proxyTo)
	printStatusBox("Running", details)

	addr := ":" + config.port
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	if config.tunnelToken != "" {
		startCloudflareTunnel(config.tunnelToken)
		return
	}

	fmt.Println(yellow + "\n  Press Ctrl+C to stop server..." + reset)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	<-sigChan
	fmt.Println(yellow + "  ⏹ Shutting down server..." + reset)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}

func main() {
	if len(os.Args) > 1 {
		for i, arg := range os.Args {
			switch arg {
			case "-path", "--path":
				if i+1 < len(os.Args) {
					config.htmlPath = os.Args[i+1]
				}
			case "-port", "--port":
				if i+1 < len(os.Args) {
					config.port = os.Args[i+1]
				}
			case "-proxy", "--proxy":
				if i+1 < len(os.Args) {
					config.proxyTo = os.Args[i+1]
				}
			case "-token", "--token":
				if i+1 < len(os.Args) {
					config.tunnelToken = os.Args[i+1]
				}
			case "-h", "-help", "--help":
				printBanner()
				fmt.Println()
				fmt.Println(purple + "  Usage: go-server.exe [options]" + reset)
				fmt.Println()
				fmt.Println(cyan + "  Options:" + reset)
				fmt.Println(cyan + "    -path <folder>   " + reset + "HTML files folder")
				fmt.Println(cyan + "    -port <num>     " + reset + "Port number (default 8080)")
				fmt.Println(cyan + "    -proxy <addr>   " + reset + "Proxy to server")
				fmt.Println(cyan + "    -token <token>    " + reset + "Cloudflare tunnel token")
				fmt.Println()
				fmt.Println(yellow + "  Run without options for interactive mode!" + reset)
				return
			}
		}
		printBanner()
		startServer()
	} else {
		runInteractive()
	}
}
