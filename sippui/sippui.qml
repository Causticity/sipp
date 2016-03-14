import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

ApplicationWindow {
	id: app

    title: "SIPP"

    // Hack to make the top window invisible, hopefully temporary until I can
    // figure out how to properly do SDI with QtQuick.
    width: 1
    height: 1
    
    x: Screen.width/2 - 100
    y: Screen.height/4

    menuBar: MenuBar {
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
            		// and must apply to the "currrent" tree
            		// and must call into Go to manage Go structs
            		app.close()
            	}
            	enabled: true
            }
			enabled: true
        }
    }
    
	FileDialog {
        id: srcFileDialog
        nameFilters: [ "Image files (*.png)" ]
        folder: "../testdata"
        title: "Select a Greyscale image (8- or 16-bit)"
        onAccepted: {
	    	// The URL comes back with a "file://" prefix, so we remove that.
	    	app.gotFile(srcFileDialog.fileUrl.toString().substring(6))
	    }
    }
    
    // These indirections are necessary because I can't seem to get access to 
    // the FileDialog object from Go.
    
	signal gotFile(string filename) 
	
    function getFile() {
    	srcFileDialog.open()
    }
    

}