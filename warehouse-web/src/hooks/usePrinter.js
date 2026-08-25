import { ref } from "vue";

const PRINTER_IP_KEY = "printerIP";
const PRINTER_MODE_KEY = "printerMode";
const printMode = ref(localStorage.getItem(PRINTER_MODE_KEY) || "wifi");

const isValidIpv4 = (value) => {
  const parts = String(value || "").split(".");
  return (
    parts.length === 4 &&
    parts.every((part) => /^\d{1,3}$/.test(part) && Number(part) >= 0 && Number(part) <= 255)
  );
};

export function usePrinter() {
  const setPrinterMode = (nextMode) => {
    printMode.value = nextMode || (printMode.value === "wifi" ? "ip" : "wifi");
    localStorage.setItem(PRINTER_MODE_KEY, printMode.value);
  };

  const setPrinterIP = async (ip) => {
    if (!isValidIpv4(ip)) return false;
    localStorage.setItem(PRINTER_IP_KEY, ip);
    return true;
  };

  const clearPrinterIP = () => {
    localStorage.removeItem(PRINTER_IP_KEY);
  };

  const printZPL = async (zpl) => {
    const ip = localStorage.getItem(PRINTER_IP_KEY);
    if (!ip || !isValidIpv4(ip) || !zpl) return false;

    try {
      const response = await fetch(`http://${ip}/pstprnt`, {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: zpl,
      });
      return response.ok;
    } catch (error) {
      console.error("Unable to send ZPL to printer:", error);
      return false;
    }
  };

  return {
    printMode,
    setPrinterMode,
    setPrinterIP,
    clearPrinterIP,
    printZPL,
  };
}
