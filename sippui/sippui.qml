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
    
    menuBar: MenuBar {
        Menu {
            title: "File"
            MenuItem { 
            	text: "New Tree" 
            	shortcut: StandardKey.New
            	onTriggered: newSippTree.start()
            }
            MenuItem {
            	text: "Open Tree"
            	shortcut: StandardKey.Open
            	enabled: false
            }
            MenuItem {
            	text: "Save Tree"
            	enabled: false
            }
            MenuItem {
            	text: "Save Image"
            	enabled: false
            }
            MenuItem {
            	text: "Close Tree"
                shortcut: StandardKey.Close
            	onTriggered: {
            		app.close()
            	}
            	enabled: false
            }
			enabled: true
        }
    }

    SippTree { id: newSippTree }
}