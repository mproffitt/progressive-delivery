package crd

import (
	"context"
	"sync"

	"github.com/weaveworks/weave-gitops/core/clustersmngr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func NewNoCacheFetcher(clientPool clustersmngr.Client) Fetcher {
	fetcher := &noCacheFetcher{
		client: clientPool,
		crds:   map[string][]v1.CustomResourceDefinition{},
	}

	return fetcher
}

type noCacheFetcher struct {
	sync.RWMutex
	client clustersmngr.Client
	crds   map[string][]v1.CustomResourceDefinition
}

func (s *noCacheFetcher) UpdateCRDList() {
	s.Lock()
	defer s.Unlock()

	ctx := context.Background()

	for clusterName, client := range s.client.ClientsPool().Clients() {
		crdList := &v1.CustomResourceDefinitionList{}

		err := client.List(ctx, crdList)
		if err != nil {
			// Handle error, or log at least.
			continue
		}

		s.crds[clusterName] = crdList.Items
	}
}

func (s *noCacheFetcher) IsAvailable(clusterName, name string) bool {
	s.UpdateCRDList()

	s.Lock()
	defer s.Unlock()

	for _, crd := range s.crds[clusterName] {
		if crd.Name == name {
			return true
		}
	}

	return false
}