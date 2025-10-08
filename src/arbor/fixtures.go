package arbor

import (
    "bytes"
    "os"
    "path/filepath"
    "text/template"
)

// getFixture renders a fixture template if present under TMS_FIXTURES_DIR,
// otherwise renders the provided fallback template string.
// The fixture is looked up as <dir>/<name>.txt.
func getFixture(name string, data any, fallback string) string {
    dir := os.Getenv("TMS_FIXTURES_DIR")
    tplStr := fallback
    if dir != "" {
        path := filepath.Join(dir, name+".txt")
        if b, err := os.ReadFile(path); err == nil {
            tplStr = string(b)
        }
    }
    tpl, err := template.New(name).Parse(tplStr)
    if err != nil {
        // If a bad template is provided by user, fall back to raw text.
        return tplStr
    }
    var buf bytes.Buffer
    if err := tpl.Execute(&buf, data); err != nil {
        return tplStr
    }
    return buf.String()
}

