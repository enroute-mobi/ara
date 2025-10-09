package siri

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
)

var (
	stringLt  = []byte("&lt;")
	stringGt  = []byte("&gt;")
	stringAmp = []byte("&amp;")

	// Regexp could be simplified in `{{ ?((?:\.(?!Build)[^. }]+)+) ?}}`
	// but Golang doesn't support lookahead
	re        = regexp.MustCompile(`{{ ?((?:\.(?:(?:[^. }]{0,5})|(?:Buil[^d][^. }]*)|(?:Bui[^l][^. }]*)|(?:Bu[^i][^. }]*)|(?:B[^u][^. }]*)|(?:[^B][^. }]*)))+) ?}}`)
	strEscape = "{{ $1 | stringEscape }}"

	errorType       = reflect.TypeFor[error]()
	fmtStringerType = reflect.TypeFor[fmt.Stringer]()
)

var templates *template.Template

func init() {
	templatePath, err := config.GetTemplateDirectory()
	if err != nil {
		logger.Log.Panicf("Error while loading templates: %v", err)
	}

	t := template.New("").Funcs(template.FuncMap{"stringEscape": StringEscape})
	templates = template.Must(parseGlob(t, filepath.Join(templatePath, "*.template")))
}

func parseGlob(t *template.Template, pattern string) (*template.Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
	}
	return parseFiles(t, readFileOS, filenames...)
}

func parseFiles(t *template.Template, readFile func(string) (string, []byte, error), filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		name, b, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		_, err = tmpl.Parse(re.ReplaceAllString(s, strEscape))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

// see https://github.com/golang/go/blob/8488309192b0ed4b393e2f7b2a93491139ff8ad0/src/text/template/funcs.go#L611
func StringEscape(a any) string {
	s := evalArg(a)
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	var b strings.Builder
	stringEscape(&b, []byte(s))
	return b.String()
}

func stringEscape(w io.Writer, b []byte) {
	last := 0
	for i, c := range b {
		var string []byte
		switch c {
		case '<':
			string = stringLt
		case '>':
			string = stringGt
		case '&':
			string = stringAmp
		default:
			continue
		}
		w.Write(b[last:i])
		w.Write(string)
		last = i + 1
	}
	w.Write(b[last:])
}

func evalArg(arg any) (s string) {
	ok := false
	// Fast path for simple common case.
	s, ok = arg.(string)
	if !ok {
		a, ok := printableValue(reflect.ValueOf(arg))
		if ok {
			arg = a
		} // else let fmt do its thing
		s = fmt.Sprint(arg)
	}
	return s
}

func printableValue(v reflect.Value) (any, bool) {
	if v.Kind() == reflect.Pointer {
		v, _ = indirect(v) // fmt.Fprint handles nil.
	}
	if !v.IsValid() {
		return "<no value>", true
	}

	if !v.Type().Implements(errorType) && !v.Type().Implements(fmtStringerType) {
		if v.CanAddr() && (reflect.PointerTo(v.Type()).Implements(errorType) || reflect.PointerTo(v.Type()).Implements(fmtStringerType)) {
			v = v.Addr()
		} else {
			switch v.Kind() {
			case reflect.Chan, reflect.Func:
				return nil, false
			}
		}
	}
	return v.Interface(), true
}

func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}
