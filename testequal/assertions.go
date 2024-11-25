//go:build !solution

package testequal

import (
	"reflect"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if reflect.TypeOf(expected) == reflect.TypeOf(actual) && reflect.TypeOf(expected) == reflect.TypeOf(struct{}{}) {
		return false
	}
	switch expected := expected.(type) {
	case []int:
		if actualSlice, ok := actual.([]int); ok {
			if expected == nil || actualSlice == nil {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
			if len(expected) != len(actualSlice) {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
			for i, expVal := range expected {
				if actVal := actualSlice[i]; expVal != actVal {
					t.Helper()
					if len(msgAndArgs) > 0 {
						if format, ok := msgAndArgs[0].(string); ok {
							t.Errorf(format, msgAndArgs[1:]...)
						}
					} else {
						t.Errorf("")
					}
					return false
				}
			}
			return true
		}
	case []byte:
		if actualSlice, ok := actual.([]byte); ok {
			if expected == nil || actualSlice == nil {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
			if len(expected) != len(actualSlice) {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
			for i, expVal := range expected {
				if actVal := actualSlice[i]; expVal != actVal {
					t.Helper()
					if len(msgAndArgs) > 0 {
						if format, ok := msgAndArgs[0].(string); ok {
							t.Errorf(format, msgAndArgs[1:]...)
						}
					} else {
						t.Errorf("")
					}
					return false
				}
			}
			return true
		}
	case map[string]string:
		if actualMap, ok := actual.(map[string]string); ok {
			if expected == nil || actualMap == nil {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
			if len(expected) != len(actualMap) {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
			for key, expVal := range expected {
				if actVal, found := actualMap[key]; !found || expVal != actVal {
					t.Helper()
					if len(msgAndArgs) > 0 {
						if format, ok := msgAndArgs[0].(string); ok {
							t.Errorf(format, msgAndArgs[1:]...)
						}
					} else {
						t.Errorf("")
					}
					return false
				}
			}
			return true
		}
	}
	if expected != actual {
		t.Helper()
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				t.Errorf(format, msgAndArgs[1:]...)
			}
		} else {
			t.Errorf("")
		}
		return false
	}

	return true
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if reflect.TypeOf(expected) == reflect.TypeOf(actual) && reflect.TypeOf(expected) == reflect.TypeOf(struct{}{}) {
		return false
	}
	if reflect.TypeOf(expected) == reflect.TypeOf([]int{}) {
		if actualSlice, ok := actual.([]int); ok {
			if (actualSlice == nil && expected.([]int) != nil) || (actualSlice != nil && expected.([]int) == nil) {
				return true
			}
			if len(expected.([]int)) != len(actualSlice) {
				return true
			}
			equal := true
			for i, expVal := range expected.([]int) {
				if actVal := actualSlice[i]; expVal != actVal {
					equal = false
				}
			}
			if !equal {
				return true
			} else {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
		}
	}
	if reflect.TypeOf(expected) == reflect.TypeOf([]byte{}) {
		if actualSlice, ok := actual.([]byte); ok {
			if (actualSlice == nil && expected.([]byte) != nil) || (actualSlice != nil && expected.([]byte) == nil) {
				return true
			}
			if len(expected.([]byte)) != len(actualSlice) {
				return true
			}
			equal := true
			for i, expVal := range expected.([]byte) {
				if actVal := actualSlice[i]; expVal != actVal {
					equal = false
				}
			}
			if !equal {
				return true
			} else {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
		}
	}
	if reflect.TypeOf(expected) == reflect.TypeOf(map[string]string{}) {
		if actualSlice, ok := actual.(map[string]string); ok {
			if (actualSlice == nil && expected.(map[string]string) != nil) || (actualSlice != nil && expected.(map[string]string) == nil) {
				return true
			}
			if len(expected.(map[string]string)) != len(actualSlice) {
				return true
			}
			equal := true
			for i, expVal := range expected.(map[string]string) {
				if actVal := actualSlice[i]; expVal != actVal {
					equal = false
				}
			}
			if !equal {
				return true
			} else {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				return false
			}
		}
	}
	if expected == actual {
		t.Helper()
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				t.Errorf(format, msgAndArgs[1:]...)
			}
		} else {
			t.Errorf("")
		}
		return false
	}

	return true
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	switch expected := expected.(type) {
	case []int:
		if actualSlice, ok := actual.([]int); ok {
			if expected == nil || actualSlice == nil {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()

			}
			if len(expected) != len(actualSlice) {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
			for i, expVal := range expected {
				if actVal := actualSlice[i]; expVal != actVal {
					t.Helper()
					if len(msgAndArgs) > 0 {
						if format, ok := msgAndArgs[0].(string); ok {
							t.Errorf(format, msgAndArgs[1:]...)
						}
					} else {
						t.Errorf("")
					}
					t.FailNow()
				}
			}
			return
		}
	case []byte:
		if actualSlice, ok := actual.([]byte); ok {
			if expected == nil || actualSlice == nil {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
			if len(expected) != len(actualSlice) {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
			for i, expVal := range expected {
				if actVal := actualSlice[i]; expVal != actVal {
					t.Helper()
					if len(msgAndArgs) > 0 {
						if format, ok := msgAndArgs[0].(string); ok {
							t.Errorf(format, msgAndArgs[1:]...)
						}
					} else {
						t.Errorf("")
					}
					t.FailNow()
				}
			}
			return
		}
	case map[string]string:
		if actualMap, ok := actual.(map[string]string); ok {
			if expected == nil || actualMap == nil {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
			if len(expected) != len(actualMap) {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
			for key, expVal := range expected {
				if actVal, found := actualMap[key]; !found || expVal != actVal {
					t.Helper()
					if len(msgAndArgs) > 0 {
						if format, ok := msgAndArgs[0].(string); ok {
							t.Errorf(format, msgAndArgs[1:]...)
						}
					} else {
						t.Errorf("")
					}
					t.FailNow()
				}
			}
			return
		}
	}
	if expected != actual {
		t.Helper()
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				t.Errorf(format, msgAndArgs[1:]...)
			}
		} else {
			t.Errorf("")
		}
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if reflect.TypeOf(expected) == reflect.TypeOf(actual) && reflect.TypeOf(expected) == reflect.TypeOf(struct{}{}) {
		return
	}
	if reflect.TypeOf(expected) == reflect.TypeOf([]int{}) {
		if actualSlice, ok := actual.([]int); ok {
			if (actualSlice == nil && expected.([]int) != nil) || (actualSlice != nil && expected.([]int) == nil) {
				return
			}
			if len(expected.([]int)) != len(actualSlice) {
				return
			}
			equal := true
			for i, expVal := range expected.([]int) {
				if actVal := actualSlice[i]; expVal != actVal {
					equal = false
				}
			}
			if !equal {
				return
			} else {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
		}
	}
	if reflect.TypeOf(expected) == reflect.TypeOf([]byte{}) {
		if actualSlice, ok := actual.([]byte); ok {
			if (actualSlice == nil && expected.([]byte) != nil) || (actualSlice != nil && expected.([]byte) == nil) {
				return
			}
			if len(expected.([]byte)) != len(actualSlice) {
				return
			}
			equal := true
			for i, expVal := range expected.([]byte) {
				if actVal := actualSlice[i]; expVal != actVal {
					equal = false
				}
			}
			if !equal {
				return
			} else {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
		}
	}
	if reflect.TypeOf(expected) == reflect.TypeOf(map[string]string{}) {
		if actualSlice, ok := actual.(map[string]string); ok {
			if (actualSlice == nil && expected.(map[string]string) != nil) || (actualSlice != nil && expected.(map[string]string) == nil) {
				return
			}
			if len(expected.(map[string]string)) != len(actualSlice) {
				return
			}
			equal := true
			for i, expVal := range expected.(map[string]string) {
				if actVal := actualSlice[i]; expVal != actVal {
					equal = false
				}
			}
			if !equal {
				return
			} else {
				t.Helper()
				if len(msgAndArgs) > 0 {
					if format, ok := msgAndArgs[0].(string); ok {
						t.Errorf(format, msgAndArgs[1:]...)
					}
				} else {
					t.Errorf("")
				}
				t.FailNow()
			}
		}
	}
	if expected == actual {
		t.Helper()
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				t.Errorf(format, msgAndArgs[1:]...)
			}
		} else {
			t.Errorf("")
		}
		t.FailNow()
	}
}
