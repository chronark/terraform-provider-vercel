package utility

// Difference returns a slice containing all of the elements of 'a' that aren't present in 'b'.
//
// See https://siongui.github.io/2018/03/14/go-set-difference-of-two-arrays/
func Difference(a, b []string) (diff []string) {
  m := make(map[string]bool)

  for _, item := range b {
    m[item] = true
  }

  for _, item := range a {
    if _, ok := m[item]; !ok {
      diff = append(diff, item)
    }
  }

  return diff
}
