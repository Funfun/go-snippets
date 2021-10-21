package good

/*
	GOOD VARIANT #1 WITH FILTERING USING INTERFACES
*/

type President struct {
	status   Status
	currency Currency
	person   string
}
type Status string
type Currency int

//GoodFiltering1 - main filtering struct
type GoodFiltering1 struct {
	//some dependencies
}

//Parametr interface - lets us use parameterized filtering, without changing any filtering methods
type Parametr interface {
	Match(*President) bool
}

//Filter - it's just one method, which gets any kind of filtering params without modifying any logic
func (*GoodFiltering1) Filter(presidents []President, param Parametr) []President {
	var res []President
	for _, pres := range presidents {
		if param.Match(&pres) {
			res = append(res, pres)
		}

	}
	return res
}

//StatusParam - specification for filtering by Status
type StatusParam Status

func (sp StatusParam) Match(p *President) bool {
	return Status(sp) == p.status
}

//CurrencyParam - specification for filtering by Currency
type CurrencyParam Currency

func (cp CurrencyParam) Match(p *President) bool {
	return Currency(cp) == p.currency
}

//PresNameParam - specification for filtering by president's name
type PresNameParam string

func (pn PresNameParam) Match(p *President) bool {
	return string(pn) == p.person
}
