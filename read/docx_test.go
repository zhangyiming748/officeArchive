package read

import (
	"fmt"
	"testing"
)

func TestReadDocx(t *testing.T) {
	fp := "C:\\Users\\zhang\\Desktop\\2021电子版不含文件\\党员对照检视发言提纲—刘德全20210801.docx"

	// 使用便捷函数读取 DOCX 文件
	text, err := ReadDocxFile(fp)
	if err != nil {
		t.Fatalf("读取 DOCX 文件失败: %v", err)
	}

	// 打印提取的文本内容
	fmt.Println("提取的文本内容:")
	fmt.Println(text)
	fmt.Println("---")

	// 验证是否成功提取到文本
	if text == "" {
		t.Error("未能提取到任何文本内容")
	}

	// 输出文本长度
	fmt.Printf("提取的文本长度: %d 字符\n", len(text))
}
