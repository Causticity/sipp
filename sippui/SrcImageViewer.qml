import QtQuick 2.2
import QtQuick.Window 2.1

Window {
    id: imageViewer
    minimumWidth: srcImage.width
    minimumHeight: srcImage.height
    x: Screen.width/2 - width/2
    y: Screen.height/2 - height/2
    function open(source) {
        srcImage.source = source
        width = srcImage.implicitWidth
        height = srcImage.implicitHeight
        title = source
        visible = true
    }
    Image {
        id: srcImage
        anchors.centerIn: parent
    }
}
