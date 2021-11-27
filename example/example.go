package example

import (
	"errors"
	"fmt"
)

func Foo(prefix string) (result []string, err error) {
	if prefix == "foo" {
		return nil, errors.New("bad prefix")
	}

	for i := 0; i < 10; i++ {
		result = append(result, fmt.Sprintf("%s-%d", prefix, i))
	}
	return
}
