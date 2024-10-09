// Code generated by github.com/worldiety/macro. DO NOT EDIT.

package domain

import (
	json "encoding/json"
	fmt "fmt"
)

// This variable is declared to let linters know, that [_Component] is used at compile time to generate [Component].
type _ _Component

// A Component is a sum type or tagged union.
// Actually, we can generate different flavors, so that Go makes fun for modelling business stuff.
type Component struct {
	ordinal int
	value   any
}

func (e Component) Unwrap() any {
	return e.value
}

func (e Component) Ordinal() int {
	return e.ordinal
}

func (e Component) Valid() bool {
	return e.ordinal > 0
}

// Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case.
func (e Component) Switch(onButton func(Button), onTextField func(TextField), onText func(Text), onChapter func(Chapter), _onDefault func(any)) {
	switch e.ordinal {
	case 1:
		if onButton != nil {
			onButton(e.value.(Button))
			return
		}
	case 2:
		if onTextField != nil {
			onTextField(e.value.(TextField))
			return
		}
	case 3:
		if onText != nil {
			onText(e.value.(Text))
			return
		}
	case 4:
		if onChapter != nil {
			onChapter(e.value.(Chapter))
			return
		}
	}

	if _onDefault != nil {
		_onDefault(e.value)
	}
}

func (e Component) AsButton() (Button, bool) {
	var zero Button
	if e.ordinal == 1 {
		return e.value.(Button), true
	}

	return zero, false
}

func (e Component) WithButton(v Button) Component {
	e.ordinal = 1
	e.value = v
	return e
}

func (e Component) AsTextField() (TextField, bool) {
	var zero TextField
	if e.ordinal == 2 {
		return e.value.(TextField), true
	}

	return zero, false
}

func (e Component) WithTextField(v TextField) Component {
	e.ordinal = 2
	e.value = v
	return e
}

func (e Component) AsText() (Text, bool) {
	var zero Text
	if e.ordinal == 3 {
		return e.value.(Text), true
	}

	return zero, false
}

func (e Component) WithText(v Text) Component {
	e.ordinal = 3
	e.value = v
	return e
}

func (e Component) AsChapter() (Chapter, bool) {
	var zero Chapter
	if e.ordinal == 4 {
		return e.value.(Chapter), true
	}

	return zero, false
}

func (e Component) WithChapter(v Chapter) Component {
	e.ordinal = 4
	e.value = v
	return e
}

func (e Component) MarshalJSON() ([]byte, error) {
	if e.ordinal == 0 {
		return nil, fmt.Errorf("marshalling a zero value is not allowed")
	}

	// note, that by definition, this kind of encoding does not work with union types which evaluates to null, arrays or primitives.
	// Chose adjacent encoding instead.
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}

	switch e.ordinal {
	case 1:
		return json.Marshal(adjacentlyTagged[Button]{
			Type:  "helloBtn",
			Value: e.value.(Button),
		})
	case 2:
		return json.Marshal(adjacentlyTagged[TextField]{
			Type:  "ATef",
			Value: e.value.(TextField),
		})
	case 3:
		return json.Marshal(adjacentlyTagged[Text]{
			Type:  "str",
			Value: e.value.(Text),
		})
	case 4:
		return json.Marshal(adjacentlyTagged[Chapter]{
			Type:  "Chappy",
			Value: e.value.(Chapter),
		})
	default:
		return nil, fmt.Errorf("unknown type ordinal variant '%d'", e.ordinal)
	}
}

func (e *Component) UnmarshalJSON(bytes []byte) error {
	typeOnly := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}
	switch typeOnly.Type {
	case "helloBtn":
		var value adjacentlyTagged[Button]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Button'")
		}
		e.ordinal = 1
		e.value = value.Value
	case "ATef":
		var value adjacentlyTagged[TextField]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'TextField'")
		}
		e.ordinal = 2
		e.value = value.Value
	case "str":
		var value adjacentlyTagged[Text]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Text'")
		}
		e.ordinal = 3
		e.value = value.Value
	case "Chappy":
		var value adjacentlyTagged[Chapter]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Chapter'")
		}
		e.ordinal = 4
		e.value = value.Value
	default:
		return fmt.Errorf("unknown type variant name '%s'", typeOnly.Type)
	}

	return nil
}

