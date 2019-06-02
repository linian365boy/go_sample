package handler

import "go_sample/logkit"

func TestLog(name string) {
	logkit.Infof("I am recevie a name string %s", name)
}