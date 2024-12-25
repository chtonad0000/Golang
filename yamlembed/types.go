package yamlembed

import (
	"fmt"
	_ "gopkg.in/yaml.v2"
	"strings"
)

// Foo handles its custom (de)serialization logic.
type Foo struct {
	A string `yaml:"aa"`
	p int64  // private field, ignored
}

func (f *Foo) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		A string `yaml:"aa"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	f.A = raw.A
	return nil
}

func (f Foo) MarshalYAML() (interface{}, error) {
	return struct {
		A string `yaml:"aa"`
	}{
		A: f.A,
	}, nil
}

// Bar handles its custom (de)serialization logic.
type Bar struct {
	I      int64    `yaml:"i,omitempty"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"` // Не сериализуется в YAML
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f,omitempty"`
}

// UnmarshalYAML для Bar, чтобы обрабатывать приватное поле UpperB и сбрасывать I
func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		B  string   `yaml:"b"`
		OI []string `yaml:"oi,omitempty"`
		F  []any    `yaml:"f,omitempty"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	b.B = raw.B
	b.UpperB = strings.ToUpper(raw.B) // Устанавливаем UpperB на основе B
	b.OI = raw.OI
	b.F = raw.F
	b.I = 0 // Игнорируем поле I из YAML
	return nil
}

// MarshalYAML для Bar с кастомной сериализацией для F
func (b Bar) MarshalYAML() (interface{}, error) {
	// Преобразуем F как строку с квадратными скобками
	f := fmt.Sprintf("%v", b.F)
	return struct {
		B  string   `yaml:"b"`
		I  int64    `yaml:"i,omitempty"`
		OI []string `yaml:"oi,omitempty"`
		F  string   `yaml:"f,omitempty"` // f будет строкой с квадратными скобками
	}{
		B:  b.B,
		I:  b.I,
		OI: b.OI,
		F:  f, // Преобразуем в строку с квадратными скобками
	}, nil
}

// Baz combines Foo and Bar with inline embedding.
type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}
