import QtQuick 2.2
import QtQuick.Window 2.1

Window {
    minimumWidth: srcImage.width
    minimumHeight: srcImage.height
    x: Screen.width/2 - width/2
    y: Screen.height/2 - height/2
    function open(source) {
        srcImage.source = source
        width = srcImage.implicitWidth
        height = srcImage.implicitHeight
        //title = source
        title = "dunno yet"
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
    onActiveChanged: {
    	if (active) {
    		title = "I'm active!"
    		//app.appMenuBar.closeMenuItem.text = "Close Image"
    		// Send Go a signal here, so Go can change the menu
    	} else {
    		title = "I'm NOT active!"
    		// Send Go a signal here, so Go can change the menu
    	}
    }
    // TODO: This doesn't work for the standard key, but does for a regular key,
    // because the app menu catches the key. This is going to have to happen 
    // from Go.
	Item {
		focus: true
		Keys.onPressed: {
			if (event.key == StandardKey.Close) {
				imageViewer.close()
				event.accepted = true;
			}
		}
	}
	
}