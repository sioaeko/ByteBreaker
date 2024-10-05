package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"go-downloader/lang" // 다국어 지원을 위해 lang 패키지 임포트
)

// 설정 구조체
type Config struct {
	DefaultSavePath    string `json:"default_save_path"`
	DefaultAgentNum    int    `json:"default_agent_num"`
	MaxDownloadSpeed   int    `json:"max_download_speed"`
	AutoStart          bool   `json:"auto_start"`
	EnableLogFile      bool   `json:"enable_log_file"`
	Proxy              string `json:"proxy"`
	MaxRetryCount      int    `json:"max_retry_count"`
	DownloadTimeout    int    `json:"download_timeout"` // 초 단위
	PostDownloadAction string `json:"post_download_action"`
	CustomUserAgent    string `json:"custom_user_agent"`
}

var config Config
var isPaused bool
var downloadSpeed float64
var downloadHistory []string

func main() {
	// 애플리케이션 생성 및 테마 설정
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow(lang.GetText("title"))

	// 설정 로드
	loadConfig()

	// 탭 1: 다운로드 인터페이스
	downloadTab := makeDownloadTab(w)

	// 탭 2: 다운로드 목록 (상단부터 표시)
	downloadListTab := makeDownloadListTab()

	// 탭 3: 설정 탭
	settingsTab := makeSettingsTab(w)

	// 탭 생성
	tabs := container.NewAppTabs(
		container.NewTabItem(lang.GetText("download"), downloadTab),
		container.NewTabItem("다운로드 목록", downloadListTab),
		container.NewTabItem("환경설정", settingsTab),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	// 창 크기 및 컨텐츠 설정 (탭은 상단, 각 탭 내용은 중앙에 배치)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(500, 600))
	w.ShowAndRun()
}

// 설정 저장
func saveConfig() {
	file, _ := os.Create("config.json")
	defer file.Close()
	json.NewEncoder(file).Encode(config)
}

// 설정 로드
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		config = Config{
			DefaultSavePath:    "",
			DefaultAgentNum:    4,
			MaxDownloadSpeed:   0, // 0은 제한 없음
			AutoStart:          false,
			EnableLogFile:      false,
			Proxy:              "",
			MaxRetryCount:      3,
			DownloadTimeout:    60,     // 60초
			PostDownloadAction: "none", // 옵션: none, shutdown, open_file
			CustomUserAgent:    "SegmenGetDownloader/1.0",
		}
		saveConfig()
		return
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&config)
}

// 다운로드 탭 생성 (중앙 정렬)
func makeDownloadTab(w fyne.Window) *fyne.Container {
	// 타이틀
	title := canvas.NewText(lang.GetText("title"), theme.PrimaryColor())
	title.TextSize = 30
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// URL 입력 필드
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder(lang.GetText("urlPlaceholder"))

	// 저장 경로 선택
	savePathEntry := widget.NewEntry()
	savePathEntry.SetText(config.DefaultSavePath) // 기본 저장 경로
	savePathEntry.SetPlaceHolder(lang.GetText("savePath"))
	browseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			savePathEntry.SetText(uri.Path())
		}, w)
	})

	// 파일명 수동 입력 옵션
	customFilenameCheck := widget.NewCheck("파일명 직접 지정", nil)
	customFilenameEntry := widget.NewEntry()
	customFilenameEntry.SetPlaceHolder("파일명 입력")
	customFilenameEntry.Hide()

	customFilenameCheck.OnChanged = func(checked bool) {
		if checked {
			customFilenameEntry.Show()
		} else {
			customFilenameEntry.Hide()
		}
	}

	// 에이전트 수 입력
	agentEntry := widget.NewEntry()
	agentEntry.SetText(fmt.Sprintf("%d", config.DefaultAgentNum)) // 기본 에이전트 수
	agentEntry.SetPlaceHolder(lang.GetText("agentNum"))

	// 진행 상황 표시
	progress := widget.NewProgressBar()
	progress.Hide()

	// 상태 및 속도 표시
	status := widget.NewLabelWithStyle(lang.GetText("statusWaiting"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	speedLabel := widget.NewLabelWithStyle("속도: 0 KB/s", fyne.TextAlignLeading, fyne.TextStyle{})

	// 다운로드 버튼
	var downloadBtn *widget.Button
	downloadBtn = widget.NewButtonWithIcon(lang.GetText("download"), theme.DownloadIcon(), func() {
		downloadURL := urlEntry.Text
		savePath := savePathEntry.Text
		agents, err := strconv.Atoi(agentEntry.Text) // 에이전트 수를 입력받음
		if err != nil || agents <= 0 {
			dialog.ShowError(fmt.Errorf("잘못된 에이전트 수입니다"), w)
			return
		}

		var filename string
		if customFilenameCheck.Checked {
			filename = customFilenameEntry.Text
		} else {
			filename = getFilenameFromURL(downloadURL)
		}

		if filename == "" {
			dialog.ShowError(fmt.Errorf("파일명을 지정해주세요"), w)
			return
		}

		outputFile := filepath.Join(savePath, filename)

		progress.Show()
		status.SetText(lang.GetText("statusDownloading"))
		downloadBtn.Disable()

		go func() {
			err := downloadFile(downloadURL, outputFile, agents, progress, speedLabel) // 에이전트 수로 다운로드 함수 호출
			if err != nil {
				dialog.ShowError(err, w)
				status.SetText("다운로드 실패")
			} else {
				// 다운로드 성공 시 목록에 추가
				downloadHistory = append(downloadHistory, outputFile)
				dialog.ShowInformation("완료", lang.GetText("statusCompleted"), w)
				status.SetText(lang.GetText("statusCompleted"))
			}
			progress.Hide()
			downloadBtn.Enable()
		}()
	})

	// 레이아웃 구성
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("📡 "+lang.GetText("urlPlaceholder")+":", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		urlEntry,
		widget.NewLabelWithStyle("📁 "+lang.GetText("savePath")+":", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, browseBtn, savePathEntry),
		customFilenameCheck,
		customFilenameEntry,
		container.NewGridWithColumns(2,
			widget.NewLabelWithStyle(lang.GetText("agentNum")+":", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			agentEntry,
		),
		progress,
		speedLabel,
		status,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), downloadBtn),
	)

	// 내용만 중앙 배치
	return container.NewCenter(content)
}

