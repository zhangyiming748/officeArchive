package read

import (
	"fmt"
	"testing"
)

func TestReadDoc(t *testing.T) {
	// 注意：这个测试需要一个实际的 .doc 文件
	// 请根据实际情况修改文件路径
	fp := "C:\\Users\\zhang\\Desktop\\test.doc"

	// 使用 ReadDocFile 函数读取 DOC 文件
	text, err := ReadDocFile(fp)
	if err != nil {
		t.Fatalf("读取 DOC 文件失败: %v", err)
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
