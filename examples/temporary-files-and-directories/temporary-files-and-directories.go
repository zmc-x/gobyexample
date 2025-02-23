// 在程序运行时，我们经常创建一些运行时用到，程序结束后就不再使用的数据。
// *临时目录和文件* 对于上面的情况很有用，因为它不会随着时间的推移而污染文件系统。

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// 创建临时文件最简单的方法是调用 `os.CreateTemp` 函数。
	// 它会创建并打开文件，我们可以对文件进行读写。
	// 函数的第一个参数传 `""`，`os.CreateTemp` 会在操作系统的默认位置下创建该文件。
	f, err := os.CreateTemp("", "sample")
	check(err)

	// 打印临时文件的名称。
	// 文件名以 `os.CreateTemp` 函数的第二个参数作为前缀，
	// 剩余的部分会自动生成，以确保并发调用时，生成不重复的文件名。
	// 在类 Unix 操作系统下，临时目录一般是 `/tmp`。
	fmt.Println("Temp file name:", f.Name())

	// defer 删除该文件。
	// 尽管操作系统会自动在某个时间清理临时文件，但手动清理是一个好习惯。
	defer os.Remove(f.Name())
	
	// 当然如果这样去移除所创建的临时文件，那么其实会出现error（`The process cannot access the file because it is being used by another process.`），
	// 首先解释一下这个error为什么会出现或者说为什么只有一个`defer os.Remove()`是无法删掉临时文件的。
	// 其实很简单，就是所创建的file可以还在“写”的过程，那么当我们去删掉该文件的时候，就会出现error，这是很显而易见的error；
	// 那么解决方法也很简单，结束“写”的过程；那么就可以正常的删除了。
	defer f.Close()

	// 我们可以向文件写入一些数据。
	_, err = f.Write([]byte{1, 2, 3, 4})
	check(err)

	// 如果需要写入多个临时文件，最好是为其创建一个临时 *目录* 。
	// `os.MkdirTemp` 的参数与 `CreateTemp` 相同，
	// 但是它返回的是一个 *目录名* ，而不是一个打开的文件。
	dname, err := os.MkdirTemp("", "sampledir")
	check(err)
	fmt.Println("Temp dir name:", dname)

	// 对于目录不存在上述问题。
	defer os.RemoveAll(dname)

	// 现在，我们可以通过拼接临时目录和临时文件合成完整的临时文件路径，并写入数据。
	fname := filepath.Join(dname, "file1")
	err = os.WriteFile(fname, []byte{1, 2}, 0666)
	check(err)
}
