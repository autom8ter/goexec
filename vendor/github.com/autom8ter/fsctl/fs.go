package fsctl

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/autom8ter/fsctl/clone"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var errFmt = "ERROR: %s\n%s\n"

type Fs struct {
	assetFunc AssetFunc
	dirFunc   AssetDirFunc
	*afero.Afero
	*viper.Viper
}

func NewFs(c clone.CloneFunc) (*Fs, error) {
	fs := &afero.Afero{
		Fs: afero.NewOsFs(),
	}
	v := viper.GetViper()
	if v == nil {
		v = viper.New()
	}
	v.AutomaticEnv()
	v.SetFs(fs)
	if err := v.BindEnv("CFGTOKEN"); err != nil {
		panic(err)
	}

	f := &Fs{
		Afero: fs,
		Viper: v,
	}
	if err := f.CloneConfig(c); err != nil {
		return f, err
	}

	f.Sync()
	return f, nil
}

func (fs *Fs) WalkTemplates(dir string, outDir string) {

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fs.Exit(1, errFmt, err, "error walking path")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := fs.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, fs.AllSettings())
		}
		return nil
	}); err != nil {
		fs.Exit(1, errFmt, err, "failed to walk templates")
	}
}

func (f *Fs) CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := f.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := f.Afero.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, f.Chmod(dstfile, 0755)
}

func (c *Fs) ScanAndReplaceFile(f afero.File, replacements ...string) {
	nm := f.Name()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}
	if err := c.Remove(f.Name()); err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(strings.NewReader(fmt.Sprintf("%s", d)))
	rep := strings.NewReplacer(replacements...)
	var newstr string
	for scanner.Scan() {
		newstr = rep.Replace(scanner.Text())
		if err := scanner.Err(); err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	newf, err := c.Create(nm)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.WriteString(newf, newstr)
	if err != nil {
		c.Exit(1, errFmt, err, "failed to write string to new file")
	}
	fmt.Println("successfully scanned and replaced: " + f.Name())
}

func (f *Fs) ScanAndReplace(r io.Reader, replacements ...string) string {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	var text string
	for scanner.Scan() {
		text = rep.Replace(scanner.Text())
	}
	return text
}

func (c *Fs) Exit(code int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(code)
}

func (c *Fs) Cd(path string) {
	if err := os.Chdir(path); err != nil {
		c.Exit(1, errFmt, err, "failed to change working directory")
	}
}