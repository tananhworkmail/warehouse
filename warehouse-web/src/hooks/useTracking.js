// src/hooks/useTracking.js
import axios from "@/plugins/axios";
import { ref, reactive } from "vue";

const API_PREFIX = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/+$/, "");

const url_scan_nhap = `${API_PREFIX}/scan/tracking/meterialimport`;
const url_scan_xuat = `${API_PREFIX}/scan/tracking/meterialexport`;
const url_get_KCLLScan_RFSSS = `${API_PREFIX}/scan/tracking/getkcllscan_rfsss`;
const url_get_KCLLScan_RFSS = `${API_PREFIX}/scan/tracking/getkcllscan_rfss`;
const url_update_KCLLScan_RFSS = `${API_PREFIX}/scan/tracking/updatekcllscan_rfss`;
const url_update_KCRKScan_RFSS = `${API_PREFIX}/scan/tracking/updatekcrkscan_rfss`;
const url_update_KCLLScan_RFSSS = `${API_PREFIX}/scan/tracking/updatekcllscan_rfsss`;
const url_get_print_data = `${API_PREFIX}/scan/tracking/merterialprint`;
const url_update_KCBH_by_barcode = `${API_PREFIX}/scan/tracking/updatekcbhbybarcode`;

const KCLLScan_RFSS_row_list = ref([]);
const KCLLScan_RFSSS_row_list = ref([]);
const Barcode_old = ref("");
const Track_data_list = ref([]);
const Print_data_list = ref([]);

const TrackSearchCondition = reactive({
  inputValues: "",
  selectTimeMode: "",
});

const left_table_selceted_row_colume = reactive({
  row: "",
  colume: "",
});

const right_table_selceted_row_colume = reactive({
  row: "",
  colume: "",
});

const print_file_layout_show = ref(false);
const file_type = ref();
const isloading = ref(false);
const trackingError = ref("");

