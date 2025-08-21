package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	keyCount = make(map[uint32]int)
)

// 将按键频率映射到红色梯度（白 → 红）
func getHeatColor(count int, max int) string {
	if max == 0 {
		return "#FFFFFF" // 空白
	}
	// 计算红色强度（255 ~ 0）
	intensity := 255 - (count * 255 / max)
	if intensity < 0 {
		intensity = 0
	}
	// RGB(255, intensity, intensity) -> 红色强度随频率增强
	return fmt.Sprintf("#%02X%02X%02X", 255, intensity, intensity)
}

func writeXlsx(tmpFilePath string, data map[uint32]int) string {
	f, err := excelize.OpenFile(tmpFilePath)
	if err != nil {
		log.Fatal("无法打开文件:", err)
	}

	// 获取最大值用于颜色比例计算
	max := 0
	for _, v := range data {
		if v > max {
			max = v
		}
	}

	for k, v := range data {
		cellLoc, ok := KeyCodeMap[k]
		if !ok {
			continue // 如果映射不存在，跳过
		}

		// 写入按键次数
		err = f.SetCellValue("export", cellLoc, v)
		if err != nil {
			log.Println("写入失败:", err)
			continue
		}

		// ===== 设置热力背景色 =====
		bgColor := getHeatColor(v, max)

		style, err := f.GetCellStyle("export", cellLoc)
		if err != nil {
			log.Printf("无法获取单元格样式 %s: %v", cellLoc, err)
			continue
		}

		// 获取旧样式详细内容
		oldStyle, err := f.GetStyle(style)
		if err != nil {
			log.Printf("无法解析样式: %v", err)
			continue
		}

		// 替换背景颜色，保留原样式（如边框、字体等）
		oldStyle.Fill = excelize.Fill{
			Type:    "pattern",
			Color:   []string{bgColor},
			Pattern: 1,
		}

		// 创建新的样式
		newStyleID, err := f.NewStyle(oldStyle)
		if err == nil {
			_ = f.SetCellStyle("export", cellLoc, cellLoc, newStyleID)
		}

	}

	// 另存为带时间戳的新文件
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("modified_%s.xlsx", timestamp)
	err = f.SaveAs(filename)
	if err != nil {
		log.Fatal("保存失败:", err)
	}
	return filename
}

func main() {
	data, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("未找到 data.json，或读取失败")
		return
	}

	err = json.Unmarshal(data, &keyCount)

	if err != nil {
		fmt.Println("❌ JSON 编码失败:", err)
	}

	filename := writeXlsx("export_temp.xlsx", keyCount)
	fmt.Printf("已保存为 %s\n按任意键退出...", filename)
	var b = make([]byte, 1)
	_, _ = os.Stdin.Read(b)
}

var KeyCodeMap = map[uint32]string{
	8:  "N4",  // "Backspace",
	9:  "A6",  // "Tab",
	13: "N8",  // "Enter",
	19: "R2",  // "Pause",
	20: "A8",  // "CapsLock",
	27: "A2",  // "Esc",
	32: "D12", // "Space",
	33: "R4",  // "PageUp",
	34: "R6",  // "PageDown",
	35: "Q6",  // "End",
	36: "Q4",  // "Home",
	37: "P12", // "Left",
	38: "Q10", // "Up",
	39: "R12", // "Right",
	40: "Q12", // "Down",
	45: "P4",  // "Insert",
	46: "P6",  // "Delete",

	48: "K4", // "0",
	49: "B4", // "1",
	50: "C4", // "2",
	51: "D4", // "3",
	52: "E4", // "4",
	53: "F4", // "5",
	54: "G4", // "6",
	55: "H4", // "7",
	56: "I4", // "8",
	57: "J4", // "9",

	65: "C8",  // A",
	66: "G10", // B",
	67: "E10", // C",
	68: "E8",  // D",
	69: "D6",  // E",
	70: "F8",  // F",
	71: "G8",  // G",
	72: "H8",  // H",
	73: "I6",  // I",
	74: "I8",  // J",
	75: "J8",  // K",
	76: "K8",  // L",
	77: "I10", // M",
	78: "H10", // N",
	79: "J6",  // O",
	80: "K6",  // P",
	81: "B6",  // Q",
	82: "E6",  // R",
	83: "D8",  // S",
	84: "F6",  // T",
	85: "H6",  // U",
	86: "F10", // V",
	87: "C6",  // W",
	88: "D10", // X",
	89: "G6",  // Y",
	90: "C10", // Z",

	91:  "B12", // "LeftWin",
	92:  "L12", // "RightWin",
	93:  "",    // "Apps",
	96:  "T12", // "Numpad0",
	97:  "T10", // "Numpad1",
	98:  "U10", // "Numpad2",
	99:  "V10", // "Numpad3",
	100: "T8",  // "Numpad4",
	101: "U8",  // "Numpad5",
	102: "V8",  // "Numpad6",
	103: "T6",  // "Numpad7",
	104: "U6",  // "Numpad8",
	105: "V6",  // "Numpad9",
	106: "V4",  // "Numpad*",
	107: "W6",  // "Numpad+",
	109: "W4",  // "Numpad-",
	110: "V12", // "Numpad.",
	111: "U4",  // "Numpad/",
	112: "C2",  // "F1",
	113: "D2",  // "F2",
	114: "E2",  // "F3",
	115: "F2",  // "F4",
	116: "G2",  // "F5",
	117: "H2",  // "F6",
	118: "I2",  // "F7",
	119: "J2",  // "F8",
	120: "K2",  // "F9",
	121: "L2",  // "F10",
	122: "M2",  // "F11",
	123: "N2",  // "F12",
	144: "T4",  // "NumLock",
	145: "Q2",  // "ScrollLock",
	160: "A10", // "LeftShift",
	161: "M10", // "RightShift",
	162: "A12", // "LeftCtrl",
	163: "N12", // "RightCtrl",
	164: "C12", // "LeftAlt",
	165: "K12", // "RightAlt",
	186: "L8",  // ";",
	187: "M4",  // "=",
	188: "J10", // ",",
	189: "L4",  // "-",
	190: "K10", // ".",
	191: "L10", // "/",
	192: "A4",  // "`",
	219: "L6",  // "[",
	220: "N6",  // "\\",
	221: "M6",  // "]",
	222: "M8",  // "'",
}
