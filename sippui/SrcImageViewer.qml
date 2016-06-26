// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.4
import QtQuick.Window 2.1

ApplicationWindow {

	id:src

    minimumWidth: srcImage.width
    minimumHeight: srcImage.height
    x: Screen.width/2 - width/2
    y: Screen.height/2 - height/2
    
    function open(source) {
        srcImage.source = "image://src/"+source
        width = srcImage.implicitWidth
        height = srcImage.implicitHeight
        title = source
        visible = true
    }
    Image {
        id: srcImage
        anchors.centerIn: parent
        cache: false
    }
    Component.onCompleted: {
    	requestActivate()
    }

    menuBar: MenuBar {
    	id: nodeMenuBar
        Menu {
	        title: "File"
            MenuItem { 
           	    text: "New Tree" 
           	    shortcut: StandardKey.New
           	    objectName: "newTree"
           	    onTriggered: nodeFileDialog.open()
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
            	text: "Close Image"
                shortcut: StandardKey.Close
            	objectName: "closeImage"
            	onTriggered: src.close()
            	enabled: true
            }
			enabled: true
        }
    }
}