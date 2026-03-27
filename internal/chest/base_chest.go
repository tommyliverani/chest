package chest

type baseChest struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

func (b *baseChest) GetId() string          { return b.Id }
func (b *baseChest) GetName() string        { return b.Name }
func (b *baseChest) GetKind() string        { return b.Kind }
func (b *baseChest) GetDescription() string { return b.Description }
