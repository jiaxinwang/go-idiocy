package platform

// InstanceCall ...
type InstanceCall struct {
	Pkg      string
	FuncName string
}

var GinInstanceCall = []InstanceCall{
	{`gin`, `Default`},
	{`gin`, `New`},
}
