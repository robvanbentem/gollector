package collector

type CollectorCreator func() Collector

type Collector interface {
	Handle([]byte) string
}
