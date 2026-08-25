import { ref, onMounted, onBeforeUnmount } from 'vue'
const isFullscreen = ref(false)
const isMobile = ref(window.innerWidth <= 900)
export default function useScreenAndResizableTable(
    heightRatioMobile,
    widthRatioMobile,
    heightRatioDesktop,
    widthRatioDesktop,
) {
    // 判断是否为移动设备

    // 根据设备类型设置初始表格尺寸
    const heightRatio = isMobile.value
        ? heightRatioMobile ?? 1
        : heightRatioDesktop ?? 1
    const widthRatio = isMobile.value
        ? widthRatioMobile ?? 1
        : widthRatioDesktop ?? 1

    const mainTableHeight = ref(window.innerHeight * heightRatio)
    const mainTableWidth = ref(window.innerWidth * widthRatio)

    // 更新表格尺寸
    const updateTableSize = () => {
        const currentHeightRatio = isMobile.value
            ? heightRatioMobile ?? 1
            : heightRatioDesktop ?? 1
        const currentWidthRatio = isMobile.value
            ? widthRatioMobile ?? 1
            : widthRatioDesktop ?? 1
        mainTableHeight.value = window.innerHeight * currentHeightRatio
        mainTableWidth.value = window.innerWidth * currentWidthRatio
    }

    // 更新是否为移动设备
    const updateIsMobile = () => {
        isMobile.value = window.innerWidth <= 900
        updateTableSize()
    }

    // 全屏控制
    const enterFullscreen = () => {
        const elem = document.documentElement
        const requestFullscreen =
            elem.requestFullscreen ||
            elem.mozRequestFullScreen ||
            elem.webkitRequestFullscreen ||
            elem.msRequestFullscreen
        if (requestFullscreen) {
            requestFullscreen
                .call(elem)
                .then(() => {
                    isFullscreen.value = true
                })
                .catch((err) => {
                    console.error(
                        `尝试启用全屏模式时出错: ${err.message} (${err.name})`,
                    )
                })
        }
    }

    const exitFullscreen = () => {
        const exitFullscreen =
            document.exitFullscreen ||
            document.mozCancelFullScreen ||
            document.webkitExitFullscreen ||
            document.msExitFullscreen
        if (exitFullscreen) {
            exitFullscreen
                .call(document)
                .then(() => {
                    isFullscreen.value = false
                })
                .catch((err) => {
                    console.error(
                        `尝试退出全屏模式时出错: ${err.message} (${err.name})`,
                    )
                })
        }
    }

    const showAlertAndExitFullscreen = (message) => {
        if (isFullscreen.value) {
            exitFullscreen()
        }
        alert(message)
    }

    const handleFullscreenChange = () => {
        isFullscreen.value = !!(
            document.fullscreenElement ||
            document.mozFullScreenElement ||
            document.webkitFullscreenElement ||
            document.msFullscreenElement
        )
    }

    onMounted(() => {
        window.addEventListener('resize', updateIsMobile)
        window.addEventListener('resize', updateTableSize)
        document.documentElement.addEventListener(
            'fullscreenchange',
            handleFullscreenChange,
        )
    })

    onBeforeUnmount(() => {
        window.removeEventListener('resize', updateIsMobile)
        window.removeEventListener('resize', updateTableSize)
        document.documentElement.removeEventListener(
            'fullscreenchange',
            handleFullscreenChange,
        )
    })

    return {
        isFullscreen,
        isMobile,
        mainTableHeight,
        mainTableWidth,
        enterFullscreen,
        exitFullscreen,
        showAlertAndExitFullscreen,
    }
}
