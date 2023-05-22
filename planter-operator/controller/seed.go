package controller

import (
	resources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
)

func (self *Controller) processSeed(seed *resources.Seed) (bool, error) {
	if err := self.plantSeed(seed); err == nil {
		if err := self.reconcileSeed(seed); err == nil {
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}
}

func (self *Controller) plantSeed(seed *resources.Seed) error {
	if seed.Spec.Planted {
		if content, err := self.Client.GetContent(seed.Spec.SeedURL); err == nil {
			if content, err = self.execPlugins(content); err == nil {
				_, err := self.Client.PlantSeed(seed, content)
				return err
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return nil
	}
}

func (self *Controller) reconcileSeed(seed *resources.Seed) error {
	if seed.Spec.Planted && (seed.Status.PlantedPath != "") {
		if content, err := self.Client.GetContent(seed.Status.PlantedPath); err == nil {
			if resources, err := NewResources(content, self); err == nil {
				for {
					resource := resources.next()
					if resource == nil {
						break
					}
					if err := resources.Reconcile(resource); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
