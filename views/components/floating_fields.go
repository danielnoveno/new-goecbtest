/*
   file:           views/components/floating_fields.go
   description:    Komponen UI umum untuk floating fields
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type floatingEntry struct {
	widget.Entry
	focused        bool
	onStateChanged func()
}

func newFloatingEntry(password bool, onStateChanged func()) *floatingEntry {
	entry := &floatingEntry{onStateChanged: onStateChanged}
	entry.Password = password
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *floatingEntry) FocusGained() {
	e.Entry.FocusGained()
	e.focused = true
	if e.onStateChanged != nil {
		e.onStateChanged()
	}
}

func (e *floatingEntry) FocusLost() {
	e.Entry.FocusLost()
	e.focused = false
	if e.onStateChanged != nil {
		e.onStateChanged()
	}
}

func (e *floatingEntry) TypedRune(r rune) {
	e.Entry.TypedRune(r)
	if e.onStateChanged != nil {
		e.onStateChanged()
	}
}

func (e *floatingEntry) TypedKey(k *fyne.KeyEvent) {
	e.Entry.TypedKey(k)
	if e.onStateChanged != nil {
		e.onStateChanged()
	}
}

type floatingLabelLayout struct {
	input       fyne.CanvasObject
	label       *canvas.Text
	shouldFloat func() bool
}

func (l *floatingLabelLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	floatUp := l.shouldFloat()
	if floatUp {
		l.label.TextSize = theme.TextSize() - 2
		l.label.Color = theme.PrimaryColor()
	} else {
		l.label.TextSize = theme.TextSize()
		l.label.Color = theme.PlaceHolderColor()
	}
	l.label.Refresh()

	labelSize := l.label.MinSize()
	labelPadding := theme.InnerPadding()
	entryTop := labelSize.Height / 2
	entryHeight := size.Height - entryTop
	if entryHeight < 0 {
		entryHeight = 0
	}

	l.input.Resize(fyne.NewSize(size.Width, entryHeight))
	l.input.Move(fyne.NewPos(0, entryTop))

	if floatUp {
		l.label.Move(fyne.NewPos(labelPadding, 0))
	} else {
		centerY := entryTop + (entryHeight-labelSize.Height)/2
		l.label.Move(fyne.NewPos(labelPadding, centerY))
	}
}

func (l *floatingLabelLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	labelSize := l.label.MinSize()
	inputSize := l.input.MinSize()
	return fyne.NewSize(inputSize.Width, inputSize.Height+labelSize.Height/2)
}

type FloatingEntryField struct {
	Entry *floatingEntry
	root  *fyne.Container
}

func NewFloatingEntry(label string, password bool) *FloatingEntryField {
	var refresh func()
	entry := newFloatingEntry(password, func() {
		if refresh != nil {
			refresh()
		}
	})
	entry.SetPlaceHolder("")

	labelText := canvas.NewText(label, theme.PlaceHolderColor())
	labelText.TextStyle = fyne.TextStyle{Bold: true}

	layout := &floatingLabelLayout{
		input: entry,
		label: labelText,
		shouldFloat: func() bool {
			return entry.focused || entry.Text != ""
		},
	}
	container := container.New(layout, entry, labelText)
	refresh = container.Refresh

	return &FloatingEntryField{
		Entry: entry,
		root:  container,
	}
}

func (f *FloatingEntryField) Object() fyne.CanvasObject {
	return f.root
}

func (f *FloatingEntryField) Text() string {
	return f.Entry.Text
}

func (f *FloatingEntryField) SetText(val string) {
	f.Entry.SetText(val)
	f.root.Refresh()
}

type FloatingSelectField struct {
	Select *widget.Select
	root   *fyne.Container
}

func NewFloatingSelect(label string, options []string, onChanged func(string)) *FloatingSelectField {
	var refresh func()
	selectWidget := widget.NewSelect(options, func(val string) {
		if refresh != nil {
			refresh()
		}
		if onChanged != nil {
			onChanged(val)
		}
	})
	selectWidget.PlaceHolder = ""

	labelText := canvas.NewText(label, theme.PlaceHolderColor())
	labelText.TextStyle = fyne.TextStyle{Bold: true}

	layout := &floatingLabelLayout{
		input: selectWidget,
		label: labelText,
		shouldFloat: func() bool {
			return selectWidget.Selected != ""
		},
	}
	container := container.New(layout, selectWidget, labelText)
	refresh = container.Refresh

	return &FloatingSelectField{
		Select: selectWidget,
		root:   container,
	}
}

func (f *FloatingSelectField) Object() fyne.CanvasObject {
	return f.root
}
