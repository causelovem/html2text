package html2Text

import "testing"

func TestHTML2Text(t *testing.T) {
	type args struct {
		htmlString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0", args{""}, ""},
		{"1", args{"<div>"}, ""},
		{"2", args{"<p></p>"}, ""},
		{"3", args{"<p>text </p>"}, "text \n"},
		{"4", args{"<p>some new line <br>here and <br />here also</p>"}, "some new line \nhere and \nhere also\n"},
		{"5", args{"<h1>text</h1> <div> <noscript> <p> text inside noscript <p>  </noscript> </div>"}, "text\n text inside noscript \n"},
		{"6", args{"<h1>header</h1> <script> func some func() {some + code = return 1010} </script> <div> text </div>"}, "header\n text "},
		{"7", args{"<p>text1 </p> <p> text2 </p>"}, "text1 \n text2 \n"},
		{"8", args{"<div> <div> some text</div> <div> another text </div> </div>"}, " some text another text "},
		{"9", args{"<div> <div>some &amp; text</div> <div>another text</div> </div>"}, "some & text another text"},
		{"10", args{"<div> <p>some &lt;&gt; text</p> <div>another text</div> </div>"}, "some <> text\nanother text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTML2Text(tt.args.htmlString); got != tt.want {
				t.Errorf("HTML2Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