func MatchComponent[R any](e Component, onButton func(Button) R, onTextField func(TextField) R, onText func(Text) R, onChapter func(Chapter) R, _onDefault func(any) R) R {
	if _onDefault == nil {
		panic(`missing default match: cannot guarantee exhaustive matching`)
	}

	switch e.ordinal {
	case 1:
		if onButton != nil {
			return onButton(e.value.(Button))
		}
	case 2:
		if onTextField != nil {
			return onTextField(e.value.(TextField))
		}
	case 3:
		if onText != nil {
			return onText(e.value.(Text))
		}
	case 4:
		if onChapter != nil {
			return onChapter(e.value.(Chapter))
		}
	}

	return _onDefault(e.value)
}

// This variable is declared to let linters know, that [_privateComponent] is used at compile time to generate [privateComponent].
type _ _privateComponent
type privateComponent struct {
	ordinal int
	value   any
}

func (e privateComponent) Unwrap() any {
	return e.value
}

func (e privateComponent) Ordinal() int {
	return e.ordinal
}

func (e privateComponent) Valid() bool {
	return e.ordinal > 0
}

// Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case.
func (e privateComponent) Switch(onButton func(Button), onTextField func(TextField), onText func(Text), onChapter func(Chapter), _onDefault func(any)) {
	switch e.ordinal {
	case 1:
		if onButton != nil {
			onButton(e.value.(Button))
			return
		}
	case 2:
		if onTextField != nil {
			onTextField(e.value.(TextField))
			return
		}
	case 3:
		if onText != nil {
			onText(e.value.(Text))
			return
		}
	case 4:
		if onChapter != nil {
			onChapter(e.value.(Chapter))
			return
		}
	}

	if _onDefault != nil {
		_onDefault(e.value)
	}
}

func (e privateComponent) AsButton() (Button, bool) {
	var zero Button
	if e.ordinal == 1 {
		return e.value.(Button), true
	}

	return zero, false
}

func (e privateComponent) WithButton(v Button) privateComponent {
	e.ordinal = 1
	e.value = v
	return e
}

func (e privateComponent) AsTextField() (TextField, bool) {
	var zero TextField
	if e.ordinal == 2 {
		return e.value.(TextField), true
	}

	return zero, false
}

func (e privateComponent) WithTextField(v TextField) privateComponent {
	e.ordinal = 2
	e.value = v
	return e
}

func (e privateComponent) AsText() (Text, bool) {
	var zero Text
	if e.ordinal == 3 {
		return e.value.(Text), true
	}

	return zero, false
}

func (e privateComponent) WithText(v Text) privateComponent {
	e.ordinal = 3
	e.value = v
	return e
}

func (e privateComponent) AsChapter() (Chapter, bool) {
	var zero Chapter
	if e.ordinal == 4 {
		return e.value.(Chapter), true
	}

	return zero, false
}

func (e privateComponent) WithChapter(v Chapter) privateComponent {
	e.ordinal = 4
	e.value = v
	return e
}

