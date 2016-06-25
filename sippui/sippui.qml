// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

import "uiComponents"

ApplicationWindow {
	id: app

    title: "SIPP"
    
    // We fiddle the location of the top window to keep it off screen most of
    // the time, or at least unobtrusively in a corner, and bringing it on
    // screen when we open the file dialog. We also make it tiny because it has
    // no content. I wish I could figure out a way to do this without a useless
    // tiny window.
    width: 1
    height: 1
    x: 0
    y: 0

    SippFileDialog { 
    	id:srcFileDialog
    	onAccepted: {
    		app.x = -10
    		app.y = -10
    		app.gotFile(srcFileDialog.fileUrl)
    	}
    }	
    
    signal gotFile(url name) 

    function getFile() {
    	app.x = Screen.width/2 - 100
    	app.y = Screen.height/4
    	srcFileDialog.open()
    }

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
            MenuItem {
            	text: "Close Tree"
                shortcut: StandardKey.Close
            	objectName: "closeTree"
            	enabled: false
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
            //MenuItem {
            //	text: "Close Image"
            //   shortcut: StandardKey.Close
            //	objectName: "closeImage"
            //	onTriggered: {
            //		// Really needs an "Are you sure?" dialog
            //		app.closeCurrentImage()
            //	}
            //	enabled: false
            //}
			enabled: true
        }
    }

}