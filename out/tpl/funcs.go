package tpl

import (
	"fmt"
	"html/template"
	"time"
)

func tableCellElapsed(d time.Duration) template.HTML {
	if d.Nanoseconds() == 0 {
		return `<td class="elapsed empty">&mdash;</td>`
	}

	sec := d.Seconds()
	if sec < 0.001 {
		return `<td class="elapsed small">&lt; 1ms</td>`
	}

	intsec := int(sec)
	min := intsec / 60
	return template.HTML(fmt.Sprintf(
		`<td class="elapsed duration">%02d:%06.3f</td>`,
		min,
		sec-float64(min*60),
	))
}

func tableCellCount(value int, bad bool) template.HTML {
	if value == 0 {
		return `<td class="count empty">&nbsp;</td>`
	}

	style := "normal"
	if bad {
		style = "failed"
	}

	return template.HTML(fmt.Sprintf(`<td class="count %s">%d</td>`, style, value))
}