func (e privateComponent) MarshalJSON() ([]byte, error) {
	if e.ordinal == 0 {
		return nil, fmt.Errorf("marshalling a zero value is not allowed")
	}

	// note, that by definition, this kind of encoding does not work with union types which evaluates to null, arrays or primitives.
	// Chose adjacent encoding instead.
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}

	switch e.ordinal {
	case 1:
		return json.Marshal(adjacentlyTagged[Button]{
			Type:  "Button",
			Value: e.value.(Button),
		})
	case 2:
		return json.Marshal(adjacentlyTagged[TextField]{
			Type:  "TextField",
			Value: e.value.(TextField),
		})
	case 3:
		return json.Marshal(adjacentlyTagged[Text]{
			Type:  "Text",
			Value: e.value.(Text),
		})
	case 4:
		return json.Marshal(adjacentlyTagged[Chapter]{
			Type:  "Chapter",
			Value: e.value.(Chapter),
		})
	default:
		return nil, fmt.Errorf("unknown type ordinal variant '%d'", e.ordinal)
	}
}

func (e *privateComponent) UnmarshalJSON(bytes []byte) error {
	typeOnly := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}
	switch typeOnly.Type {
	case "Button":
		var value adjacentlyTagged[Button]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Button'")
		}
		e.ordinal = 1
		e.value = value.Value
	case "TextField":
		var value adjacentlyTagged[TextField]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'TextField'")
		}
		e.ordinal = 2
		e.value = value.Value
	case "Text":
		var value adjacentlyTagged[Text]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Text'")
		}
		e.ordinal = 3
		e.value = value.Value
	case "Chapter":
		var value adjacentlyTagged[Chapter]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Chapter'")
		}
		e.ordinal = 4
		e.value = value.Value
	default:
		return fmt.Errorf("unknown type variant name '%s'", typeOnly.Type)
	}

	return nil
}

func MatchprivateComponent[R any](e privateComponent, onButton func(Button) R, onTextField func(TextField) R, onText func(Text) R, onChapter func(Chapter) R, _onDefault func(any) R) R {
	if _onDefault == nil {
		panic(`missing default match: cannot guarantee exhaustive matching`)
	}

	switch e.ordinal {
	case 1:
		if onButton != nil {
			return onButton(e.value.(Button))
		}
	case 2:
		if onTextField != nil {
			return onTextField(e.value.(TextField))
		}
	case 3:
		if onText != nil {
			return onText(e.value.(Text))
		}
	case 4:
		if onChapter != nil {
			return onChapter(e.value.(Chapter))
		}
	}

	return _onDefault(e.value)
}

// This variable is declared to let linters know, that [_ÜmlautComponent] is used at compile time to generate [ÜmlautComponent].
type _ _ÜmlautComponent
type ÜmlautComponent struct {
	ordinal int
	value   any
}

func (e ÜmlautComponent) Unwrap() any {
	return e.value
}

func (e ÜmlautComponent) Ordinal() int {
	return e.ordinal
}

func (e ÜmlautComponent) Valid() bool {
	return e.ordinal > 0
}

// Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case.
func (e ÜmlautComponent) Switch(onButton func(Button), onTextField func(TextField), onText func(Text), onChapter func(Chapter), onÄpfel func(Äpfel), _onDefault func(any)) {
	switch e.ordinal {
	case 1:
		if onButton != nil {
			onButton(e.value.(Button))
			return
		}
	case 2:
		if onTextField != nil {
			onTextField(e.value.(TextField))
			return
		}
	case 3:
		if onText != nil {
			onText(e.value.(Text))
			return
		}
	case 4:
		if onChapter != nil {
			onChapter(e.value.(Chapter))
			return
		}
	case 5:
		if onÄpfel != nil {
			onÄpfel(e.value.(Äpfel))
			return
		}
	}

	if _onDefault != nil {
		_onDefault(e.value)
	}
}

func (e ÜmlautComponent) AsButton() (Button, bool) {
	var zero Button
	if e.ordinal == 1 {
		return e.value.(Button), true
	}

	return zero, false
}

func (e ÜmlautComponent) WithButton(v Button) ÜmlautComponent {
	e.ordinal = 1
	e.value = v
	return e
}

func (e ÜmlautComponent) AsTextField() (TextField, bool) {
	var zero TextField
	if e.ordinal == 2 {
		return e.value.(TextField), true
	}

	return zero, false
}

