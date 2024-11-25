package iprange

import "testing"

func FuzzParseList(f *testing.F) {
	f.Add("10.0.0.1")
	f.Add("10.0.0.0/24")
	f.Add("10.0.0.*")
	f.Add("10.0.0.1-10")
	f.Add("10.0.0.-10")
	f.Add("10.0.0.256=257")
	f.Add("10.0.0.256-257")
	f.Add("1453264632.37423598794.497593750754.4")
	f.Add("10.0.0.2561235467876543456787654567875678768788")
	f.Add("sdfghjk")
	f.Add("10.0.0.3/32")
	f.Add("10.0.0.3/33")
	f.Add("10.0.0.3/31")

	// Основной фузз-тест
	f.Fuzz(func(t *testing.T, input string) {
		_, err := ParseList(input)
		if err != nil {
			t.Errorf("unexpected error: %v, input: %s", err, input)
		}
	})
}
