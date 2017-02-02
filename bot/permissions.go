package bot

import "github.com/sorcix/irc"

type Permissions interface {
	Can(perm string) bool
}

type PermissionsFunc func(perm string) bool

func (pf PermissionsFunc) Can(perm string) bool {
	return pf(perm)
}

type Auther interface {
	Auth(mask *irc.Prefix) (Permissions, error)
}

type AuthFunc func(mask *irc.Prefix) (Permissions, error)

func (af AuthFunc) Auth(mask *irc.Prefix) (Permissions, error) {
	return af(mask)
}
