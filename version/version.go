package version

import (
	"fmt"
	"strconv"
)

type Version uint16

var (
	Application = "Indexer"
	Description = "Offchain index Pactus blockchain to centralized database"

	CommitID  string
	BuildTime string

	Major Version = 0
	Minor Version = 1
	Patch Version = 0
)

func (v Version) String() string {
	return strconv.Itoa(int(v))
}

func Semantic() string {
	return fmt.Sprintf("v%v.%v.%v", Major, Minor, Patch)
}

func Full() string {
	return fmt.Sprintf("%s v%v.%v.%v", Application, Major, Minor, Patch)
}
