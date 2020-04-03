// Copyright Raul Vera 2020

// Tests for package sgrad.

package sipptesting

import (
    . "github.com/Causticity/sipp/simage"
    "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

var Sgray *SippGray
var Sgray16 *SippGray16
var SgrayZero *SippGray

func init() {
    Sgray = new(SippGray)
    Sgray.Gray = &sipptestcore.Gray
    Sgray16 = new(SippGray16)
    Sgray16.Gray16 = &sipptestcore.Gray16
    SgrayZero = new(SippGray)
    SgrayZero.Gray = &sipptestcore.GrayZero
}
