package regia

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) byName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

func (ps Params) Get(key string) Value {
	v := ps.byName(key)
	return Value(v)
}
