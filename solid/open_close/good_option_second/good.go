package good

/*
	GOOD VARIANT #2 WITH FILTERING USING PARAM FUNCS
*/
type President struct {
}
type Status string
type Currency int

type ParamFn func(*President) bool

//GoodFiltering1 - main filtering struct
type GoodFiltering2 struct {
	//some dependencies
}

//Filter - it's just one method, which gets any kind of filtering params without modifying any logic
func (*GoodFiltering2) Filter(presidents []President, param ParamFn) []President {
	var res []President
	for _, pres := range presidents {
		if param(&pres) {
			res = append(res, pres)
		}
	}
	return res
}

/*
	As an example:
	gf2 := GoodFiltering2{}
	paramCurrencyFn := ParamFn(func(p *President) bool {
		return p.currency == RUB
	})
	result := gf2.Filter(seedData, paramCurrencyFn)
*/
