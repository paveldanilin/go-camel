package camel

type Route struct {
	id       string
	from     string
	producer Producer
}

func NewRoute(id string, from string, producer Producer) *Route {
	return &Route{
		id:       id,
		from:     from,
		producer: producer,
	}
}

func (r *Route) Id() string {

	return r.id
}

func (r *Route) Producer() Producer {

	return r.producer
}
