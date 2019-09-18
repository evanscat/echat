package echat

type MessageHandleFunc func(dest string, bt []byte) error

type Grouper interface {
	Join(id string, group string)
	Leave(id string, group string)
	Groups(client string) []string
	Clients(group string) []string
	ForEach(group string, bt []byte, handlerFunc MessageHandleFunc)
}

type SimpleGrouper struct {
	Client2Groups map[string]map[string]struct{}
	Group2Clients map[string]map[string]struct{}
}

func NewSimpleGrouper() *SimpleGrouper {
	r := &SimpleGrouper{Client2Groups: make(map[string]map[string]struct{}), Group2Clients: make(map[string]map[string]struct{})}
	return r
}

func (r *SimpleGrouper) ForEach(group string, bt []byte, handlerFunc MessageHandleFunc) {
	if dest, ok := r.Group2Clients[group]; !ok {
		return
	} else {
		for key := range dest {
			_ = handlerFunc(key, bt)
		}
	}
}

func (r *SimpleGrouper) Join(id string, group string) {
	if _, ok := r.Client2Groups[id]; !ok {
		r.Client2Groups[id] = make(map[string]struct{})
	}
	r.Client2Groups[id][group] = struct{}{}
	if _, ok := r.Group2Clients[group]; !ok {
		r.Group2Clients[group] = make(map[string]struct{})
	}
	r.Group2Clients[group][id] = struct{}{}
}

func (r *SimpleGrouper) Leave(id string, group string) {
	delete(r.Client2Groups[id], group)
	if _, ok := r.Client2Groups[id]; !ok {
		delete(r.Client2Groups, id)
	}
	delete(r.Group2Clients[group], id)
	if _, ok := r.Group2Clients[group]; !ok {
		delete(r.Group2Clients, group)
	}
}

func (r *SimpleGrouper) Clients(group string) []string {
	return nil
}

func (r *SimpleGrouper) Groups(client string) []string {
	return nil
}
