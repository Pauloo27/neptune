package explorer

import (
	"github.com/Pauloo27/neptune/db"
	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createQueueEntry(track *db.Track, queueIndex int) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	trackLabel, err := gtk.LabelNew(utils.Fmt("%s - %s", utils.EnforceSize(track.Album.Artist.Name, 15), utils.EnforceSize(track.Title, 50)))
	utils.HandleError(err, "Cannot create label")

	moveUpButton, err := gtk.ButtonNewFromIconName("go-up", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	moveUpButton.Connect("clicked", func() {
		player.MoveUpInQueue(queueIndex)
	})

	moveDownButton, err := gtk.ButtonNewFromIconName("go-down", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	moveDownButton.Connect("clicked", func() {
		player.MoveDownInQueue(queueIndex)
	})

	removeButton, err := gtk.ButtonNewFromIconName("delete", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	removeButton.Connect("clicked", func() {
		player.RemoveFromQueue(queueIndex)
	})

	container.PackStart(trackLabel, false, false, 1)
	container.PackEnd(removeButton, false, false, 1)
	container.PackEnd(moveDownButton, false, false, 1)
	container.PackEnd(moveUpButton, false, false, 1)

	return container
}

func createQueuePage() *gtk.ScrolledWindow {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	headerContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	utils.HandleError(err, "Cannot create box")

	headerContainer.SetHAlign(gtk.ALIGN_CENTER)

	shuffleButton, err := gtk.ButtonNewFromIconName("shuffle", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	shuffleButton.Connect("clicked", func() {
		player.Shuffle()
	})

	clearQueueButton, err := gtk.ButtonNewFromIconName("delete", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	clearQueueButton.Connect("clicked", func() {
		player.ClearQueue()
	})

	loopButton, err := gtk.ButtonNew()
	utils.HandleError(err, "Cannot create button")

	loopNoneIcon, err := gtk.ImageNewFromIconName("format-justify-fill", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create image")

	loopTrackIcon, err := gtk.ImageNewFromIconName("media-playlist-repeat-song", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create image")

	loopQueueIcon, err := gtk.ImageNewFromIconName("media-playlist-repeat", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create image")

	loopButton.SetImage(loopNoneIcon)

	loopButton.Connect("clicked", func() {
		player.NextLoopStatus()
	})

	hook.RegisterHook(hook.HOOK_LOOP_STATUS_CHANGED, func(params ...interface{}) {
		loopStatus := player.GetLoopStatus()
		glib.IdleAdd(func() {
			var icon *gtk.Image
			switch loopStatus {
			case player.LOOP_NONE:
				icon = loopNoneIcon
			case player.LOOP_TRACK:
				icon = loopTrackIcon
			case player.LOOP_QUEUE:
				icon = loopQueueIcon
			}
			loopButton.SetImage(icon)
		})
	})

	headerContainer.PackStart(shuffleButton, false, false, 1)
	headerContainer.PackStart(clearQueueButton, false, false, 1)
	headerContainer.PackStart(loopButton, false, false, 1)

	queueContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	utils.HandleError(err, "Cannot create box")

	container.PackStart(headerContainer, false, false, 1)
	container.PackStart(queueContainer, false, false, 1)

	hook.RegisterHook(
		hook.HOOK_QUEUE_UPDATE_FINISHED,
		func(params ...interface{}) {
			glib.IdleAdd(func() {
				queueContainer.GetChildren().Foreach(func(item interface{}) {
					item.(*gtk.Widget).Destroy()
				})
				for i := 0; i < len(player.State.Queue); i++ {
					track := player.GetTrackAt(i)

					queueContainer.PackStart(createQueueEntry(track, i), false, false, 1)
				}
				queueContainer.ShowAll()
			})
		},
	)

	scrolledContainer, err := gtk.ScrolledWindowNew(nil, nil)
	utils.HandleError(err, "Cannot create scrolled window")

	scrolledContainer.Add(container)

	return scrolledContainer
}
