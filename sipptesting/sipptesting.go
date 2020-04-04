// Copyright Raul Vera 2020

// Test infrastructure for the sipp package. Note that this depends on the
// sipp/simage package, so it can't be imported by the tests for that package.
// Test infrastructure used by the simage tests and shared with other tests
// can be found in the sipptestcore package.

package sipptesting

import (
    . "github.com/Causticity/sipp/simage"
    "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

var Sgray *SippGray
var Sgray16 *SippGray16
var SgrayZero *SippGray
var SgrayCosxCosyTiny *SippGray

func init() {
    Sgray = new(SippGray)
    Sgray.Gray = &sipptestcore.Gray
    Sgray16 = new(SippGray16)
    Sgray16.Gray16 = &sipptestcore.Gray16
    SgrayZero = new(SippGray)
    SgrayZero.Gray = &sipptestcore.GrayZero
    SgrayCosxCosyTiny = new(SippGray)
    SgrayCosxCosyTiny.Gray = &sipptestcore.CosxCosyTiny
}
