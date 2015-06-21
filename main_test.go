package main

import (
	"bytes"
	"strings"
	"testing"
)

var tests = []struct {
	in   string
	want string
}{
	// 	{
	// 		in:   ``,
	// 		want: ``,
	// 	},
	// 	{
	// 		in: `
	// // this comment is longer than 70 lines and is not followed by a comment itself
	// `,
	// 		want: `
	// // this comment is longer than 70 lines and is not followed by a
	// // comment itself
	// `,
	// 	},
	// 	{
	// 		in: `
	// // this comment is longer than 70 lines and is followed by a comment itself,
	// // so this bit should be prefixed
	// `,
	// 		want: `
	// // this comment is longer than 70 lines and is followed by a comment
	// // itself, so this bit should be prefixed
	// `,
	// 	},
	// 	{
	// 		in: `
	// // this comment is longer than 70 lines and is followed by a comment itself, and is in fact many times the length of the maximum length that is allowed, so this should be truncated or chomped at least three times,
	// // so this bit should be prefixed
	// this is some code, lets say
	// `,
	// 		want: `
	// // this comment is longer than 70 lines and is followed by a comment
	// // itself, and is in fact many times the length of the maximum length
	// // that is allowed, so this should be truncated or chomped at least
	// // three times, so this bit should be prefixed
	// this is some code, lets say
	// `,
	// 	},
	// 	{
	// 		in: `
	// // this comment is longer than 70 lines and is followed by a comment itself, and is in fact many times the length of the maximum length that is allowed, so this should be truncated or chomped at least three times,
	// // so this bit should be prefixed
	// this is some code, lets say

	// type Something struct {
	//     ImAField    string // i'm a comment on the field, and i'm longer than 70 chars
	//     ImAField    string // i'm a short comment on the field
	//     ImAField    string // i'm a comment on the field, and i'm longer than 70 chars, like really a lot longer, even I should be wrapped

	//     // i'm a comment *before* the field, and i'm longer than 70 chars, like really a lot longer, even I should be wrapped
	//     ImAField    string
	// }
	// `,
	// 		want: `
	// // this comment is longer than 70 lines and is followed by a comment
	// // itself, and is in fact many times the length of the maximum length
	// // that is allowed, so this should be truncated or chomped at least
	// // three times, so this bit should be prefixed
	// this is some code, lets say

	// type Something struct {
	//     // i'm a comment on the field, and i'm longer than 70 chars
	//     ImAField    string
	//     ImAField    string // i'm a short comment on the field
	//     // i'm a comment on the field, and i'm longer than 70 chars, like
	//     // really a lot longer, even I should be wrapped
	//     ImAField    string

	//     // i'm a comment *before* the field, and i'm longer than 70 chars,
	//     // like really a lot longer, even I should be wrapped
	//     ImAField    string
	// }
	// `,
	// 	},

	// 	{
	// 		in: `
	// // this comment is longer than 70 lines and is followed by a comment itself,
	// // so this bit should be prefixed
	// this is some code, lets say

	// type Something struct {
	//     ImAField    string // i'm a comment on the field, and i'm longer than 70 chars
	//     ImAField    string // i'm a short comment on the field
	//     ImAField    string // i'm a comment on the field, and i'm longer than 70 chars, like really a lot longer, even I should be wrapped

	//     // i'm a comment *before* the field, and i'm longer than 70 chars, like really a lot longer, even I should be wrapped
	//     ImAField    string

	//     // im just a lone indented comment, this comment is longer than 70 lines and is not followed by a comment itself, so this bit should be prefixed
	// }
	// `,
	// 		want: `
	// // this comment is longer than 70 lines and is followed by a comment
	// // itself, so this bit should be prefixed
	// this is some code, lets say

	// type Something struct {
	//     // i'm a comment on the field, and i'm longer than 70 chars
	//     ImAField    string
	//     ImAField    string // i'm a short comment on the field
	//     // i'm a comment on the field, and i'm longer than 70 chars, like
	//     // really a lot longer, even I should be wrapped
	//     ImAField    string

	//     // i'm a comment *before* the field, and i'm longer than 70 chars,
	//     // like really a lot longer, even I should be wrapped
	//     ImAField    string

	//     // im just a lone indented comment, this comment is longer than 70
	//     // lines and is not followed by a comment itself, so this bit
	//     // should be prefixed
	// }
	// `,
	// 	},

	{
		in: `
       // AsyncClose triggers a shutdown of the producer, flushing any messages it may have
       // buffered. The shutdown has completed when both the Errors and Successes channels
       // have been closed. When calling AsyncClose, you *must* continue to read from those
`,
		want: `
       // AsyncClose triggers a shutdown of the producer, flushing any
       // messages it may have buffered. The shutdown has completed
       // when both the Errors and Successes channels have been
       // closed. When calling AsyncClose, you *must* continue to
       // read from those channels in order to drain the results of
       // any messages in flight.
        `,
	},
}

func TestCanFixComments(t *testing.T) {

	for i, tt := range tests {
		t.Logf("test #%d", i+1)

		out := bytes.NewBuffer(nil)
		_ = "breakpoint"
		if _, err := wrapComments(out, strings.NewReader(tt.in)); err != nil {
			t.Fatal(err)
		}
		want := tt.want
		got := out.String()
		if want != got {
			t.Errorf("want=\n%s", want)
			t.Errorf(" got=\n%s", got)
		}
	}
}
