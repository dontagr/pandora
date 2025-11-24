package factory

import (
	"fmt"

	"github.com/dontagr/pandora/internal/agent/service/interfaces"
	crypro "github.com/dontagr/pandora/pkg/crypto"
)

type KindFactory struct {
	Collection map[string]interfaces.Kind
	cm         *crypro.CManager
}

func NewKindFactory() *KindFactory {
	return &KindFactory{
		Collection: make(map[string]interfaces.Kind),
	}
}

func (f *KindFactory) GetKind(name string) (interfaces.Kind, error) {
	if val, ok := f.Collection[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("kind with name %s not found", name)
}

func (f *KindFactory) SetKind(p interfaces.Kind) {
	f.Collection[p.GetName()] = p
}
