package min

// github.com/tdewolff/minify binding for glub.
// No Configuration required.

import (
	"sync"
   "bytes"
   "log"
	"github.com/tdewolff/minify"
   "github.com/tdewolff/minify/css"
   "github.com/tdewolff/minify/html"
   "github.com/tdewolff/minify/js"
	"github.com/omeid/slurp"
)

func Min(c *slurp.C, mediatype string) slurp.Stage {
	return func(in <-chan slurp.File, out chan<- slurp.File) {

      supported := make(map[string]uint8)
      supported["text/css"] = 1
      supported["text/html"] = 1
      supported["text/javascript"] = 1

		var wg sync.WaitGroup
		defer wg.Wait()

      _, ok := supported[mediatype]
      if ok {
         for file := range in {
            m := minify.New()
            m.AddFunc("text/css", css.Minify)
            m.AddFunc("text/html", html.Minify)
            m.AddFunc("text/javascript", js.Minify)
            wg.Add(1)
            go func(file slurp.File) {
               defer wg.Done()

               minified := new(bytes.Buffer)
               if err := m.Minify(mediatype, minified, file.Reader); err != nil {
                  log.Fatal("js.Minify:", err)
               }

               file.Reader = minified
               file.FileInfo.SetSize(int64(minified.Len()))
               out <- file
            }(file)
         }
		}
	}
}
