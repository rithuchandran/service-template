package hotel

type RegionServiceInt interface {
	Update() error
	Search(destination string) (Region, error)
}

type regionService struct {
	repository regionRepositoryInt
	client     clientInt
}

func NewRegionService(repo regionRepositoryInt, client clientInt) *regionService {
	return &regionService{
		repository: repo,
		client:     client,
	}
}

func (s *regionService) Search(destination string) (Region, error) {
	return s.repository.get(destination)
}

func (s *regionService) Update() error {
	reg := s.client.getRegions()
	return s.repository.update(reg)
}