export default function useTracking() {
  const init_status = () => {
    KCLLScan_RFSS_row_list.value = [];
    KCLLScan_RFSSS_row_list.value = [];
    Barcode_old.value = "";
    Track_data_list.value = [];
    Print_data_list.value = [];
    TrackSearchCondition.inputValues = "";
    TrackSearchCondition.selectTimeMode = "";
    left_table_selceted_row_colume.row = "";
    left_table_selceted_row_colume.colume = "";
    right_table_selceted_row_colume.row = "";
    right_table_selceted_row_colume.colume = "";
    print_file_layout_show.value = false;
    file_type.value = undefined;
    isloading.value = false;
    trackingError.value = "";
  };

  // Tracking public => cố định công ty VDH
  const getGSBH = () => {
    return "VDH";
  };

  const buildDateParams = (inputValues, selectTimeMode) => {
    const [fromDate1, toDate1] = (inputValues?.[0] || "").split("~");
    const [fromDate2, toDate2] = (inputValues?.[1] || "").split("~");

    return selectTimeMode === 0
      ? { SCFMDATE: fromDate1 || "", ECFMDATE: toDate1 || "" }
      : { SCFMDATE: fromDate2 || "", ECFMDATE: toDate2 || "" };
  };

  const buildBaseParams = (inputValues) => {
    const GSBH = getGSBH();

    return {
      CLBH: inputValues?.[2] || "",
      YWBM: inputValues?.[3] || "",
      ZSDH: inputValues?.[4] || "",
      ZSYWJC: inputValues?.[5] || "",
      SCBH: inputValues?.[6] || "",
      CGNO: inputValues?.[7] || "",
      GSBH,
    };
  };

  async function getKCLLScan(row_id) {
    let KCLLScan_RFSS_Data;
    let KCLLScan_RFSSS_Data;

    KCLLScan_RFSS_row_list.value.splice(0, KCLLScan_RFSS_row_list.value.length);
    KCLLScan_RFSSS_row_list.value.splice(
      0,
      KCLLScan_RFSSS_row_list.value.length,
    );

    // try {
    //   const [result1, result2] = await Promise.all([
    //     axios.post(url_get_KCLLScan_RFSS, { Barcode: row_id }),
    //     axios.post(url_get_KCLLScan_RFSSS, { Barcode: row_id }),
    //   ]);

    //   KCLLScan_RFSS_Data = Array.isArray(result1.data.data)
    //     ? result1.data.data
    //     : [result1.data.data];
    //   KCLLScan_RFSSS_Data = Array.isArray(result2.data.data)
    //     ? result2.data.data
    //     : [result2.data.data];

    //   const processItemRFSS = (item) => ({
    //     LLNO: item.LLNO,
    //     ScanQty_input: parseFloat(item.SCANQty),
    //     SCBH: item.SCBH,
    //     DFL: item.DFL,
    //     CLBH: item.CLBH,
    //     SCANUSERID: item.SCANUSERID,
    //     SCANUSERDATE: item.SCANUSERDATE,
    //     Barcode: row_id,
    //     USERID: getUserId(),
    //     SCANQty: item.SCANQty,
    //   });

    //   const processItemRFSSS = (item) => ({
    //     LLNO_OWE: item.LLNO_OWE,
    //     ScanQty_input: parseFloat(item.SCANQty),
    //     SCBH: item.SCBH,
    //     DFL: item.DFL,
    //     SCANUSERID: item.SCANUSERID,
    //     SCANUSERDATE: item.SCANUSERDATE,
    //     Barcode: row_id,
    //     LLNO: item.LLNO,
    //     CLBH: item.CLBH,
    //     USERID: getUserId(),
    //     SCANQty: item.SCANQty,
    //   });

    //   KCLLScan_RFSS_row_list.value.push(
    //     ...KCLLScan_RFSS_Data.filter(Boolean).map(processItemRFSS)
    //   );
    //   KCLLScan_RFSSS_row_list.value.push(
    //     ...KCLLScan_RFSSS_Data.filter(Boolean).map(processItemRFSSS)
    //   );
    // } catch {
    //   alert(
    //     "操作失敗，請確認網路連線後再操作！ Thao tác không thành công, vui lòng kiểm tra lại!",
    //   );
    // }
  }

  function update_KCLLScan_RFSS_row_list(new_KCLLScan_RFSS_row, index) {
    KCLLScan_RFSS_row_list.value[index] = new_KCLLScan_RFSS_row;
  }

  function update_KCLLScan_RFSSS_row_list(new_KCLLScan_RFSSS_row, index) {
    Barcode_old.value = KCLLScan_RFSSS_row_list.value[index].Barcode;
    KCLLScan_RFSSS_row_list.value[index] = new_KCLLScan_RFSSS_row;
  }

  function removeFields(obj, fieldsToRemove) {
    const newObj = { ...obj };
    fieldsToRemove.forEach((field) => {
      delete newObj[field];
    });
    return newObj;
  }

  async function updateKCLLScan(url, list, fieldsToRemove) {
    try {
      await list.reduce(async (previousPromise, item) => {
        await previousPromise;

        const updatedItem = removeFields(item, fieldsToRemove);
        return axios.post(url, updatedItem);
      }, Promise.resolve());
    } catch {
      alert(
        "操作失敗，請確認網路連線後再操作！ Thao tác không thành công, vui lòng kiểm tra lại!",
      );
    }
  }

  async function UpdateKCLLScan_RFSS() {
    const fieldsToRemove = ["SCANUSERDATE"];
    await updateKCLLScan(
      url_update_KCLLScan_RFSS,
      KCLLScan_RFSS_row_list.value,
      fieldsToRemove,
    );

    if (KCLLScan_RFSS_row_list.value.length > 0) {
      UpdateKCRKScan_RFSS(KCLLScan_RFSS_row_list.value[0].Barcode);
    }
  }

  const UpdateKCRKScan_RFSS = async (qrcode) => {
    isloading.value = true;

    try {
      const params = { Barcode: qrcode };
      await axios.post(url_update_KCRKScan_RFSS, params);

      KCLLScan_RFSS_row_list.value.length = 0;
      KCLLScan_RFSSS_row_list.value.length = 0;
      left_table_selceted_row_colume.row = "";
      left_table_selceted_row_colume.colume = "";

      await refreshTrackingData();
    } catch (error) {
      console.error(error);
      alert(
        "操作失敗，請確認網路連線後再操作！ Thao tác không thành công, vui lòng kiểm tra lại!",
      );
    } finally {
      isloading.value = false;
    }
  };

  async function UpdateKCLLScan_RFSSS() {
    const fieldsToRemove = ["SCANUSERDATE", "DFL"];
    await updateKCLLScan(
      url_update_KCLLScan_RFSSS,
      KCLLScan_RFSS_row_list.value,
      fieldsToRemove,
    );

    KCLLScan_RFSS_row_list.value.length = 0;
    KCLLScan_RFSSS_row_list.value.length = 0;
    left_table_selceted_row_colume.row = "";
    left_table_selceted_row_colume.colume = "";

    await refreshTrackingData();
  }

  const updateTrackDataList = (data, selectTimeMode) => {
    if (!data) return;

    const rows = Array.isArray(data) ? data : [data];

    const commonFields = [
      "Memo_RY",
      "GSBH",
      "CKBH",
      "MAKE",
      "CLBH",
      "YWPM",
      "DWBH",
      "Pack",
      "Qty",
      "DQty",
      "RemQty",
      "CGNO",
      "DOCNO",
      "MEMO",
      "ZSBH",
      "ZSYWJC",
      "CFMID",
      "Barcode",
      "Pack1",
    ];

    rows.filter(Boolean).forEach((element) => {
      const Track_Data_Row = {};

      commonFields.forEach((field) => {
        Track_Data_Row[field] = element[field];
      });

      if (selectTimeMode === 1) {
        Track_Data_Row.CFMDATE1 = element.CFMDATE1;
      } else if (selectTimeMode === 0) {
        Track_Data_Row.CFMDATE = element.CFMDATE;
      }

      Track_Data_Row.File_Name = element.File_Name;

      Track_data_list.value.push(Track_Data_Row);
    });
  };

  const setTrackSearchCondition = (inputValues, selectTimeMode) => {
    TrackSearchCondition.inputValues = inputValues;
    TrackSearchCondition.selectTimeMode = selectTimeMode;
  };

  const updatePrintDataList = (data) => {
    Print_data_list.value.length = 0;
    if (data) {
      data.forEach((element) => {
        Print_data_list.value.push(element);
      });
    }
  };

  async function getTrackingData(inputValues, selectTimeMode) {
    Track_data_list.value.length = 0;
    KCLLScan_RFSS_row_list.value.length = 0;
    KCLLScan_RFSSS_row_list.value.length = 0;
    isloading.value = true;
    trackingError.value = "";

    try {
      const dateParams = buildDateParams(inputValues, selectTimeMode);
      const keyParams = buildBaseParams(inputValues);

      const MAKE = inputValues?.[8] || "";
      const Barcode = inputValues?.[9] || "";

      if (MAKE && Barcode) {
        alert("Chỉ được nhập MAKE hoặc Barcode");
        return;
      }

      if (MAKE) keyParams.MAKE = MAKE;
      if (Barcode) keyParams.Barcode = Barcode;

      const url = selectTimeMode === 0 ? url_scan_nhap : url_scan_xuat;
      const result = await axios.post(url, {
        ...dateParams,
        ...keyParams,
      });

      updateTrackDataList(result.data.data, selectTimeMode);
      setTrackSearchCondition(inputValues, selectTimeMode);
    } catch (error) {
      console.error("Tracking search failed:", error);
      trackingError.value =
        error.response?.data?.message ||
        "Không thể tải dữ liệu tracking. Vui lòng thử lại.";
    } finally {
      isloading.value = false;
    }
  }

  async function refreshTrackingData() {
    Track_data_list.value.length = 0;
    isloading.value = true;
    trackingError.value = "";

    try {
      const inputValues = TrackSearchCondition.inputValues;
      const selectTimeMode = TrackSearchCondition.selectTimeMode;

      if (
        !inputValues ||
        !Array.isArray(inputValues) ||
        inputValues.length < 10
      ) {
        return;
      }

      const dateParams = buildDateParams(inputValues, selectTimeMode);
      const keyParams = buildBaseParams(inputValues);

      const MAKE = inputValues?.[8] || "";
      const Barcode = inputValues?.[9] || "";

      if (MAKE) keyParams.MAKE = MAKE;
      if (Barcode) keyParams.Barcode = Barcode;

      const url = selectTimeMode === 0 ? url_scan_nhap : url_scan_xuat;
      const result = await axios.post(url, {
        ...dateParams,
        ...keyParams,
      });

      updateTrackDataList(result.data.data, selectTimeMode);
    } catch (error) {
      console.error("Tracking refresh failed:", error);
      trackingError.value =
        error.response?.data?.message ||
        "Không thể tải lại dữ liệu tracking. Vui lòng thử lại.";
    } finally {
      isloading.value = false;
    }
  }

  async function getTrackDataToPrint(inputValues, selectTimeMode) {
    isloading.value = true;

    try {
      const dateParams = buildDateParams(inputValues, selectTimeMode);
      const keyParams = buildBaseParams(inputValues);

      const MAKE = inputValues?.[8] || "";
      const Barcode = inputValues?.[9] || "";

      if (MAKE && Barcode) {
        alert("Chỉ được nhập MAKE hoặc Barcode");
        return;
      }

      if (MAKE) keyParams.MAKE = MAKE;
      if (Barcode) keyParams.Barcode = Barcode;

      const url = selectTimeMode === 0 ? url_get_print_data : url_scan_xuat;
      file_type.value = selectTimeMode === 0 ? "import" : "export";

      const result = await axios.post(url, {
        ...dateParams,
        ...keyParams,
      });

      updatePrintDataList(result.data.data);
      print_file_layout_show.value = Print_data_list.value.length > 0;
    } catch (error) {
      console.error("Print Error:", error);
      alert("Thao tác không thành công, vui lòng kiểm tra kết nối mạng!");
    } finally {
      isloading.value = false;
    }
  }

  async function updateKCBHByBarcode(payload) {
    if (!payload.barcode || !payload.newKCBH) {
      alert("Vui lòng nhập đầy đủ Barcode và MAKE mới!");
      return;
    }

    isloading.value = true;

    try {
      const formattedPayload = {
        barcode: payload.barcode.trim(),
        newKCBH: payload.newKCBH.toUpperCase().trim(),
      };

      const response = await axios.post(
        url_update_KCBH_by_barcode,
        formattedPayload,
      );

      return response.data;
    } catch (err) {
      console.error("Update MAKE Error:", err);
      const msg =
        err.response?.data?.message ||
        "Thao tác không thành công, vui lòng kiểm tra mạng!";
      alert(msg);
      throw err;
    } finally {
      isloading.value = false;
    }
  }

  return {
    Track_data_list,
    Print_data_list,
    KCLLScan_RFSS_row_list,
    KCLLScan_RFSSS_row_list,
    Barcode_old,
    getTrackingData,
    getKCLLScan,
    left_table_selceted_row_colume,
    right_table_selceted_row_colume,
    update_KCLLScan_RFSS_row_list,
    update_KCLLScan_RFSSS_row_list,
    UpdateKCLLScan_RFSS,
    UpdateKCLLScan_RFSSS,
    getTrackDataToPrint,
    init_status,
    print_file_layout_show,
    isloading,
    trackingError,
    file_type,
    updateKCBHByBarcode,
    refreshTrackingData,
  };
}
