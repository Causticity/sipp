import QtQuick 2.2
import QtQuick.Dialogs 1.0
import QtQuick.Window 2.1

Window {    
    width: 500
    height: 500
    
    x: Screen.width/15
    y: Screen.height/15
    
    id: newSippTree
    
    FileDialog {
        id: srcFileDialog
        nameFilters: [ "Image files (*.png)" ]
        folder: "../testdata"
        onAccepted: {
	    	// The URL comes back with a "file://" prefix, so we remove that.
	    	newSippTree.gotFile(srcFileDialog.fileUrl.toString().substring(6))
	    }
    }
    
    // These indirections are necessary because I can't seem to get access to 
    // the FileDialog object from Go.
	signal gotFile(string filename) 
	
    function getFile() {
    	srcFileDialog.open()
    }
    
    function setThumbSource(name) {
		thumb.source = "image://thumb/" + name
	}
    
    Image {
    	id: thumb
    	objectName: "thumby"
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