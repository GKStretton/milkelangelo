package socialmedia

import (
	"os"

	// This api is great!!!
	// https://github.com/Davincible/goinsta/wiki/01.-Getting-Started
	"github.com/Davincible/goinsta/v3"
)

const path = ".goinsta"

func main() {
	insta, err := goinsta.Import(path)
	defer insta.Export(path)
	bail(err)
	f, err := os.Open("/mnt/md0/light-stores/session_content/latest_production/dslr/post/0001.jpg")
	bail(err)
	item, err := insta.Upload(&goinsta.UploadOptions{
		File:    f,
		Caption: "api test",
	})
	_ = item
	bail(err)
}

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func login() {
	insta := goinsta.New("astudyoflight_", "[password]")

	// Only call Login the first time you login. Next time import your config
	if err := insta.Login(); err != nil {
		panic(err)
	}

	// Export your configuration
	// after exporting you can use Import function instead of New function.
	// insta, err := goinsta.Import("~/.goinsta")
	// it's useful when you want use goinsta repeatedly.
	// Export is deffered because every run insta should be exported at the end of the run
	//   as the header cookies change constantly.
	defer insta.Export(".goinsta")
}
