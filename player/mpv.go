package player

import (
	"math"

	"github.com/Pauloo27/my-tune/utils"
	"github.com/YouROK/go-mpv/mpv"
)

const (
	maxVolume = 150.0
)

var MpvInstance *mpv.Mpv
var State *PlayerState
var DataFolder string

func Initialize(dataFolder string) {
	var err error

	DataFolder = dataFolder

	initialVolume := 50.0

	// create a mpv instance
	MpvInstance = mpv.Create()

	// set options
	// disable video
	err = MpvInstance.SetOptionString("video", "no")
	utils.HandleError(err, "Cannot set mpv video option")

	// disable cache
	err = MpvInstance.SetOptionString("cache", "no")
	utils.HandleError(err, "Cannot set mpv cache option")

	// set default volume value
	err = MpvInstance.SetOption("volume", mpv.FORMAT_DOUBLE, initialVolume)
	utils.HandleError(err, "Cannot set mpv volume option")

	// set default volume value
	err = MpvInstance.SetOption("volume-max", mpv.FORMAT_DOUBLE, maxVolume)
	utils.HandleError(err, "Cannot set mpv volume-max option")

	// set quality to worst
	err = MpvInstance.SetOptionString("ytdl-format", "worst")
	utils.HandleError(err, "Cannot set mpv ytdl-format option")

	// start event listener
	startEventHandler()

	// create the state
	State = &PlayerState{
		false,
		nil,
		initialVolume,
		0.0,
	}

	// internal hooks
	RegisterHook(HOOK_FILE_LOAD_STARTED, func(err error, params ...interface{}) {
		Play()
		ClearPlaylist()
	})

	// start the player
	err = MpvInstance.Initialize()
	utils.HandleError(err, "Cannot initialize mpv")

	callHooks(HOOK_PLAYER_INITIALIZED, err)
}

func ClearPlaylist() error {
	return MpvInstance.Command([]string{"playlist-clear"})
}

func RemoveCurrentFromPlaylist() error {
	return MpvInstance.Command([]string{"playlist-remove", "current"})
}

func LoadFile(filePath string) error {
	err := MpvInstance.Command([]string{"loadfile", filePath})
	callHooks(HOOK_FILE_LOAD_STARTED, err, filePath)
	return err
}

func PlayPause() error {
	if State.Paused {
		return Play()
	} else {
		return Pause()
	}
}

func Pause() error {
	err := MpvInstance.SetProperty("pause", mpv.FORMAT_FLAG, true)
	callHooks(HOOK_PLAYBACK_PAUSED, err)
	return err
}

func Play() error {
	err := MpvInstance.SetProperty("pause", mpv.FORMAT_FLAG, false)
	callHooks(HOOK_PLAYBACK_RESUMED, err)
	return err
}

func SetVolume(volume float64) error {
	volume = math.Min(maxVolume, volume)
	err := MpvInstance.SetProperty("volume", mpv.FORMAT_DOUBLE, volume)
	callHooks(HOOK_VOLUME_CHANGED, err, volume)
	return err
}

func GetPosition() (float64, error) {
	position, err := MpvInstance.GetProperty("time-pos", mpv.FORMAT_DOUBLE)
	if err != nil {
		return 0.0, err
	}
	return position.(float64), err
}

func SetPosition(pos float64) error {
	err := MpvInstance.SetProperty("time-pos", mpv.FORMAT_DOUBLE, pos)
	callHooks(HOOK_POSITION_CHANGED, err, pos)
	return err
}
