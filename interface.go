package dix

type iContainerData interface {
	setAccessed() (isFirstAccess bool)
	triggerOnCloseHook()
}