// 다운로드 목록 탭 생성 (상단 정렬)
func makeDownloadListTab() *fyne.Container {
	list := widget.NewList(
		func() int {
			return len(downloadHistory)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("다운로드된 파일")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(downloadHistory[i])
		},
	)

	// 내용은 상단에 정렬
	content := container.NewVBox(
		widget.NewLabelWithStyle("📁 다운로드 목록", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		list,
	)

	// 상단부터 시작 (중앙 정렬 X)
	return content
}

// 설정 탭 생성 (프록시, 재시도, 타임아웃, 다운로드 후 액션 추가)
func makeSettingsTab(w fyne.Window) *fyne.Container {
	// 기본 저장 경로 설정
	savePathEntry := widget.NewEntry()
	savePathEntry.SetText(config.DefaultSavePath)
	savePathEntry.SetPlaceHolder(lang.GetText("savePath"))
	browseBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			savePathEntry.SetText(uri.Path())
		}, w)
	})

	// 기본 에이전트 수 설정
	agentEntry := widget.NewEntry()
	agentEntry.SetText(fmt.Sprintf("%d", config.DefaultAgentNum))
	agentEntry.SetPlaceHolder(lang.GetText("agentNum"))

	// 최대 다운로드 속도 설정 (KB/s)
	speedEntry := widget.NewEntry()
	speedEntry.SetText(fmt.Sprintf("%d", config.MaxDownloadSpeed))
	speedEntry.SetPlaceHolder("최대 다운로드 속도 (KB/s)")

	// 프록시 설정
	proxyEntry := widget.NewEntry()
	proxyEntry.SetText(config.Proxy)
	proxyEntry.SetPlaceHolder("프록시 (예: http://127.0.0.1:8080)")

	// 최대 재시도 횟수 설정
	retryEntry := widget.NewEntry()
	retryEntry.SetText(fmt.Sprintf("%d", config.MaxRetryCount))
	retryEntry.SetPlaceHolder("최대 재시도 횟수")

	// 다운로드 타임아웃 설정
	timeoutEntry := widget.NewEntry()
	timeoutEntry.SetText(fmt.Sprintf("%d", config.DownloadTimeout))
	timeoutEntry.SetPlaceHolder("다운로드 타임아웃 (초 단위)")

	// 다운로드 후 액션 설정
	actionSelect := widget.NewSelect([]string{"none", "shutdown", "open_file"}, func(value string) {
		config.PostDownloadAction = value
	})
	actionSelect.SetSelected(config.PostDownloadAction)

	// 사용자 정의 User-Agent 설정
	userAgentEntry := widget.NewEntry()
	userAgentEntry.SetText(config.CustomUserAgent)
	userAgentEntry.SetPlaceHolder("사용자 정의 User-Agent")

	// 언어 설정 추가
	langSelect := widget.NewSelect([]string{"English", "Korean", "Japanese", "Chinese"}, func(selected string) {
		switch selected {
		case "English":
			lang.SetLanguage(lang.English)
		case "Korean":
			lang.SetLanguage(lang.Korean)
		case "Japanese":
			lang.SetLanguage(lang.Japanese)
		case "Chinese":
			lang.SetLanguage(lang.Chinese)
		}
		updateLanguageUI(w)
	})
	langSelect.SetSelected("English")

	// 설정 저장 버튼
	saveBtn := widget.NewButton("저장", func() {
		config.DefaultSavePath = savePathEntry.Text
		agentNum, err := strconv.Atoi(agentEntry.Text)
		if err != nil || agentNum <= 0 {
			dialog.ShowError(fmt.Errorf("유효한 에이전트 수를 입력해주세요."), w)
			return
		}
		config.DefaultAgentNum = agentNum

		speed, err := strconv.Atoi(speedEntry.Text)
		if err != nil || speed < 0 {
			dialog.ShowError(fmt.Errorf("유효한 속도를 입력해주세요."), w)
			return
		}
		config.MaxDownloadSpeed = speed
		config.Proxy = proxyEntry.Text
		config.MaxRetryCount, _ = strconv.Atoi(retryEntry.Text)
		config.DownloadTimeout, _ = strconv.Atoi(timeoutEntry.Text)
		config.CustomUserAgent = userAgentEntry.Text

		saveConfig()
		dialog.ShowInformation("저장 완료", "설정이 저장되었습니다.", w)
	})

	// 레이아웃 구성
	content := container.NewVBox(
		widget.NewLabelWithStyle("설정", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("📁 "+lang.GetText("savePath"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, browseBtn, savePathEntry),
		widget.NewLabelWithStyle("🔧 "+lang.GetText("agentNum"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		agentEntry,
		widget.NewLabelWithStyle("🚀 최대 다운로드 속도 (KB/s)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		speedEntry,
		widget.NewLabelWithStyle("프록시 설정", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		proxyEntry,
		widget.NewLabelWithStyle("최대 재시도 횟수", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		retryEntry,
		widget.NewLabelWithStyle("다운로드 타임아웃 (초)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		timeoutEntry,
		widget.NewLabelWithStyle("다운로드 완료 후 액션", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		actionSelect,
		widget.NewLabelWithStyle("사용자 정의 User-Agent", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		userAgentEntry,
		widget.NewLabelWithStyle("언어 설정", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		langSelect,
		layout.NewSpacer(),
		saveBtn,
	)

	// 내용 중앙 배치
	return container.NewCenter(content)
}

// URL에서 파일명 추출
func getFilenameFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	path := parsedURL.Path
	segments := strings.Split(path, "/")
	filename := segments[len(segments)-1]

	// URL 디코딩
	filename, err = url.QueryUnescape(filename)
	if err != nil {
		return ""
	}

	return filename
}

// 다운로드 함수 (에이전트 수 지정)
func downloadFile(downloadURL, outputFile string, agents int, progress *widget.ProgressBar, speedLabel *widget.Label) error {
	client := &http.Client{
		Timeout: time.Duration(config.DownloadTimeout) * time.Second,
	}

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	resp, err := client.Head(downloadURL)
	if err != nil {
		return fmt.Errorf("파일 정보를 가져오는 중 오류 발생: %v", err)
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	if fileSize <= 0 {
		return fmt.Errorf("파일 크기를 결정할 수 없습니다")
	}

	partSize := fileSize / int64(agents)
	var wg sync.WaitGroup
	errs := make(chan error, agents)
	progress.Max = 1.0

	tempDir, err := os.MkdirTemp("", "download-")
	if err != nil {
		return fmt.Errorf("임시 디렉토리 생성 실패: %v", err)
	}
	defer os.RemoveAll(tempDir)

	downloadedBytes := int64(0)
	for i := 0; i < agents; i++ {
		wg.Add(1)
		go func(partNum int) {
			defer wg.Done()
			start := int64(partNum) * partSize
			end := start + partSize - 1
			if partNum == agents-1 {
				end = fileSize - 1
			}
			tempFile := filepath.Join(tempDir, fmt.Sprintf("part%d", partNum))
			downloadedPart, err := downloadPart(client, downloadURL, tempFile, start, end)
			if err != nil {
				errs <- fmt.Errorf("파트 %d 다운로드 중 오류: %v", partNum, err)
			}
			atomic.AddInt64(&downloadedBytes, downloadedPart)

			// 다운로드 속도 계산
			elapsed := time.Since(time.Now()).Seconds()
			downloadSpeed = float64(downloadedBytes) / 1024.0 / elapsed
			speedLabel.SetText(fmt.Sprintf("속도: %.2f KB/s", downloadSpeed))

			progress.SetValue(float64(downloadedBytes) / float64(fileSize))
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	err = mergeParts(outputFile, tempDir, agents)
	if err != nil {
		return fmt.Errorf("파일 병합 중 오류: %v", err)
	}

	// 다운로드 후 액션 처리
	switch config.PostDownloadAction {
	case "shutdown":
		fmt.Println("다운로드 완료 후 시스템을 종료합니다.")
	case "open_file":
		fmt.Printf("다운로드한 파일을 엽니다: %s\n", outputFile)
	}

	return nil
}

// 각 파일 파트 다운로드
func downloadPart(client *http.Client, downloadURL, outputFile string, start, end int64) (int64, error) {
	if isPaused {
		time.Sleep(100 * time.Millisecond)
	}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	req.Header.Set("User-Agent", config.CustomUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, resp.Body)
}

// 파일 병합
func mergeParts(outputFile, tempDir string, agents int) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for i := 0; i < agents; i++ {
		partFileName := filepath.Join(tempDir, fmt.Sprintf("part%d", i))
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, partFile)
		partFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// 언어 변경 시 UI 업데이트
func updateLanguageUI(w fyne.Window) {
	w.SetTitle(lang.GetText("title"))
}
