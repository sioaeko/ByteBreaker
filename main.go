package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("분할 다운로더")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("다운로드 URL")

	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("저장할 파일 이름")

	partsEntry := widget.NewEntry()
	partsEntry.SetText("4")
	partsEntry.SetPlaceHolder("분할 수")

	progress := widget.NewProgressBar()
	progress.Hide()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "URL", Widget: urlEntry},
			{Text: "출력 파일", Widget: outputEntry},
			{Text: "분할 수", Widget: partsEntry},
		},
	}

	downloadBtn := widget.NewButton("다운로드", func() {
		url := urlEntry.Text
		outputFile := outputEntry.Text
		parts, err := strconv.Atoi(partsEntry.Text)
		if err != nil || parts <= 0 {
			dialog.ShowError(fmt.Errorf("잘못된 분할 수입니다"), w)
			return
		}

		progress.Show()
		go func() {
			err := downloadFile(url, outputFile, parts, progress)
			if err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("완료", "다운로드가 완료되었습니다", w)
			}
			progress.Hide()
		}()
	})

	content := container.NewVBox(form, progress, downloadBtn)
	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 200))
	w.ShowAndRun()
}

func downloadFile(url, outputFile string, parts int, progress *widget.ProgressBar) error {
	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("파일 정보를 가져오는 중 오류 발생: %v", err)
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	if fileSize <= 0 {
		return fmt.Errorf("파일 크기를 결정할 수 없습니다")
	}

	partSize := fileSize / int64(parts)
	var wg sync.WaitGroup
	errs := make(chan error, parts)
	progress.Max = float64(fileSize)

	for i := 0; i < parts; i++ {
		wg.Add(1)
		go func(partNum int) {
			defer wg.Done()
			start := int64(partNum) * partSize
			end := start + partSize - 1
			if partNum == parts-1 {
				end = fileSize - 1
			}

			err := downloadPart(url, fmt.Sprintf("%s.part%d", outputFile, partNum), start, end, progress)
			if err != nil {
				errs <- fmt.Errorf("파트 %d 다운로드 중 오류: %v", partNum, err)
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	err = mergeParts(outputFile, parts)
	if err != nil {
		return fmt.Errorf("파일 병합 중 오류: %v", err)
	}

	return nil
}

func downloadPart(url, outputFile string, start, end int64, progress *widget.ProgressBar) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, werr := out.Write(buf[:n])
			if werr != nil {
				return werr
			}
			progress.SetValue(progress.Value + float64(n))
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func mergeParts(outputFile string, parts int) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for i := 0; i < parts; i++ {
		partFileName := fmt.Sprintf("%s.part%d", outputFile, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, partFile)
		partFile.Close()
		os.Remove(partFileName)
		if err != nil {
			return err
		}
	}

	return nil
}
