const { app, BrowserWindow, Menu, ipcMain, Notification } = require('electron')
const path = require('path')
const exec = require('child_process')

const NODE_ENV = process.env.NODE_ENV
    //let port = '10580'

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
    return mainWin
}

app.whenReady().then(() => {
    ChildExec()
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
    ipcMain.on('move-title', (e, pos) => {
        mainWin && mainWin.setPosition(pos.posX, pos.posY)
    })

    ipcMain.on('quit-app', () => {
        app && app.quit()
    })

    ipcMain.on('min-app', () => {
            mainWin && mainWin.minimize()
        })
        /*
        ipcMain.handle('send-port', () => {
            return port
        })
        */
}

async function ChildExec() {
    let e = await require('child_process').spawn(path.join(process.cwd(), '/resources/gRpcTool.exe'))
    e.stdout.on('data', (data) => {
        let myNotification = new Notification({
            title: 'grpc-tool',
            icon: path.join(process.cwd(), '/resources/err.png'),
            body: "错误: " + data.toString()
        })
        myNotification.show()
        app && app.quit()
    })
}