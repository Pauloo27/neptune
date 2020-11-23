package player

const (
	HOOK_PLAYER_INITIALIZED = iota
	HOOK_FILE_LOADED
)

type HookCallback func()

var hooks = make(map[int][]*HookCallback)

func RegisterHook(hookType int, cb HookCallback) {
	if currentHooks, ok := hooks[hookType]; ok {
		hooks[hookType] = append(currentHooks, &cb)
	} else {
		hooks[hookType] = []*HookCallback{&cb}
	}
}

func callHooks(hookType int) {
	if hooks, ok := hooks[hookType]; ok {
		for _, hook := range hooks {
			(*hook)()
		}
	}
}
