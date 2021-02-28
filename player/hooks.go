package player

const (
	HOOK_PLAYER_INITIALIZED = iota
	HOOK_RESULT_FETCH_STARTED
	HOOK_FILE_LOAD_STARTED
	HOOK_FILE_LOADED
	HOOK_FILE_ENDED
	HOOK_FILE_APPENDED
	HOOK_PLAYBACK_PAUSED
	HOOK_PLAYBACK_RESUMED
	HOOK_VOLUME_CHANGED
	HOOK_POSITION_CHANGED
	HOOK_RESULT_DOWNLOAD_STARTED
	HOOK_QUEUE_UPDATE_FINISHED
	HOOK_PLAYER_EXIT
	HOOK_LOOP_STATUS_CHANGED
)

type HookCallback func(params ...interface{})

var hooks = make(map[int][]*HookCallback)

func RegisterHook(hookType int, cb HookCallback) {
	if currentHooks, ok := hooks[hookType]; ok {
		hooks[hookType] = append(currentHooks, &cb)
	} else {
		hooks[hookType] = []*HookCallback{&cb}
	}
}

func RegisterHooks(hookTypes []int, cb HookCallback) {
	for _, hookType := range hookTypes {
		RegisterHook(hookType, cb)
	}
}

func callHooks(hookType int, params ...interface{}) {
	if hooks, ok := hooks[hookType]; ok {
		for _, hook := range hooks {
			(*hook)(params...)
		}
	}
}
