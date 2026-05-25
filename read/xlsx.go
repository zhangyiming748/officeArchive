package read

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// XlsxReader 用于读取 XLSX 文件的结构体
type XlsxReader struct {
	filePath string
}

// NewXlsxReader 创建一个新的 XLSX 读取器
func NewXlsxReader(filePath string) *XlsxReader {
	return &XlsxReader{
		filePath: filePath,
	}
}

// ReadText 读取 XLSX 文件中的所有文本内容
func (x *XlsxReader) ReadText() (string, error) {
	// 打开 XLSX 文件（实际上是 ZIP 文件）
	r, err := zip.OpenReader(x.filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开 XLSX 文件: %v", err)
	}
	defer r.Close()

	var textContent strings.Builder

	// 遍历 ZIP 文件中的所有条目，查找工作表 XML 文件
	for _, f := range r.File {
		// XLSX 的工作表文件位于 xl/worksheets/sheet*.xml
		if strings.HasPrefix(f.Name, "xl/worksheets/sheet") && strings.HasSuffix(f.Name, ".xml") {
			text, err := x.extractTextFromXML(f)
			if err != nil {
				return "", fmt.Errorf("提取工作表文本时出错: %v", err)
			}
			textContent.WriteString(text)
			textContent.WriteString("\n") // 不同工作表之间添加换行
		}
	}

	if textContent.Len() == 0 {
		return "", fmt.Errorf("未在 XLSX 文件中找到文本内容")
	}

	return textContent.String(), nil
}

// extractTextFromXML 从工作表 XML 文件中提取文本内容
func (x *XlsxReader) extractTextFromXML(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	// 读取 XML 内容
	content, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	// 解析 XML 并提取文本
	var worksheet xlsxWorksheet
	err = xml.Unmarshal(content, &worksheet)
	if err != nil {
		return "", fmt.Errorf("解析 XML 时出错: %v", err)
	}

	// 提取所有单元格文本
	var textBuilder strings.Builder
	x.extractTextFromSheetData(worksheet.SheetData, &textBuilder)

	return textBuilder.String(), nil
}

// extractTextFromSheetData 从工作表数据中提取文本
func (x *XlsxReader) extractTextFromSheetData(sheetData xlsxSheetData, builder *strings.Builder) {
	for _, row := range sheetData.Rows {
		for _, cell := range row.Cells {
			if cell.Value != "" {
				builder.WriteString(cell.Value)
				builder.WriteString("\t") // 单元格之间用制表符分隔
			}
		}
		builder.WriteString("\n") // 行之间用换行符分隔
	}
}

// XML 结构定义 - 对应 XLSX 的 XML 格式

type xlsxWorksheet struct {
	XMLName   xml.Name      `xml:"worksheet"`
	SheetData xlsxSheetData `xml:"sheetData"`
}

type xlsxSheetData struct {
	Rows []xlsxRow `xml:"row"`
}

type xlsxRow struct {
	Cells []xlsxCell `xml:"c"`
}

type xlsxCell struct {
	Value string `xml:"v"`
}

// ReadXlsxFile 便捷函数，直接读取 XLSX 文件并返回文本内容
func ReadXlsxFile(filePath string) (string, error) {
	reader := NewXlsxReader(filePath)
	return reader.ReadText()
}
