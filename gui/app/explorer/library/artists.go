package library

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	artistsPage = &LibraryPage{"Artists", showArtists}
)

func displayArtist(artist *db.Artist) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(artist.Name)
	utils.HandleError(err, "Cannot create button")

	btn.Connect("clicked", func() {
		displayPage(createArtistPage(artist))
	})

	return btn
}

func showArtists() *gtk.Grid {
	container, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	container.SetRowSpacing(5)
	container.SetColumnHomogeneous(true)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	go func() {
		artists, err := db.ListArtists()
		utils.HandleError(err, "Cannot list artist")

		glib.IdleAdd(func() {
			for i, artist := range artists {
				container.Attach(displayArtist(artist), 0, i, 1, 1)
			}
			container.ShowAll()
		})
	}()

	return container
}
