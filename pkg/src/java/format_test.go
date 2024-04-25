package java

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	src0 := `
package myTest ;

// bla
public class App
{
/**
Another comment
another line

another lang
*/
public static void  main(args ...String)
{
}
}
`

	src, err := Format([]byte(src0))
	if err != nil {
		t.Fatal(err, string(src))
	}
	fmt.Println(string(src))

	src1 := `invalid stuff`
	if _, err := Format([]byte(src1)); err == nil {
		t.Fatal("should have failed")
	}

}