func (e ÜmlautComponent) WithTextField(v TextField) ÜmlautComponent {
	e.ordinal = 2
	e.value = v
	return e
}

func (e ÜmlautComponent) AsText() (Text, bool) {
	var zero Text
	if e.ordinal == 3 {
		return e.value.(Text), true
	}

	return zero, false
}

func (e ÜmlautComponent) WithText(v Text) ÜmlautComponent {
	e.ordinal = 3
	e.value = v
	return e
}

func (e ÜmlautComponent) AsChapter() (Chapter, bool) {
	var zero Chapter
	if e.ordinal == 4 {
		return e.value.(Chapter), true
	}

	return zero, false
}

func (e ÜmlautComponent) WithChapter(v Chapter) ÜmlautComponent {
	e.ordinal = 4
	e.value = v
	return e
}

func (e ÜmlautComponent) AsÄpfel() (Äpfel, bool) {
	var zero Äpfel
	if e.ordinal == 5 {
		return e.value.(Äpfel), true
	}

	return zero, false
}

func (e ÜmlautComponent) WithÄpfel(v Äpfel) ÜmlautComponent {
	e.ordinal = 5
	e.value = v
	return e
}

func (e ÜmlautComponent) MarshalJSON() ([]byte, error) {
	if e.ordinal == 0 {
		return nil, fmt.Errorf("marshalling a zero value is not allowed")
	}

	// note, that by definition, this kind of encoding does not work with union types which evaluates to null, arrays or primitives.
	// Chose adjacent encoding instead.
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}

	switch e.ordinal {
	case 1:
		return json.Marshal(adjacentlyTagged[Button]{
			Type:  "Button",
			Value: e.value.(Button),
		})
	case 2:
		return json.Marshal(adjacentlyTagged[TextField]{
			Type:  "TextField",
			Value: e.value.(TextField),
		})
	case 3:
		return json.Marshal(adjacentlyTagged[Text]{
			Type:  "Text",
			Value: e.value.(Text),
		})
	case 4:
		return json.Marshal(adjacentlyTagged[Chapter]{
			Type:  "Chapter",
			Value: e.value.(Chapter),
		})
	case 5:
		return json.Marshal(adjacentlyTagged[Äpfel]{
			Type:  "Äpfel",
			Value: e.value.(Äpfel),
		})
	default:
		return nil, fmt.Errorf("unknown type ordinal variant '%d'", e.ordinal)
	}
}

func (e *ÜmlautComponent) UnmarshalJSON(bytes []byte) error {
	typeOnly := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}
	switch typeOnly.Type {
	case "Button":
		var value adjacentlyTagged[Button]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Button'")
		}
		e.ordinal = 1
		e.value = value.Value
	case "TextField":
		var value adjacentlyTagged[TextField]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'TextField'")
		}
		e.ordinal = 2
		e.value = value.Value
	case "Text":
		var value adjacentlyTagged[Text]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Text'")
		}
		e.ordinal = 3
		e.value = value.Value
	case "Chapter":
		var value adjacentlyTagged[Chapter]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Chapter'")
		}
		e.ordinal = 4
		e.value = value.Value
	case "Äpfel":
		var value adjacentlyTagged[Äpfel]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Äpfel'")
		}
		e.ordinal = 5
		e.value = value.Value
	default:
		return fmt.Errorf("unknown type variant name '%s'", typeOnly.Type)
	}

	return nil
}

func MatchÜmlautComponent[R any](e ÜmlautComponent, onButton func(Button) R, onTextField func(TextField) R, onText func(Text) R, onChapter func(Chapter) R, onÄpfel func(Äpfel) R, _onDefault func(any) R) R {
	if _onDefault == nil {
		panic(`missing default match: cannot guarantee exhaustive matching`)
	}

	switch e.ordinal {
	case 1:
		if onButton != nil {
			return onButton(e.value.(Button))
		}
	case 2:
		if onTextField != nil {
			return onTextField(e.value.(TextField))
		}
	case 3:
		if onText != nil {
			return onText(e.value.(Text))
		}
	case 4:
		if onChapter != nil {
			return onChapter(e.value.(Chapter))
		}
	case 5:
		if onÄpfel != nil {
			return onÄpfel(e.value.(Äpfel))
		}
	}

	return _onDefault(e.value)
}

