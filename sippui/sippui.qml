// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

// TODO: Try making each window it's own application window, with slightly
// different menubars for the two types. Then each action can be tied to
// a receiver function and we can leave tracking the current window to the
// window system, where it belongs.

ApplicationWindow {
	id: app

    title: "SIPP"

    // Hack to make the top window invisible. But it shouldn't be visible
    // anyway, according to the QML docs.
    width: 1
    height: 1
    
    x: Screen.width/2 - 100
    y: Screen.height/4
    
    menuBar: MenuBar {
    	id: appMenuBar
        Menu {
            title: "File"
            MenuItem { 
            	text: "New Tree" 
            	shortcut: StandardKey.New
            	objectName: "newTree"
            	onTriggered: app.getFile()
            	enabled: true
            }
            //MenuItem {
            //	text: "Open Tree"
            //	shortcut: StandardKey.Open
            //	objectName: "openTree"
            //	enabled: false
            //}
            //MenuItem {
            //	text: "Save Tree"
            //	objectName: "saveTree"
            //	enabled: false
            //}
            //MenuItem {
            //	text: "Save Image"
            //	objectName: "saveImage"
            //	enabled: false
            //}
            MenuItem {
            	text: "Close Tree"
                shortcut: StandardKey.Close
            	objectName: "closeTree"
            	onTriggered: {
            		// Really needs an "Are you sure?" dialog
            		app.closeCurrentTree()
            	}
            	enabled: true
            }
            MenuItem {
            	text: "Close Image"
                shortcut: StandardKey.Close
            	objectName: "closeImage"
            	onTriggered: {
            		// Really needs an "Are you sure?" dialog
            		app.closeCurrentImage()
            	}
            	enabled: false
            }
			enabled: true
        }
    }

    signal closeCurrentTree()
    
	FileDialog {
        id: srcFileDialog
        nameFilters: [ "Image files (*.png)" ]
        folder: "../testdata"
        title: "Select a Greyscale image (8- or 16-bit)"
        onAccepted: {
	    	// The URL comes back with a "file://" prefix, so we remove that.
	    	//app.gotFile(srcFileDialog.fileUrl.toString().substring(6))
	    	app.gotFile(srcFileDialog.fileUrl)
	    }
    }
    
    // These indirections are necessary because I can't seem to get access to 
    // the FileDialog object from Go.
    
	signal gotFile(url name) 
	
    function getFile() {
    	srcFileDialog.open()
    }
    
}