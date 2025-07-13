package dix

type iContainerData interface {
	setAccessed() (isFirstAccess bool)
	lock()
	unlock()
	triggerOnCloseHook()
}
