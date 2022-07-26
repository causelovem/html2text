# html2text
Simple html to text converter.

## How to install
```bash
go get github.com/causelovem/html2text
```

## How to use
```go
package main

import (
	"fmt"

	"github.com/causelovem/html2text"
)

func main() {
	html := `<div><h1>some text</h1></div> <p>to be <b>cleared</b> </br> from &quot;html&quot; &lt;tags&gt;</p>`

	text := html2text.HTML2Text(html)

	fmt.Println(text)
}
```
