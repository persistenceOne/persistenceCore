package mint

type Keeper interface {
	Mint()
}

type BaseKeeper struct {
}

func NewKeeper() Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)

func (baseKeeper BaseKeeper) Mint() {

}
