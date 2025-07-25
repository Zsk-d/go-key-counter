//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/gofrs/flock"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	// kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procSetWindowsHookEx = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx   = user32.NewProc("CallNextHookEx")
	procGetMessageW      = user32.NewProc("GetMessageW")
	// procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")

	keyCount   = make(map[uint32]int)
	keyCountMu sync.Mutex
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_SYSKEYDOWN  = 0x0104 // Alt / F10 / 系统键
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

var hookCallback = syscall.NewCallback(func(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 && (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN) {
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		keyCountMu.Lock()
		keyCount[kbd.VkCode]++
		keyCountMu.Unlock()

		fmt.Printf("按键 VK:%d 次数: %d\n", kbd.VkCode, keyCount[kbd.VkCode])
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
})

func installHook() error {
	hook, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		hookCallback,
		0,
		0,
	)
	if hook == 0 {
		return fmt.Errorf("hook 安装失败: %v", err)
	}
	return nil
}

func autoSaveJSON() {
	for {
		time.Sleep(time.Minute * 30)

		keyCountMu.Lock()
		data, err := json.MarshalIndent(keyCount, "", "  ")
		keyCountMu.Unlock()

		if err != nil {
			fmt.Println("JSON 编码失败:", err)
			continue
		}

		err = os.WriteFile("data.json", data, 0644)
		if err != nil {
			fmt.Println("写入文件失败:", err)
		} else {
			fmt.Println("数据已保存到 data.json")
		}
	}
}
func loadFromFile() {
	data, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("未找到 data.json，或读取失败，使用空数据开始。")
		return
	}

	keyCountMu.Lock()
	defer keyCountMu.Unlock()

	err = json.Unmarshal(data, &keyCount)
	if err != nil {
		fmt.Println("无法解析 data.json:", err)
		keyCount = make(map[uint32]int) // fallback
	} else {
		fmt.Println("已加载 data.json，记录数量:", len(keyCount))
	}
}
func main() {
	lock := flock.New("app.lock")

	// 尝试获取锁，设置超时时间
	locked, err := lock.TryLock()
	if err != nil {
		fmt.Println("获取锁失败:", err)
		os.Exit(1)
	}

	// 如果未获取到锁，说明已有实例运行
	if !locked {
		fmt.Println("程序已在运行，禁止多开。")
		os.Exit(0)
	}

	// 退出时释放锁
	defer lock.Unlock()

	fmt.Println("键盘统计器已启动")

	loadFromFile()
	err = installHook()
	if err != nil {
		fmt.Println(err)
		return
	}

	go autoSaveJSON()

	var msg struct {
		Hwnd    uintptr
		Message uint32
		WParam  uintptr
		LParam  uintptr
		Time    uint32
		Pt      struct{ X, Y int32 }
	}

	// 监听消息循环
	procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
}
