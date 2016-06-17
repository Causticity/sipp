// Copyright Raul Vera 2015-2016

import QtQuick 2.2
import QtQuick.Controls 1.1
import QtQuick.Controls.Styles 1.1
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

Window {    
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
        
	Text { 
		id: myText
		anchors.bottom: parent.bottom
	}
    
    onActiveChanged: {
    	if (active) {
    		myText.text = "I have active focus and " + (activeFocusItem==null ? 
    			"there is not an activeFocusItem" : "there IS an activeFocusItem")
    		//if (activeFocusItem != null) {
    		//	newSippTree.parent.menuBar.closeMenuItem.text = "Close Tree"
    		//}
    	} else {
    		myText.text = "I do not have active focus"
		}
    	newSippTree.focusChanged()
    }

	signal focusChanged()

}