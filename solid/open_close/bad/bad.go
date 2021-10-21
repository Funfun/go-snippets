package bad

type President struct {
	status   Status
	currency Currency
}
type Status string
type Currency int

type BadFiltering struct {
	//some dependencies
}

func (*BadFiltering) ListByStatus(presidents []President, status Status) []President {
	var res []President
	for _, p := range presidents {
		if p.status == status {
			res = append(res, p)
		}
	}
	return res
}

func (*BadFiltering) ListByCurrency(presidents []President, currency Currency) []President {
	var res []President
	for _, p := range presidents {
		if p.currency == currency {
			res = append(res, p)
		}
	}
	return res
}
