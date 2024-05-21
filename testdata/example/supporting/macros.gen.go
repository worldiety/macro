// Code generated by github.com/worldiety/macro. DO NOT EDIT.

package supporting

// Permissions provides a complete slice of all annotated permissions of all bounded contexts.
// Each use case, which requires some sort of auditing, has its individual permission.
func Permissions() []Permission {
	return []Permission{
		{"de.worldiety.aufstehen2", "Zeiten loggen", "Cooles Zeitbuchen ist angesagt."},
		{"de.worldiety.aufstehen", "Aufstehen", ""},
	}
}

// Permission represents a permission to call a distinct use case. It provides method accessors,
// so that other permission consumers can accept their own interfaces.
type Permission struct {
	id   string
	name string
	desc string
}

func (p Permission) ID() string {
	return p.id
}

func (p Permission) Name() string {
	return p.name
}

func (p Permission) Desc() string {
	return p.desc
}
