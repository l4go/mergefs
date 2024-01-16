# l4go/mergefs ライブラリ

複数のfs.FSを透過的に重ね合わせた(マージした)fs.FSを提供します。

以下のようなコードで、複数のfs.FSを1つにマージすることができます。

``` go
mfs, err := mergefs.New(os.DirFS("/home/hoge"), os.DirFS("/home.old/hoge"))
```

## 詳細仕様

### type MergeFS

複数のfs.FSを透過的に1つのfs.FSへ重ね合わたfs.FS interfaceを提供します。  
また、MergeFSは、以下のinterfaceの機能に対応しています。

- [fs.FS](https://pkg.go.dev/io/fs#FS)
- [fs.StatFS](https://pkg.go.dev/io/fs#StatFS)
- [fs.ReadDirFS](https://pkg.go.dev/io/fs#ReadDirFS)
- [fs.ReadFileFS](https://pkg.go.dev/io/fs#ReadFileFS)
- [fs.Glob()](https://pkg.go.dev/io/fs#Glob)

### func New(fss ...fs.FS) fs.FS

複数のfs.FSからMergeFSを作成します。

MergeFSは、作成時の引数順が前のfs.FSを優先して処理するため、
同じパスのファイルがある場合、引数順の前なfs.FSのファイルがアクセスされます。

