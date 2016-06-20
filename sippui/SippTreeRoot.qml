// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.1
import QtQuick.Controls.Styles 1.1
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

ApplicationWindow {    
    width: 500
    height: 500
    
    x: Screen.width/15
    y: Screen.height/15
    
    id: newSippTree
        
    function setThumbSource(name) {
		thumb.source = "image://thumb/" + name
	}
    
    Image {
    	id: thumb
    	objectName: "thumb"
        anchors.centerIn: parent
        cache: false
    	MouseArea {
        	anchors.fill: parent
        	onClicked: newSippTree.thumbClicked()
        }
	}

	Button {
		text: "Gradient"
		style: ButtonStyle { }
		anchors.right: parent.right
		anchors.verticalCenter: parent.verticalCenter
		onClicked: newSippTree.gradientClicked()
	}
    
    signal thumbClicked()
    signal gradientClicked()
        
    menuBar: MenuBar {
    	id: nodeMenuBar
        Menu {
	        title: "File"
//            MenuItem { 
 //           	text: "New Tree" 
  //          	shortcut: StandardKey.New
   //         	objectName: "newTree"
    //        	onTriggered: app.getFile()
     //       	enabled: true
      //      }
            MenuItem {
            	text: "Close Tree"
                shortcut: StandardKey.Close
            	objectName: "closeTree"
            	onTriggered: {
            		newSippTree.close()
            	}
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