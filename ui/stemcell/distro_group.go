package stemcell

import (
	"sort"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type DistroGroup struct {
	Distro Distro
	ByName UniqueNameStemcells

	ss []bhstemsrepo.Stemcell // temp state
}

type DistroGroups []DistroGroup

type DistroGroupSorting []DistroGroup

func NewDistroGroups(ss []bhstemsrepo.Stemcell, filter StemcellFilter) DistroGroups {
	var groups []DistroGroup

	for _, d := range allDistros {
		groups = append(groups, DistroGroup{Distro: d})
	}

	for _, s := range ss {
		for i, g := range groups {
			if g.Distro.Matches(s) {
				groups[i].ss = append(groups[i].ss, s)
				break
			}
		}
	}

	var supportedGroups []DistroGroup

	for _, g := range groups {
		uniqueStems := NewUniqueNameStemcells(g.ss, filter)

		if uniqueStems.HasAnyStemcells() {
			g.ByName = uniqueStems
			supportedGroups = append(supportedGroups, g)
		}
	}

	sort.Sort(DistroGroupSorting(supportedGroups))

	return supportedGroups
}

func (g DistroGroup) HasAnyStemcells() bool {
	return g.ByName.HasAnyStemcells()
}

func (g DistroGroups) AllURL() string { return "/stemcells" }

func (s DistroGroupSorting) Len() int           { return len(s) }
func (s DistroGroupSorting) Less(i, j int) bool { return s[i].Distro.Sort < s[j].Distro.Sort }
func (s DistroGroupSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
