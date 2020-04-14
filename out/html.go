package out

import (
	"html/template"
	"os"
	"strconv"

	"github.com/mono83/oscar/out/tpl"
)

// WriteHTMLFiles writes report into separate HTML files, placed in required folder
func WriteHTMLFiles(path string, r *Report) error {
	if err := os.Mkdir(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	// Reading templates
	t, err := tpl.Load()
	if err != nil {
		return err
	}

	// CSS
	if err := execTplIntoFile("css", path, "main.css", t, nil); err != nil {
		return err
	}
	// Summary
	if err := execTplIntoFile("summary", path, "summary.html", t, r); err != nil {
		return err
	}
	// Remotes
	if err := execTplIntoFile("remotes", path, "remotes-avg.html", t, r.TopRemoteRequests(50, "avg")); err != nil {
		return err
	}
	if err := execTplIntoFile("remotes", path, "remotes-max.html", t, r.TopRemoteRequests(50, "max")); err != nil {
		return err
	}
	if err := execTplIntoFile("remotes", path, "remotes-sum.html", t, r.TopRemoteRequests(50, "sum")); err != nil {
		return err
	}
	// Suites
	for _, suite := range r.Suites() {
		if err := execTplIntoFile("suite", path, "suite-"+strconv.Itoa(suite.ID)+".html", t, suite); err != nil {
			return err
		}
	}

	return nil
}

func execTplIntoFile(templateName, path, file string, t *template.Template, data interface{}) error {
	f, err := os.Create(path + string(os.PathSeparator) + file)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.ExecuteTemplate(f, templateName, data)
}
