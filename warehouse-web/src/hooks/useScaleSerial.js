import { computed, ref } from "vue";

const SERIAL_OPTIONS = {
  baudRate: 9600,
  dataBits: 8,
  stopBits: 1,
  parity: "none",
  flowControl: "none",
};

const SCALE_USB_FILTERS = [
  { usbVendorId: 0x067b, usbProductId: 0x2303 },
];

export const SCALE_STATUS = Object.freeze({
  DISCONNECTED: "disconnected",
  CONNECTING: "connecting",
  CONNECTED: "connected",
  LOST: "lost",
  ERROR: "error",
});

export function parseGrossWeight(value) {
  const matches = [
    ...String(value || "").matchAll(
      /G\.W\.\s*:\s*([+-]?\s*\d+(?:\.\d+)?)\s*kg/gi,
    ),
  ];
  if (!matches.length) return null;

  const weight = Number(matches[matches.length - 1][1].replace(/\s/g, ""));
  return Number.isFinite(weight) ? weight : null;
}

export function useScaleSerial() {
  const currentWeight = ref(null);
  const scaleStatus = ref(SCALE_STATUS.DISCONNECTED);
  const scaleErrorCode = ref("");
  const authorizedPortCount = ref(0);

  let serialPort = null;
  let serialReader = null;
  let readTask = null;
  let serialBuffer = "";
  let disconnectRequested = false;
  let listenerAttached = false;

  const isSerialSupported = computed(
    () => typeof navigator !== "undefined" && "serial" in navigator,
  );
  const scaleConnected = computed(
    () => scaleStatus.value === SCALE_STATUS.CONNECTED,
  );
  const scaleConnecting = computed(
    () => scaleStatus.value === SCALE_STATUS.CONNECTING,
  );
  const hasSavedPort = computed(() => authorizedPortCount.value === 1);

  function isScalePort(port) {
    const info = port?.getInfo?.() || {};
    return info.usbVendorId === SCALE_USB_FILTERS[0].usbVendorId
      && info.usbProductId === SCALE_USB_FILTERS[0].usbProductId;
  }

  function updateWeightFromLine(line) {
    const parsedWeight = parseGrossWeight(line);
    if (parsedWeight !== null) currentWeight.value = parsedWeight;
  }

  function consumeSerialBuffer(flush = false) {
    const lines = serialBuffer.split(/[\r\n]+/);
    const remainder = lines.pop() || "";
    lines.forEach(updateWeightFromLine);

    // JWI-700W may leave the last G.W. line in the current serial chunk
    // without a trailing CR/LF. Parse a complete pattern immediately while
    // retaining the remainder so split packets can still be joined safely.
    if (remainder) updateWeightFromLine(remainder);
    serialBuffer = flush ? "" : remainder;
    if (serialBuffer.length > 2048) serialBuffer = serialBuffer.slice(-2048);
  }

  async function readSerialStream(openPort) {
    const decoder = new TextDecoder();
    let readFailed = false;

    try {
      while (!disconnectRequested && openPort.readable) {
        serialReader = openPort.readable.getReader();
        try {
          while (!disconnectRequested) {
            const { value, done } = await serialReader.read();
            if (done) break;
            serialBuffer += decoder.decode(value, { stream: true });
            consumeSerialBuffer();
          }
        } finally {
          serialReader.releaseLock();
          serialReader = null;
        }
      }
    } catch (error) {
      if (!disconnectRequested) {
        readFailed = true;
        console.error("Read scale serial stream failed:", error);
      }
    } finally {
      serialBuffer += decoder.decode();
      consumeSerialBuffer(true);

      if (!disconnectRequested && serialPort === openPort) {
        scaleStatus.value = readFailed
          ? SCALE_STATUS.ERROR
          : SCALE_STATUS.LOST;
        scaleErrorCode.value = readFailed ? "readFailed" : "";
        currentWeight.value = null;
        serialPort = null;
        try {
          await openPort.close();
        } catch {
          // A removed USB device may already be closed by the browser.
        }
      }
    }
  }

  async function refreshAuthorizedPorts() {
    if (!isSerialSupported.value) {
      authorizedPortCount.value = 0;
      return [];
    }

    try {
      const ports = await navigator.serial.getPorts();
      const scalePorts = ports.filter(isScalePort);
      const reusablePorts = scalePorts.length ? scalePorts : ports;
      authorizedPortCount.value = reusablePorts.length;
      return reusablePorts;
    } catch {
      authorizedPortCount.value = 0;
      return [];
    }
  }

  async function connectScale({ reuseSavedPort = false } = {}) {
    if (!isSerialSupported.value) {
      scaleStatus.value = SCALE_STATUS.ERROR;
      scaleErrorCode.value = "unsupported";
      return false;
    }
    if (scaleConnected.value || scaleConnecting.value) return scaleConnected.value;

    scaleStatus.value = SCALE_STATUS.CONNECTING;
    scaleErrorCode.value = "";
    disconnectRequested = false;
    serialBuffer = "";

    try {
      if (reuseSavedPort) {
        const ports = await refreshAuthorizedPorts();
        serialPort = ports.length === 1 ? ports[0] : null;
      }
      if (!serialPort) {
        serialPort = await navigator.serial.requestPort({ filters: SCALE_USB_FILTERS });
      }
    } catch (error) {
      serialPort = null;
      scaleStatus.value = SCALE_STATUS.DISCONNECTED;
      scaleErrorCode.value =
        error?.name === "NotFoundError" ? "permissionCancelled" : "connectionFailed";
      return false;
    }

    try {
      await serialPort.open(SERIAL_OPTIONS);
      scaleStatus.value = SCALE_STATUS.CONNECTED;
      scaleErrorCode.value = "";
      currentWeight.value = null;
      await refreshAuthorizedPorts();
      readTask = readSerialStream(serialPort);
      return true;
    } catch (error) {
      console.error("Open scale serial port failed:", error);
      serialPort = null;
      scaleStatus.value = SCALE_STATUS.ERROR;
      scaleErrorCode.value = "portBusy";
      await refreshAuthorizedPorts();
      return false;
    }
  }

  async function disconnectScale({ keepStatus = false } = {}) {
    disconnectRequested = true;
    const openPort = serialPort;
    serialPort = null;

    if (serialReader) {
      try {
        await serialReader.cancel();
      } catch {
        // The reader may already be cancelled after a physical disconnect.
      }
    }
    if (readTask) {
      try {
        await readTask;
      } catch {
        // Reading errors are reflected through scaleStatus.
      }
      readTask = null;
    }
    if (openPort) {
      try {
        await openPort.close();
      } catch {
        // The browser may already have closed a removed device.
      }
    }

    serialReader = null;
    serialBuffer = "";
    currentWeight.value = null;
    if (!keepStatus) {
      scaleStatus.value = SCALE_STATUS.DISCONNECTED;
      scaleErrorCode.value = "";
    }
    await refreshAuthorizedPorts();
  }

  function resetWeight() {
    currentWeight.value = null;
    serialBuffer = "";
  }

  function handleSerialDisconnect(event) {
    const disconnectedPort = event.port || event.target;
    if (!serialPort || (disconnectedPort !== serialPort && disconnectedPort !== navigator.serial)) {
      return;
    }

    scaleStatus.value = SCALE_STATUS.LOST;
    scaleErrorCode.value = "";
    void disconnectScale({ keepStatus: true });
  }

  async function initializeScale({ autoConnect = false } = {}) {
    if (!isSerialSupported.value) {
      scaleErrorCode.value = "unsupported";
      return;
    }
    const ports = await refreshAuthorizedPorts();
    if (!listenerAttached) {
      navigator.serial.addEventListener("disconnect", handleSerialDisconnect);
      listenerAttached = true;
    }
    if (autoConnect && ports.length === 1) {
      return connectScale({ reuseSavedPort: true });
    }
    return false;
  }

  async function disposeScale() {
    if (listenerAttached && isSerialSupported.value) {
      navigator.serial.removeEventListener("disconnect", handleSerialDisconnect);
      listenerAttached = false;
    }
    await disconnectScale();
  }

  return {
    currentWeight,
    scaleStatus,
    scaleErrorCode,
    isSerialSupported,
    scaleConnected,
    scaleConnecting,
    hasSavedPort,
    initializeScale,
    refreshAuthorizedPorts,
    connectScale,
    disconnectScale,
    resetWeight,
    disposeScale,
  };
}
