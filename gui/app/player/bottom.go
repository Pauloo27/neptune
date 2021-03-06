package player

import (
	"time"

	"github.com/Pauloo27/neptune/hook"
	"github.com/Pauloo27/neptune/player"
	"github.com/Pauloo27/neptune/utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var positionLabel *gtk.Label
var progressBar *gtk.Scale
var currentPosition float64

func createProgressBar() *gtk.Scale {
	var err error

	progressBar, err = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0.0, 1.0, 0.01)
	utils.HandleError(err, "Cannot create scale")

	progressBar.SetDrawValue(false)
	progressBar.SetHExpand(true)
	progressBar.Connect("value-changed", func() {
		value := progressBar.GetValue()
		if value == currentPosition {
			return
		}
		player.SetPosition(value * player.State.Duration)
	})

	return progressBar
}

func createVolumeController() *gtk.Box {
	volumeContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	volumeIcon, err := gtk.ImageNewFromIconName("audio-volume-medium", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create image")

	volumeController, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0.0, 100.0, 1.0)
	utils.HandleError(err, "Cannot create box")

	volumeController.SetDrawValue(false)
	volumeController.SetValue(player.State.Volume)

	hook.RegisterHook(hook.HOOK_VOLUME_CHANGED, func(params ...interface{}) {
		volume := params[0].(float64)
		glib.IdleAdd(func() {
			if volume != volumeController.GetValue() {
				volumeController.SetValue(volume)
			}
		})
	})

	volumeController.Connect("value-changed", func() {
		player.SetVolume(volumeController.GetValue())
	})

	volumeContainer.PackStart(volumeIcon, false, false, 0)
	volumeContainer.PackEnd(volumeController, true, true, 0)

	return volumeContainer
}

func createDurationLabel() *gtk.Label {
	durationLabel, err := gtk.LabelNew("--:--")
	utils.HandleError(err, "Cannot create label")

	durationLabel.SetHAlign(gtk.ALIGN_END)
	hook.RegisterHook(hook.HOOK_FILE_LOADED, func(params ...interface{}) {
		duration := params[1].(float64)
		glib.IdleAdd(func() {
			durationLabel.SetText(utils.FormatDuration(duration))
		})
	})

	return durationLabel
}

func updatePosition(position float64) {
	positionLabel.SetText(utils.FormatDuration(position))
	currentPosition = position / player.State.Duration
	progressBar.SetValue(currentPosition)
}

func progressUpdater() {
	for {
		position, err := player.GetPosition()
		if err == nil {
			glib.IdleAdd(func() {
				updatePosition(position)
			})
		}
		time.Sleep(1 * time.Second)
	}
}

func createPositionLabel() *gtk.Label {
	var err error
	positionLabel, err = gtk.LabelNew("--:--")
	utils.HandleError(err, "Cannot create label")

	positionLabel.SetHAlign(gtk.ALIGN_START)

	go progressUpdater()

	return positionLabel
}

func createSongLabel() *gtk.Label {
	songLabel, err := gtk.LabelNew("--")
	utils.HandleError(err, "Cannot create label")

	songLabel.SetHAlign(gtk.ALIGN_CENTER)

	hook.RegisterHook(hook.HOOK_RESULT_FETCH_STARTED, func(params ...interface{}) {
		entry := player.State.Fetching
		glib.IdleAdd(func() {
			songLabel.SetText(utils.Fmt("Fetching %s...", entry.Title))
		})
	})

	hook.RegisterHook(hook.HOOK_FILE_LOADED, func(params ...interface{}) {
		track := player.GetCurrentTrack()
		glib.IdleAdd(func() {
			songLabel.SetText(utils.Fmt("%s - %s", track.Album.Artist.Name, track.Title))
		})
	})

	return songLabel
}

func createTimeStampContainer() *gtk.Box {
	timeStampContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	utils.HandleError(err, "Cannot create box")

	timeStampContainer.PackStart(createPositionLabel(), false, false, 0)
	timeStampContainer.PackEnd(createDurationLabel(), false, false, 0)

	return timeStampContainer
}

func createButtonsContainer() *gtk.Box {
	buttonsContainer, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
	utils.HandleError(err, "Cannot create box")

	pausedIcon, err := gtk.ImageNewFromIconName("media-playback-start", gtk.ICON_SIZE_BUTTON)
	playingIcon, err := gtk.ImageNewFromIconName("media-playback-pause", gtk.ICON_SIZE_BUTTON)

	pauseButton, err := gtk.ButtonNew()
	utils.HandleError(err, "Cannot create button")

	pauseButton.SetImage(playingIcon)

	hook.RegisterHook(hook.HOOK_PLAYBACK_PAUSED, func(params ...interface{}) {
		glib.IdleAdd(func() {
			pauseButton.SetImage(pausedIcon)
		})
	})
	hook.RegisterHook(hook.HOOK_PLAYBACK_RESUMED, func(params ...interface{}) {
		glib.IdleAdd(func() {
			pauseButton.SetImage(playingIcon)
		})
	})

	pauseButton.Connect("clicked", func() {
		player.PlayPause()
	})

	prevButton, err := gtk.ButtonNewFromIconName("media-seek-backward", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	prevButton.Connect("clicked", func() {
		player.PreviousTrack()
	})

	nextButton, err := gtk.ButtonNewFromIconName("media-seek-forward", gtk.ICON_SIZE_BUTTON)
	utils.HandleError(err, "Cannot create button")

	nextButton.Connect("clicked", func() {
		player.NextTrack()
	})

	buttonsContainer.SetHAlign(gtk.ALIGN_CENTER)

	buttonsContainer.PackStart(prevButton, false, true, 0)
	buttonsContainer.PackStart(pauseButton, false, true, 0)
	buttonsContainer.PackStart(nextButton, false, true, 0)

	return buttonsContainer
}

func createPlayerBottom() *gtk.Grid {
	bottomContainer, err := gtk.GridNew()
	utils.HandleError(err, "Cannot create grid")

	// row 0
	bottomContainer.Attach(createSongLabel(), 0, 0, 10, 1)
	// row 1
	bottomContainer.Attach(createProgressBar(), 0, 1, 10, 1)
	// row 2
	bottomContainer.Attach(createTimeStampContainer(), 0, 2, 10, 1)
	// row 3
	bottomContainer.Attach(createVolumeController(), 0, 3, 3, 1)
	bottomContainer.Attach(createButtonsContainer(), 3, 3, 4, 1)

	return bottomContainer
}
