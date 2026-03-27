package jewel

// ok
type baseJewel struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

func (j *baseJewel) GetName() string        { return j.Name }
func (j *baseJewel) GetKind() string        { return j.Kind }
func (j *baseJewel) GetDescription() string { return j.Description }
