// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

ApplicationWindow {
	id: app

    title: "SIPP"

    // Hack to make the top window invisible. But it shouldn't be visible
    // anyway, according to the QML docs.
    width: 1
    height: 1
    
    x: Screen.width/2 - 100
    y: Screen.height/4
    
    signal gotFile(url name) 

	FileDialog {
        id: srcFileDialog
        nameFilters: [ "Image files (*.png)" ]
        folder: "../testdata"
        title: "Select a Greyscale image (8- or 16-bit)"
        onAccepted: app.gotFile(srcFileDialog.fileUrl)
    }
    
    // These indirections are necessary because I can't seem to get access to 
    // the FileDialog object from Go.
    	
    function getFile() {
    	srcFileDialog.open()
    }

}