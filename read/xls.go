package read

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ReadXlsFile 读取 XLS 文件（旧版 Excel 格式）
// 首先使用 LibreOffice 将 XLS 文件转换为 XLSX 文件
// 然后使用 ReadXlsxFile 函数读取转换后的 XLSX 文件
func ReadXlsFile(filePath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", filePath)
	}

	// 获取文件目录和文件名
	dir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)

	// 验证文件扩展名
	if strings.ToLower(fileExt) != ".xls" {
		return "", fmt.Errorf("文件格式不正确，期望 .xls 格式，实际为: %s", fileExt)
	}

	// 生成临时的 XLSX 文件路径
	tempXlsxPath := filepath.Join(dir, strings.TrimSuffix(fileName, fileExt)+".xlsx")

	// 使用 LibreOffice 将 XLS 转换为 XLSX
	err := convertXlsToXlsx(filePath, tempXlsxPath)
	if err != nil {
		return "", fmt.Errorf("转换 XLS 到 XLSX 失败: %v", err)
	}

	// 确保在函数退出时删除临时文件
	defer func() {
		if err := os.Remove(tempXlsxPath); err != nil {
			fmt.Printf("警告: 无法删除临时文件 %s: %v\n", tempXlsxPath, err)
		}
	}()

	// 读取转换后的 XLSX 文件
	text, err := ReadXlsxFile(tempXlsxPath)
	if err != nil {
		return "", fmt.Errorf("读取转换后的 XLSX 文件失败: %v", err)
	}

	return text, nil
}

// convertXlsToXlsx 使用 LibreOffice 将 XLS 文件转换为 XLSX 文件
func convertXlsToXlsx(xlsPath string, xlsxPath string) error {
	// 查找 LibreOffice 可执行文件
	libreOfficePath, err := findLibreOfficeForExcel()
	if err != nil {
		return fmt.Errorf("未找到 LibreOffice: %v", err)
	}

	// 构建命令
	// --headless: 无头模式（不显示 GUI）
	// --convert-to xlsx: 转换为 XLSX 格式
	// --outdir: 输出目录
	cmd := exec.Command(
		libreOfficePath,
		"--headless",
		"--convert-to", "xlsx",
		"--outdir", filepath.Dir(xlsxPath),
		xlsPath,
	)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreOffice 转换失败: %v, 输出: %s", err, string(output))
	}

	// 检查转换后的文件是否存在
	if _, err := os.Stat(xlsxPath); os.IsNotExist(err) {
		return fmt.Errorf("转换后的文件不存在: %s", xlsxPath)
	}

	return nil
}

// findLibreOfficeForExcel 查找 LibreOffice 可执行文件的路径
func findLibreOfficeForExcel() (string, error) {
	// Windows 系统常见的 LibreOffice 安装路径
	windowsPaths := []string{
		`C:\Program Files\LibreOffice\program\soffice.exe`,
		`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
	}

	// Linux 系统常见的 LibreOffice 路径
	linuxPaths := []string{
		"/usr/bin/libreoffice",
		"/usr/bin/soffice",
		"/opt/libreoffice/program/soffice",
	}

	// macOS 系统常见的 LibreOffice 路径
	macOSPaths := []string{
		"/Applications/LibreOffice.app/Contents/MacOS/soffice",
	}

	var searchPaths []string

	// 根据操作系统选择搜索路径
	switch os := runtime.GOOS; os {
	case "windows":
		searchPaths = windowsPaths
	case "linux":
		searchPaths = linuxPaths
	case "darwin":
		searchPaths = macOSPaths
	default:
		// 尝试从 PATH 中查找
		path, err := exec.LookPath("soffice")
		if err == nil {
			return path, nil
		}
		return "", fmt.Errorf("不支持的操作系统: %s", os)
	}

	// 检查每个可能的路径
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 尝试从 PATH 中查找
	path, err := exec.LookPath("soffice")
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("在所有常见路径中都未找到 LibreOffice")
}
