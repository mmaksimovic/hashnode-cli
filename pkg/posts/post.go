package posts

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// openPost opens a post in a new tview box
func openPost(app *tview.Application, postcuid string, list *tview.List) {
	var singlePost Post
	b, err := makeRequest(fmt.Sprintf("%s/%s", postAPI, postcuid))
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &singlePost)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})

	title := fmt.Sprintf("Title: %s", singlePost.Post.Title)
	author := fmt.Sprintf("Author: %s", singlePost.Post.Author.Name)
	upvotes := fmt.Sprintf("Upvotes: %d", singlePost.Post.Upvotes)
	ptype := fmt.Sprintf("Type: %s", singlePost.Post.Type)
	link := fmt.Sprintf("Link: https://hashnode.com/post/%s", singlePost.Post.Cuid)
	writeToTextView(textView, title,
		author,
		upvotes,
		ptype,
		link,
		"\n",
		singlePost.Post.ContentMarkdown,
	)
	textView.Box = textView.Box.SetBorder(true).SetBorderPadding(1, 1, 2, 1)
	textView.SetBorder(true)

	textView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
				app.Stop()
				panic(err)
			}
		}
	})

	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		app.Stop()
		panic(err)
	}
}

func writeToTextView(t *tview.TextView, contents ...string) {
	for _, content := range contents {
		t.Write([]byte(content))
		t.Write([]byte("\n"))
	}
}
