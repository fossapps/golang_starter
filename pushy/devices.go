package pushy
type Device struct {
	Token string
}
func (device *Device) Notify(data interface {}) {

}
func (device *Device) IsValid() {
	// don't use cache
}
func (device *Device) LastSeen() {
	// use cache 5 min
}
