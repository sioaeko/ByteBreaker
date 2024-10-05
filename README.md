
# ByteBreaker Downloader

ByteBreaker Downloader는 사용자가 원하는 URL에서 파일을 빠르고 효율적으로 다운로드할 수 있도록 돕는 고급 다운로드 관리 애플리케이션입니다. Go와 [Fyne](https://fyne.io/) 프레임워크를 사용하여 개발되었으며, 여러 파일 분할 다운로드, 프록시 설정, 다운로드 속도 제한 등의 기능을 제공합니다.

## 기능

- ⚡ **멀티 에이전트 다운로드**: 여러 에이전트를 사용하여 파일을 동시에 다운로드하여 빠른 다운로드 속도 제공.
- 📁 **저장 경로 선택**: 다운로드한 파일을 저장할 경로를 사용자 정의 가능.
- 🔧 **설정 관리**: 에이전트 수, 최대 다운로드 속도, 프록시 설정, 사용자 정의 User-Agent, 다운로드 후 작업 설정 등 다양한 설정 지원.
- 📡 **프록시 지원**: 프록시 서버를 통해 다운로드할 수 있는 기능 제공.
- ⏱️ **다운로드 타임아웃**: 다운로드 타임아웃 설정을 통해 일정 시간이 지나면 다운로드를 중단.
- 🚀 **다운로드 후 액션**: 다운로드가 완료되면 컴퓨터를 종료하거나, 파일을 자동으로 열 수 있는 옵션 제공.
- 📋 **다운로드 목록**: 다운로드한 파일 목록을 관리하여 사용자에게 제공.

## 설치 방법

### 1. Go 설치
ByteBreaker Downloader는 Go로 작성되었기 때문에 Go가 필요합니다. 아래 링크를 통해 Go를 설치할 수 있습니다.

- [Go 설치 페이지](https://golang.org/doc/install)

### 2. 프로젝트 클론 및 설치

1. 터미널을 열고, 프로젝트를 클론합니다:

```bash
git clone https://github.com/username/bytebreaker-downloader.git
```

2. 프로젝트 폴더로 이동합니다:

```bash
cd bytebreaker-downloader
```

3. 종속성 설치 (Fyne 프레임워크 포함):

```bash
go mod tidy
```

4. 프로그램 빌드:

```bash
go build -o bytebreaker
```

5. 실행:

```bash
./bytebreaker
```

## 사용 방법

ByteBreaker Downloader는 직관적인 사용자 인터페이스를 제공하여 쉽게 다운로드를 관리할 수 있습니다. 아래는 간단한 사용 방법입니다.

### 1. **다운로드 URL 입력**
   애플리케이션 상단의 "📡 다운로드 URL" 필드에 다운로드할 파일의 URL을 입력합니다.

### 2. **저장 경로 선택**
   "📁 저장 경로" 필드에서 다운로드된 파일을 저장할 경로를 선택합니다. 우측의 폴더 아이콘을 눌러 폴더 선택 대화 상자를 열 수 있습니다.

### 3. **파일명 직접 지정**
   "파일명 직접 지정" 옵션을 선택하면 파일명을 직접 입력할 수 있습니다. 입력하지 않으면 URL에서 추출된 파일명이 자동으로 사용됩니다.

### 4. **에이전트 수 설정**
   다운로드를 병렬로 처리할 에이전트 수를 설정할 수 있습니다. 에이전트 수가 많을수록 다운로드 속도가 빨라질 수 있지만, 서버 부하에 따라 성능이 달라질 수 있습니다.

### 5. **다운로드 진행 상황 확인**
   다운로드 버튼을 클릭하면 다운로드가 시작되고, 진행 상황과 다운로드 속도를 확인할 수 있습니다.

### 6. **다운로드 후 작업**
   설정 탭에서 다운로드 완료 후 시스템을 종료하거나, 파일을 자동으로 여는 등의 동작을 설정할 수 있습니다.

## 설정

ByteBreaker Downloader는 다양한 설정을 지원하여 사용자 맞춤형 환경을 제공합니다.

1. **기본 저장 경로**: 기본적으로 다운로드된 파일을 저장할 경로.
2. **기본 에이전트 수**: 다운로드할 때 기본적으로 사용할 에이전트 수.
3. **최대 다운로드 속도**: 다운로드 속도를 제한할 수 있습니다. 0은 제한 없음.
4. **프록시 설정**: 프록시를 통해 다운로드할 수 있도록 설정합니다. 예: `http://127.0.0.1:8080`.
5. **최대 재시도 횟수**: 다운로드 실패 시 재시도할 최대 횟수.
6. **다운로드 타임아웃**: 다운로드가 일정 시간이 지나면 중단될 시간을 설정할 수 있습니다.
7. **사용자 정의 User-Agent**: 사용자 정의 User-Agent를 설정하여 다운로드 요청에 포함할 수 있습니다.
8. **다운로드 후 작업**: 다운로드 완료 후 파일을 열거나 시스템을 종료하도록 설정할 수 있습니다.

## 기여

기여는 언제나 환영합니다! 버그 리포트, 기능 제안 또는 풀 리퀘스트를 통해 이 프로젝트에 기여할 수 있습니다. 이 프로젝트에 기여하려면:

1. 이 저장소를 포크합니다.
2. 새 브랜치를 만듭니다. (`git checkout -b feature-branch`)
3. 변경 사항을 커밋합니다. (`git commit -am 'Add some feature'`)
4. 브랜치를 푸시합니다. (`git push origin feature-branch`)
5. 풀 리퀘스트를 생성합니다.

## 라이선스

이 프로젝트는 MIT 라이선스에 따라 라이선스가 부여됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
