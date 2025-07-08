package dix

type iContainerData interface {
	setAccessed()
	lock()
	unlock()
	triggerOnCloseHook()
}
