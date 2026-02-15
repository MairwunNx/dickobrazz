package geo

import "fmt"

const (
	maxRegionSize = 60
	excludedSize  = 20 // намеренно исключён
)

var regions = initRegions()

func initRegions() map[int]string {
	m := make(map[int]string, maxRegionSize)
	for i := 0; i <= maxRegionSize; i++ {
		if i == excludedSize {
			continue
		}
		m[i] = fmt.Sprintf("RegionSize%d", i)
	}
	return m
}

func GetRegionBySize(regionId int) string {
	if region, ok := regions[regionId]; ok {
		return region
	}
	return "RegionSizeUnknown"
}
