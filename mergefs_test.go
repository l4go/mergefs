package mergefs_test

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/l4go/mergefs"
)

//go:embed test/fs1
var fs1Fs embed.FS

//go:embed test/fs2
var fs2Fs embed.FS

func new_fs1() fs.FS {
	fs1fs, err := fs.Sub(fs1Fs, "test/fs1")
	if err != nil {
		panic("crash embed FS")
	}

	return fs1fs
}

func new_fs2() fs.FS {
	fs2fs, err := fs.Sub(fs2Fs, "test/fs2")
	if err != nil {
		panic("crash embed FS")
	}

	return fs2fs
}

func ExampleMergeFS_Open() {
	mfs := mergefs.New(new_fs1(), new_fs2())

	var err error
	_, err = mfs.Open("foo.txt")
	fmt.Println(err == nil)
	_, err = mfs.Open("bar.txt")
	fmt.Println(err == nil)
	_, err = mfs.Open("baz.txt")
	fmt.Println(err == nil)
	// Output:
	// true
	// true
	// true
}

func ExampleMergeFS_Stat() {
	mfs := mergefs.New(new_fs1(), new_fs2())

	var err error
	_, err = fs.Stat(mfs, "foo")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "bar")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "baz")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "boo")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "foo.txt")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "bar.txt")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "baz.txt")
	fmt.Println(err == nil)
	_, err = fs.Stat(mfs, "boo.txt")
	fmt.Println(err == nil)

	// Output:
	// true
	// true
	// true
	// false
	// true
	// true
	// true
	// false
}

func ExampleMergeFS_ReadFile() {
	fs1 := new_fs1()
	fs2 := new_fs2()
	mfs := mergefs.New(fs1, fs2)

	var buf []byte
	var err error
	buf, err = fs.ReadFile(fs1, "baz.txt")
	if err != nil {
		fmt.Println("error")
	} else {
		fmt.Print(string(buf))
	}

	buf, err = fs.ReadFile(fs2, "baz.txt")
	if err != nil {
		fmt.Println("error")
	} else {
		fmt.Print(string(buf))
	}

	buf, err = fs.ReadFile(mfs, "baz.txt")
	if err != nil {
		fmt.Println("error")
	} else {
		fmt.Print(string(buf))
	}

	// Output:
	// baz
	// hoge
	// baz
}

func ExampleMergeFS_ReadDir() {
	mfs := mergefs.New(new_fs1(), new_fs2())

	ents, derr := fs.ReadDir(mfs, ".")
	if derr != nil {
		return
	}

	for _, e := range ents {
		fmt.Println(e.Name())
	}

	// Unordered output:
	// foo
	// bar
	// baz
	// foo.txt
	// bar.txt
	// baz.txt
	// hoge.txt
}

func ExampleMergeFS_Glob() {
	mfs := mergefs.New(new_fs1(), new_fs2())

	lst, err := fs.Glob(mfs, "foo/*.txt")
	if err != nil {
		return
	}

	for _, name := range lst {
		fmt.Println(name)
	}

	// Unordered output:
	// foo/foo.txt
	// foo/bar.txt
	// foo/baz.txt
}
