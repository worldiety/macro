// #[go.permission.generateTable]
package supporting

type User struct {
}

func (User) Audit(string) error {
	panic("implement me")
}