// This variable is declared to let linters know, that [_ExampleType] is used at compile time to generate [ExampleType].
type _ _ExampleType
type ExampleType struct {
	ordinal int
	value   any
}

func (e ExampleType) Unwrap() any {
	return e.value
}

func (e ExampleType) Ordinal() int {
	return e.ordinal
}

func (e ExampleType) Valid() bool {
	return e.ordinal > 0
}

// Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case.
func (e ExampleType) Switch(onButton func(Button), onTextField func(TextField), _onDefault func(any)) {
	switch e.ordinal {
	case 1:
		if onButton != nil {
			onButton(e.value.(Button))
			return
		}
	case 2:
		if onTextField != nil {
			onTextField(e.value.(TextField))
			return
		}
	}

	if _onDefault != nil {
		_onDefault(e.value)
	}
}

func (e ExampleType) AsButton() (Button, bool) {
	var zero Button
	if e.ordinal == 1 {
		return e.value.(Button), true
	}

	return zero, false
}

func (e ExampleType) WithButton(v Button) ExampleType {
	e.ordinal = 1
	e.value = v
	return e
}

func (e ExampleType) AsTextField() (TextField, bool) {
	var zero TextField
	if e.ordinal == 2 {
		return e.value.(TextField), true
	}

	return zero, false
}

func (e ExampleType) WithTextField(v TextField) ExampleType {
	e.ordinal = 2
	e.value = v
	return e
}

func (e ExampleType) MarshalJSON() ([]byte, error) {
	if e.ordinal == 0 {
		return nil, fmt.Errorf("marshalling a zero value is not allowed")
	}

	// note, that by definition, this kind of encoding does not work with union types which evaluates to null, arrays or primitives.
	// Chose adjacent encoding instead.
	buf, err := json.Marshal(e.value)
	if err != nil {
		return nil, err
	}
	var prefix []byte

	switch e.ordinal {
	case 1:
		prefix = []byte(`{"type":"Button"`)
	case 2:
		prefix = []byte(`{"type":"TextField"`)
	}

	if len(buf) > 2 {
		// we expect an empty object like {} or at least an object with a single attribute, which requires a , separator
		prefix = append(prefix, ',')
	}
	buf = append(buf[1:], prefix...)
	copy(buf[len(prefix):], buf)
	copy(buf, prefix)

	return buf, nil
}

func (e *ExampleType) UnmarshalJSON(bytes []byte) error {
	typeOnly := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}
	switch typeOnly.Type {
	case "Button":
		var value Button
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Button'")
		}
		e.ordinal = 1
		e.value = value
	case "TextField":
		var value TextField
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'TextField'")
		}
		e.ordinal = 2
		e.value = value
	default:
		return fmt.Errorf("unknown type variant name '%s'", typeOnly.Type)
	}

	return nil
}

func MatchExampleType[R any](e ExampleType, onButton func(Button) R, onTextField func(TextField) R, _onDefault func(any) R) R {
	if _onDefault == nil {
		panic(`missing default match: cannot guarantee exhaustive matching`)
	}

	switch e.ordinal {
	case 1:
		if onButton != nil {
			return onButton(e.value.(Button))
		}
	case 2:
		if onTextField != nil {
			return onTextField(e.value.(TextField))
		}
	}

	return _onDefault(e.value)
}

// FruchtLike defines a marker method for the polymorphic interface modelling of any member of the sum type.
type FruchtLike interface {
	isFrucht()
}

