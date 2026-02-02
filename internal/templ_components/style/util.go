package style

import (
  "github.com/a-h/templ"
)

func CSSProperty(value string) templ.SafeCSSProperty {
  return templ.SafeCSSProperty(value)
}
