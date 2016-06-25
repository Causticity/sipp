// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

FileDialog {
    x: Screen.width/15
    y: Screen.height/15
	nameFilters: [ "Image files (*.png)" ]
	folder: "../../testdata"
	title: "Select a Greyscale image (8- or 16-bit)"
}
