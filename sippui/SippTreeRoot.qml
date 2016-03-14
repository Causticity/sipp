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
    	MouseArea {
        	anchors.fill: parent
        	onClicked: {
        		imageViewer.open("image://src/")
        		//imageViewer.title = name
        	}
        }
	}

    SrcImageViewer { id: imageViewer }

}