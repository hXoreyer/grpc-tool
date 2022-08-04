export default function () {

    const useDragfunc = (el: HTMLDivElement) => {
        el.onmousedown = (e) => {
            let baseX = e.x
            let baseY = e.y
            document.onmousemove = (e) => {
                let x = e.screenX - baseX;
                let y = e.screenY - baseY;
                (window as any).api.send('move-title', { posX: x, posY: y })
            }

            document.onmouseup = () => {
                document.onmousemove = null
                document.onmouseup = null
            }
        }
    }

    return useDragfunc

}