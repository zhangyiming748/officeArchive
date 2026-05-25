package read

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// DocxReader 用于读取 DOCX 文件的结构体
type DocxReader struct {
	filePath string
}

// NewDocxReader 创建一个新的 DOCX 读取器
func NewDocxReader(filePath string) *DocxReader {
	return &DocxReader{
		filePath: filePath,
	}
}

// ReadText 读取 DOCX 文件中的所有文本内容
func (d *DocxReader) ReadText() (string, error) {
	// 打开 DOCX 文件（实际上是 ZIP 文件）
	r, err := zip.OpenReader(d.filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开 DOCX 文件: %v", err)
	}
	defer r.Close()

	var textContent strings.Builder

	// 遍历 ZIP 文件中的所有条目
	for _, f := range r.File {
		// 查找 document.xml 文件，它包含主要的文档内容
		if f.Name == "word/document.xml" {
			text, err := d.extractTextFromXML(f)
			if err != nil {
				return "", fmt.Errorf("提取文本时出错: %v", err)
			}
			textContent.WriteString(text)
			break
		}
	}

	if textContent.Len() == 0 {
		return "", fmt.Errorf("未在 DOCX 文件中找到文本内容")
	}

	return textContent.String(), nil
}

// extractTextFromXML 从 XML 文件中提取文本内容
func (d *DocxReader) extractTextFromXML(file *zip.File) (string, error) {
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
	var doc wordDocument
	err = xml.Unmarshal(content, &doc)
	if err != nil {
		return "", fmt.Errorf("解析 XML 时出错: %v", err)
	}

	// 提取所有文本节点
	var textBuilder strings.Builder
	d.extractTextFromElement(doc.Body, &textBuilder)

	return textBuilder.String(), nil
}

// extractTextFromElement 递归地从 XML 元素中提取文本
func (d *DocxReader) extractTextFromElement(element interface{}, builder *strings.Builder) {
	switch v := element.(type) {
	case wordBody:
		for _, child := range v.Children {
			d.extractTextFromElement(child, builder)
		}
	case wordParagraph:
		for _, child := range v.Children {
			d.extractTextFromElement(child, builder)
		}
	case wordRun:
		for _, child := range v.Children {
			d.extractTextFromElement(child, builder)
		}
	case wordText:
		builder.WriteString(v.Content)
	}
}

// XML 结构定义 - 对应 DOCX 的 XML 格式

type wordDocument struct {
	XMLName xml.Name `xml:"document"`
	Body    wordBody `xml:"body"`
}

type wordBody struct {
	Children []interface{} `xml:",any"`
}

type wordParagraph struct {
	XMLName  xml.Name      `xml:"p"`
	Children []interface{} `xml:",any"`
}

type wordRun struct {
	XMLName  xml.Name      `xml:"r"`
	Children []interface{} `xml:",any"`
}

type wordText struct {
	XMLName xml.Name `xml:"t"`
	Content string   `xml:",chardata"`
}

// UnmarshalXML 自定义解组方法，处理不同类型的子元素
func (b *wordBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				var p wordParagraph
				if err := d.DecodeElement(&p, &t); err != nil {
					return err
				}
				b.Children = append(b.Children, p)
			default:
				// 跳过未知元素
				if err := d.Skip(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// UnmarshalXML 自定义解组方法，处理段落中的子元素
func (p *wordParagraph) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "r":
				var r wordRun
				if err := d.DecodeElement(&r, &t); err != nil {
					return err
				}
				p.Children = append(p.Children, r)
			default:
				// 跳过未知元素
				if err := d.Skip(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// UnmarshalXML 自定义解组方法，处理运行中的子元素
func (r *wordRun) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "t":
				var t_elem wordText
				if err := d.DecodeElement(&t_elem, &t); err != nil {
					return err
				}
				r.Children = append(r.Children, t_elem)
			default:
				// 跳过未知元素
				if err := d.Skip(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ReadDocxFile 便捷函数，直接读取 DOCX 文件并返回文本内容
func ReadDocxFile(filePath string) (string, error) {
	reader := NewDocxReader(filePath)
	return reader.ReadText()
}
