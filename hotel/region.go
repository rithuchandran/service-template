package hotel

type Region struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	NameFull    string `json:"name_full"`
	Descriptor  string `json:"descriptor"`
	Ancestors   []Data `json:"ancestors"`
	Descendants map[string][]string
}
type Regions map[string]Region

type Data struct {
	Id   string
	Type string
}