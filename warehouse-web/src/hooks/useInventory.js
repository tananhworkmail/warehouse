import axios from "@/plugins/axios";
import { ref } from "vue";

const API_PREFIX = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/+$/, "");
const INVENTORY_DAY_URL = `${API_PREFIX}/scan/number-scan/inventory`;
const INVENTORY_MONTH_URL = `${API_PREFIX}/scan/number-scan/inventory-month`;
const RACK_LIST_URL = `${API_PREFIX}/scan/number-scan/rack-list`;

const isSuccessfulResponse = (data) => Number(data?.code) === 200;

export function useInventory() {
  const inventoryList = ref([]);
  const rackList = ref([]);
  const isPending = ref(false);

  const fetchRackList = async (gsbh, signal) => {
    if (!gsbh) {
      rackList.value = [];
      return [];
    }

    try {
      const res = await axios.get(RACK_LIST_URL, {
        params: { GSBH: gsbh },
        signal,
      });
      rackList.value = isSuccessfulResponse(res.data) ? res.data.data || [] : [];
      return rackList.value;
    } catch (error) {
      if (error?.code !== "ERR_CANCELED") console.error("fetchRackList failed:", error);
      rackList.value = [];
      if (error?.code !== "ERR_CANCELED") throw error;
      return [];
    }
  };

  const fetchInventory = async (date, gsbh, rack = "", mode = "day", signal) => {
    const url = mode === "month" ? INVENTORY_MONTH_URL : INVENTORY_DAY_URL;
    isPending.value = true;

    try {
      const res = await axios.get(url, {
        params: { Date: date, GSBH: gsbh, Rack: rack },
        signal,
      });

      if (!isSuccessfulResponse(res.data)) {
        inventoryList.value = [];
        throw new Error(res.data?.message || "Unable to load inventory data.");
      }

      const rawData = Array.isArray(res.data?.data) ? res.data.data : [];
      const company = String(gsbh || "").toUpperCase();

      if (company === "HBA") {
        inventoryList.value = rawData.filter((item) =>
          String(item?.Rack || "").toUpperCase().startsWith("B"),
        );
      } else if (company === "TBA") {
        inventoryList.value = rawData.filter((item) =>
          String(item?.Rack || "").toUpperCase().startsWith("A"),
        );
      } else {
        inventoryList.value = rawData;
      }

      return inventoryList.value;
    } catch (error) {
      if (error?.code === "ERR_CANCELED") throw error;
      console.error("fetchInventory failed:", error);
      inventoryList.value = [];
      throw error;
    } finally {
      isPending.value = false;
    }
  };

  const saveInventory = async (payload) => {
    const url = payload.MODE === "month" ? INVENTORY_MONTH_URL : INVENTORY_DAY_URL;
    isPending.value = true;

    try {
      const res = await axios.post(url, payload);
      const data = res.data || {};

      // A warning is a valid business response (for example, a duplicate scan),
      // but an explicit non-200 response must never be reported as a successful scan.
      if (data.type !== "warning" && data.code !== undefined && Number(data.code) !== 200) {
        throw new Error(data.message || "Unable to save inventory scan.");
      }

      return data;
    } catch (error) {
      console.error("saveInventory failed:", error);
      throw error;
    } finally {
      isPending.value = false;
    }
  };

  return {
    inventoryList,
    rackList,
    isPending,
    fetchInventory,
    fetchRackList,
    saveInventory,
  };
}
