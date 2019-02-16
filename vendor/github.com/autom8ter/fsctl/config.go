package fsctl

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/autom8ter/fsctl/clone"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (c *Fs) Render(s string) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(s)
		if err != nil {
			c.Exit(1, errFmt, err, "failed to render string")
		}
		buf := bytes.NewBuffer(nil)
		if err := t.Execute(buf, c.AllSettings()); err != nil {
			c.Exit(1, errFmt, err, "failed to render string")
		}
		return buf.String()
	}
	return s
}

func (c *Fs) Sync() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		c.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range c.AllSettings() {
		val, ok := v.(string)
		if ok {
			if err := os.Setenv(k, val); err != nil {
				c.Exit(1, errFmt, err, "failed to bind config to env variable")
			}
		}
	}
}


func (f *Fs) CloneConfig(c clone.CloneFunc) error {
	dir, err := f.TempDir("", "config")
	if err != nil {
		panic(err)
	}
	defer f.RemoveAll(dir)
	_ = c.Clone(dir)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		if !info.IsDir() && filepath.Ext(path) == "yaml" || filepath.Ext(path) == "json"{
			file, err := f.Open(path)
			if err != nil {
				return err
			}
			if err := f.ReadConfig(file); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	log.Println("successfully read stored remote config")
	return nil
}
