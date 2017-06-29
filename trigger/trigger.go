package trigger

type Trigger interface {
	Start()
	Mask()
	Unmask()
	Stop()
	Wait()
}
