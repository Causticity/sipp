import QtQuick 2.2
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

Window {
    id: newSippTree
    
    width: 500
    height: 500
    
    x: Screen.width/15
    y: Screen.height/15

    FileDialog {
        id: srcFileDialog
        nameFilters: [ "Image files (*.png)" ]
        onAccepted: {
        	// The URL comes back with a "file://" prefix, so we remove that.
        	var name = fileUrl.toString().substring(6)
        	thumb.source = "image://thumb/" + name
        	newSippTree.show()
        	newSippTree.title = name

        }
    }

    function start() {
    	srcFileDialog.open()
    }
    
    Image {
    	id: thumb
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