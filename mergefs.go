package mergefs

import (
	"io/fs"
)

type MergeFS struct {
	fss []fs.FS
}

func New(fss ...fs.FS) fs.FS {
	return &MergeFS{fss: fss}
}

func (mfs *MergeFS) Open(name string) (fs.File, error) {
	var f fs.File
	var err error
	for _, fsys := range mfs.fss {
		f, err = fsys.Open(name)
		if err == nil {
			break
		}
	}

	return f, err
}

func (mfs *MergeFS) Stat(name string) (fs.FileInfo, error) {
	var last_err error
	var newer_fi fs.FileInfo

	all_err := true
	for _, fsys := range mfs.fss {
		fi, err := fs.Stat(fsys, name)
		if err != nil {
			last_err = err
			continue
		}

		all_err = false
		last_err = nil
		if !fi.IsDir() {
			newer_fi = fi
			break
		}

		if newer_fi == nil || fi.ModTime().After(newer_fi.ModTime()) {
			newer_fi = fi
		}
	}

	if all_err {
		return nil, last_err
	}

	return newer_fi, nil
}

func (mfs *MergeFS) ReadFile(name string) ([]byte, error) {
	var buf []byte
	var err error
	for _, fsys := range mfs.fss {
		buf, err = fs.ReadFile(fsys, name)
		if err == nil {
			break
		}
	}

	return buf, err
}

func (mfs *MergeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	uniq := map[string]struct{}{}
	all_dent := []fs.DirEntry{}

	var err error
	for _, fsys := range mfs.fss {
		dent, e := fs.ReadDir(fsys, name)
		if e != nil {
			err = e
			continue
		}
		err = nil

		for _, ent := range dent {
			n := ent.Name()
			if _, ok := uniq[n]; ok {
				continue
			}

			uniq[n] = struct{}{}
			all_dent = append(all_dent, ent)
		}
	}

	if err != nil {
		return nil, err
	}

	return all_dent, nil
}

func (mfs *MergeFS) Glob(pattern string) ([]string, error) {
	uniq := map[string]struct{}{}

	var last_err error
	all_err := true
	for _, fsys := range mfs.fss {
		lst, e := fs.Glob(fsys, pattern)
		last_err = e
		if e == nil {
			all_err = false
		}

		for _, n := range lst {
			uniq[n] = struct{}{}
		}
	}

	if all_err {
		return nil, last_err
	}

	list := make([]string, 0, len(uniq))
	for k := range uniq {
		list = append(list, k)
	}

	return list, nil
}
