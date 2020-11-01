package main

type ServerState struct {
	IsPingable         bool
	UsersCount         int
	ServerPid          int
	LastMcCrashCount   int
	McCrashCount       int
	LastJavaCrashCount int
	JavaCrashCount     int
}

func validateState(state *ServerState) {
	// someone removed logs?
	if state.LastJavaCrashCount < state.JavaCrashCount {
		state.JavaCrashCount = state.LastJavaCrashCount
	}
	if state.LastMcCrashCount < state.McCrashCount {
		state.McCrashCount = state.LastMcCrashCount
	}
}
