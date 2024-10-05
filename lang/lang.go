package lang

// Language type to represent a language choice
type Language int

const (
	English Language = iota
	Korean
	Japanese
	Chinese
)

var currentLang = English

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	currentLang = lang
}

// GetText returns the localized text based on the current language
func GetText(key string) string {
	texts := map[Language]map[string]string{
		English: {
			"title":             "⚡ SegmenGet Downloader ⚡",
			"urlPlaceholder":    "Enter Download URL",
			"savePath":          "Save Path",
			"agentNum":          "Number of Agents",
			"download":          "Download",
			"statusWaiting":     "Waiting...",
			"statusDownloading": "Downloading...",
			"statusCompleted":   "Download Complete",
		},
		Korean: {
			"title":             "⚡ SegmenGet 다운로더 ⚡",
			"urlPlaceholder":    "다운로드 URL 입력",
			"savePath":          "저장 경로",
			"agentNum":          "에이전트 수",
			"download":          "다운로드",
			"statusWaiting":     "대기 중...",
			"statusDownloading": "다운로드 중...",
			"statusCompleted":   "다운로드 완료",
		},
		Japanese: {
			"title":             "⚡ SegmenGet ダウンローダー ⚡",
			"urlPlaceholder":    "ダウンロード URL を入力",
			"savePath":          "保存パス",
			"agentNum":          "エージェント数",
			"download":          "ダウンロード",
			"statusWaiting":     "待機中...",
			"statusDownloading": "ダウンロード中...",
			"statusCompleted":   "ダウンロード完了",
		},
		Chinese: {
			"title":             "⚡ SegmenGet 下载器 ⚡",
			"urlPlaceholder":    "输入下载URL",
			"savePath":          "保存路径",
			"agentNum":          "代理数量",
			"download":          "下载",
			"statusWaiting":     "等待中...",
			"statusDownloading": "下载中...",
			"statusCompleted":   "下载完成",
		},
	}

	return texts[currentLang][key]
}
