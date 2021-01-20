// Copyright Raul Vera 2020

// Test infrastructure for the sipp package. Note that this depends on the
// sipp/simage package, so it can't be imported by the tests for that package.
// Test infrastructure used by the simage tests and shared with other tests
// can be found in the sipptestcore package.

package sipptesting

import (
	"path/filepath"
	"strings"

	. "github.com/Causticity/sipp/simage"
	"github.com/Causticity/sipp/sipptesting/sipptestcore"
)

var Sgray *SippGray
var Sgray16 *SippGray16
var SgrayCosxCosyTiny *SippGray

func init() {
	Sgray = new(SippGray)
	Sgray.Gray = &sipptestcore.Gray
	Sgray16 = new(SippGray16)
	Sgray16.Gray16 = &sipptestcore.Gray16
	SgrayCosxCosyTiny = new(SippGray)
	SgrayCosxCosyTiny.Gray = &sipptestcore.CosxCosyTiny
}

func SaveFailedSimage(expFileName string, simg SippImage) (name string) {
	// Compose a new name from the expected file name, with "FAILED" added to
	// the base name.
	ext := filepath.Ext(expFileName)
	name = strings.TrimSuffix(expFileName, ext) + "_FAILED" + ext
	// Write the failed image out with that name, overwriting any existing one.
	err := simg.Write(&name)
	if err != nil {
		panic("Error writing failed image file");
	}
	return
}