func (_ Birne) isFrucht() {
}
func (_ Kirsche) isFrucht() {
}

// This variable is declared to let linters know, that [_Frucht] is used at compile time to generate [Frucht].
type _ _Frucht
type Frucht struct {
	ordinal int
	value   any
}

func (e Frucht) Unwrap() any {
	return e.value
}

func (e Frucht) Ordinal() int {
	return e.ordinal
}

func (e Frucht) Valid() bool {
	return e.ordinal > 0
}

// Switch provides an exhaustive and type safe closure callback mechanic. Nil callbacks are allowed. Unmatched branches are delegated into a default case.
func (e Frucht) Switch(onBirne func(Birne), onKirsche func(Kirsche), _onDefault func(any)) {
	switch e.ordinal {
	case 1:
		if onBirne != nil {
			onBirne(e.value.(Birne))
			return
		}
	case 2:
		if onKirsche != nil {
			onKirsche(e.value.(Kirsche))
			return
		}
	}

	if _onDefault != nil {
		_onDefault(e.value)
	}
}

func (e Frucht) AsBirne() (Birne, bool) {
	var zero Birne
	if e.ordinal == 1 {
		return e.value.(Birne), true
	}

	return zero, false
}

func (e Frucht) AsKirsche() (Kirsche, bool) {
	var zero Kirsche
	if e.ordinal == 2 {
		return e.value.(Kirsche), true
	}

	return zero, false
}

func (e Frucht) MarshalJSON() ([]byte, error) {
	if e.ordinal == 0 {
		return nil, fmt.Errorf("marshalling a zero value is not allowed")
	}

	// note, that by definition, this kind of encoding does not work with union types which evaluates to null, arrays or primitives.
	// Chose adjacent encoding instead.
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}

	switch e.ordinal {
	case 1:
		return json.Marshal(adjacentlyTagged[Birne]{
			Type:  "Birne",
			Value: e.value.(Birne),
		})
	case 2:
		return json.Marshal(adjacentlyTagged[Kirsche]{
			Type:  "Kirsche",
			Value: e.value.(Kirsche),
		})
	default:
		return nil, fmt.Errorf("unknown type ordinal variant '%d'", e.ordinal)
	}
}

func (e *Frucht) UnmarshalJSON(bytes []byte) error {
	typeOnly := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(bytes, &typeOnly); err != nil {
		return err
	}
	type adjacentlyTagged[T any] struct {
		Type  string `json:"type"`
		Value T      `json:"content"`
	}
	switch typeOnly.Type {
	case "Birne":
		var value adjacentlyTagged[Birne]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Birne'")
		}
		e.ordinal = 1
		e.value = value.Value
	case "Kirsche":
		var value adjacentlyTagged[Kirsche]
		if err := json.Unmarshal(bytes, &value); err != nil {
			return fmt.Errorf("cannot unmarshal variant 'Kirsche'")
		}
		e.ordinal = 2
		e.value = value.Value
	default:
		return fmt.Errorf("unknown type variant name '%s'", typeOnly.Type)
	}

	return nil
}

func NewFrucht(obj FruchtLike) Frucht {
	u := Frucht{}
	switch obj := obj.(type) {
	case Birne:
		u.ordinal = 1
		u.value = obj
	case Kirsche:
		u.ordinal = 2
		u.value = obj
	default:
		panic(fmt.Errorf("invalid value: %T", obj))
	}
	return u
}
func MatchFrucht[R any](e Frucht, onBirne func(Birne) R, onKirsche func(Kirsche) R, _onDefault func(any) R) R {
	if _onDefault == nil {
		panic(`missing default match: cannot guarantee exhaustive matching`)
	}

	switch e.ordinal {
	case 1:
		if onBirne != nil {
			return onBirne(e.value.(Birne))
		}
	case 2:
		if onKirsche != nil {
			return onKirsche(e.value.(Kirsche))
		}
	}

	return _onDefault(e.value)
}
