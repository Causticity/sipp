import QtQuick 2.2
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

	Text { id: myText }
    
    signal focusChanged()
    signal thumbClicked()
    
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

}