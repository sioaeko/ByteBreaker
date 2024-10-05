
# SegmenGet Downloader

<img src="https://github.com/user-attachments/assets/14436bc4-3c8d-41d1-8d9a-f618df0b6685" width="40%" style="display: block; margin-left: auto; margin-right: auto;">

SegmenGet Downloader is an advanced download manager application that helps users efficiently and quickly download files from any URL. It is developed using Go and the [Fyne](https://fyne.io/) framework, offering features such as multi-segment downloads, proxy settings, and download speed limits.

## Features

- ⚡ **Multi-Agent Download**: Uses multiple agents to download files simultaneously, providing faster download speeds.
- 📁 **Custom Save Path**: Users can specify the path where downloaded files are saved.
- 🔧 **Settings Management**: Offers various settings, including agent count, maximum download speed, proxy configuration, custom User-Agent, and post-download actions.
- 📡 **Proxy Support**: Allows downloads via a proxy server.
- ⏱️ **Download Timeout**: Set a download timeout to stop downloads after a specific period.
- 🚀 **Post-Download Actions**: Options to shut down the computer or automatically open files after a download is complete.
- 📋 **Download List**: Manage a list of downloaded files for easy access.

## Installation

### 1. Install Go
SegmenGet Downloader is written in Go, so you need to have Go installed. You can install Go from the following link:

- [Go Installation Guide](https://golang.org/doc/install)

### 2. Clone and Install the Project

1. Open the terminal and clone the project:

```bash
git clone https://github.com/sioaeko/SegmenGet.git
```

2. Navigate to the project folder:

```bash
cd SegmenGet
```

3. Install dependencies, including the Fyne framework:

```bash
go mod tidy
```

4. Build the program:

```bash
go build -o SegmenGet
```

5. Run the application:

```bash
./SegmenGet
```

## How to Use

SegmenGet Downloader provides an intuitive user interface for managing downloads. Here's a simple guide on how to use it:

### 1. **Enter Download URL**
   Enter the URL of the file you wish to download in the "📡 Download URL" field at the top of the application.

### 2. **Choose Save Path**
   In the "📁 Save Path" field, select the path where the downloaded file will be saved. You can open the folder selection dialog by clicking the folder icon on the right.

### 3. **Manually Set File Name**
   Select the "Set File Name" option to manually enter a file name. If not specified, the file name will be extracted automatically from the URL.

### 4. **Set Agent Count**
   You can set the number of agents to download the file in parallel. More agents can speed up downloads, but performance may vary depending on server load.

### 5. **Monitor Download Progress**
   Click the download button to start downloading, and you can monitor the progress and speed.

### 6. **Post-Download Actions**
   In the settings tab, you can configure actions such as shutting down the system or opening the file automatically after the download is complete.

## Settings

SegmenGet Downloader supports various settings to provide a personalized experience for users:

1. **Default Save Path**: Set the default path where downloaded files will be saved.
2. **Default Agent Count**: The number of agents used by default for downloads.
3. **Maximum Download Speed**: Limit the download speed. Set to 0 for unlimited speed.
4. **Proxy Settings**: Configure a proxy to use for downloads, e.g., `http://127.0.0.1:8080`.
5. **Max Retry Attempts**: Set the maximum number of retries for failed downloads.
6. **Download Timeout**: Set a time limit after which downloads will be stopped.
7. **Custom User-Agent**: Set a custom User-Agent for download requests.
8. **Post-Download Actions**: Configure the app to open the file or shut down the system after downloads complete.

## Contributing

Contributions are always welcome! You can contribute by reporting bugs, suggesting features, or submitting pull requests. To contribute to this project:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push the branch (`git push origin feature-branch`).
5. Create a pull request.

## License

This project is licensed under the MIT License. For more details, see the [LICENSE](LICENSE) file.
