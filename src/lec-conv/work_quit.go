package main

type QuitWork struct {
}

func (w QuitWork) Run() bool {
	return true
}

func (w QuitWork) IsQuit() bool {
	return true
}
