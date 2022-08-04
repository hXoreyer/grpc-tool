const { app, BrowserWindow, Menu, ipcMain } = require('electron')
const path = require('path')
const exec = require('child_process')

const NODE_ENV = process.env.NODE_ENV

const createMainWindow = () => {
    Menu.setApplicationMenu(null)
    const mainWin = new BrowserWindow({
        minWidth: 1260,
        minHeight: 700,
        width: 1260,
        height: 700,
        frame: false,
        transparent: true,
        backgroundColor: '#00000000',
        webPreferences: {
            preload: path.join(__dirname, 'preload.ts')
        }
    })

    ipcFunc(mainWin, app)

    if (NODE_ENV === 'development') {
        mainWin.loadURL('http://localhost:5173')
    } else {
        mainWin.loadFile("dist/index.html")
    }

    if (NODE_ENV === "development") {
        mainWin.webContents.openDevTools()
    }
}

app.whenReady().then(() => {
    require('child_process').spawn(path.join(process.cwd(), '/resources/gRpcTool.exe'))
    createMainWindow()

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) createWindow()
    })
})


app.on('window-all-closed', () => {
    if (process.platform != "darwin") {
        app.quit()
    }
})

function ipcFunc(mainWin, app) {
    ipcMain.on('move-title', (event, pos) => {
        mainWin && mainWin.setPosition(pos.posX, pos.posY)
    })

    ipcMain.on('quit-app', () => {
        app && app.quit()
    })

    ipcMain.on('min-app', () => {
        mainWin && mainWin.minimize()
    })

}