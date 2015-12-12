package hot

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func testHot(conf *Config, t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		tpl, err := New(conf)
		if err != nil {
			t.Error(err)
			w.Done()
			return
		}
		body := "hello {{.Name}}"
		name := filepath.Join(conf.Dir, "hello.tpl")
		err = ioutil.WriteFile(name, []byte(body), 0600)
		if err != nil {
			t.Error(err)
			w.Done()
			return
		}
		time.Sleep(time.Second)

		data := make(map[string]interface{})
		data["Name"] = "gernest"
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, "hello.tpl", data)
		if err != nil {
			t.Error(err)
			w.Done()
			return
		}
		message := "hello gernest"
		if buf.String() != message {
			t.Errorf("expcted %s got %s", message, buf.String())
		}
		w.Done()
	}(&wg)
	wg.Wait()
}

func TestHot(t *testing.T) {
	conf := &Config{
		Watch:          true,
		BaseName:       "hot",
		Dir:            "fixtures",
		FilesExtension: []string{".tpl", ".html", ".tmpl"},
	}
	testHot(conf, t)
	name := filepath.Join(conf.Dir, "hello.tpl")
	ioutil.WriteFile(name, []byte("hello"), 0600)
}
